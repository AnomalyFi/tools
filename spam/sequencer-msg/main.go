package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/AnomalyFi/hypersdk/codec"
	"github.com/AnomalyFi/hypersdk/crypto/ed25519"
	"github.com/AnomalyFi/hypersdk/fees"
	"github.com/AnomalyFi/hypersdk/utils"
	"golang.org/x/sync/errgroup"

	"github.com/AnomalyFi/hypersdk/chain"
	"github.com/AnomalyFi/hypersdk/pubsub"
	"github.com/AnomalyFi/hypersdk/rpc"
	"github.com/AnomalyFi/nodekit-seq/actions"
	"github.com/AnomalyFi/nodekit-seq/auth"
	"github.com/AnomalyFi/nodekit-seq/consts"
	trpc "github.com/AnomalyFi/nodekit-seq/rpc"
	"github.com/ava-labs/avalanchego/ids"
)

var (
	ErrTxFailed = errors.New("tx failed on-chain")
)

const (
	defaultRange          = 32
	issuerShutdownTimeout = 60 * time.Second
)

var (
	issuerWg sync.WaitGroup
	exiting  sync.Once

	l            sync.Mutex
	confirmedTxs uint64
	totalTxs     uint64

	inflight atomic.Int64
	sent     atomic.Int64
)

const (
	decimals              = 9
	maxTxBacklog          = 500
	maxFee                = -1
	randomRecipient       = false
	numAccounts           = 1
	numTxsPerAccount      = 1 // per second
	numClients            = 1
	MillisecondsPerSecond = 1000 // 1000ms = 1 sec
)

type PrivateKey struct {
	Address codec.Address
	Bytes   []byte
}

type txIssuer struct {
	c *rpc.JSONRPCClient
	d *rpc.WebSocketClient

	l              sync.Mutex
	uri            int
	abandoned      error
	outstandingTxs int
}

func main() {
	ctx := context.Background()

	// chain id
	chainID, _ := ids.FromString("cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7")
	uris := []string{"http://34.229.58.146:9650/ext/bc/cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7",
		"http://54.91.39.46:9650/ext/bc/cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7",
		"http://34.227.206.152:9650/ext/bc/cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7",
		"http://18.215.159.187:9650/ext/bc/cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7",
		"http://54.242.101.165:9650/ext/bc/cKA3rhvogANuQV6y8hXX9282tVVvDZQBoVUyRgBoXUZhnPjN7",
	}

	// root private key, with all the funds:
	privBytes, _ := codec.LoadHex(
		"323b1d8f4eed5f0da9da93071b034f2dce9d2d22692c172f3cb252a64ddfafd01b057de320297c29ad0c1f589ea216869cf1938d88c9fbd70d6748323dbf2fa7", //nolint:lll
		ed25519.PrivateKeyLen,
	)

	priv := ed25519.PrivateKey(privBytes)
	factory := auth.NewED25519Factory(priv)
	address := auth.NewED25519Address(priv.PublicKey())
	sddr, _ := codec.AddressBech32("token", address)
	cli := rpc.NewJSONRPCClient(uris[1])
	networkID, _, _, err := cli.Network(ctx)
	if err != nil {
		panic(err)
	}
	tclient, _, err := createClient(uris[1], networkID, chainID)
	if err != nil {
		panic(err)
	}
	balance, err := lookupBalance(tclient, sddr)
	if err != nil {
		panic(err)
	}
	chainIDs := [][]byte{[]byte("nkit"), []byte("everest"), []byte("combator"), []byte("marinedrive")}
	actions := getSequencerMessage(address, chainIDs[0], 300)
	parser, err := tclient.Parser(ctx)
	if err != nil {
		panic(err)
	}
	maxUnits, err := chain.EstimateUnits(parser.Rules(time.Now().UnixMilli()), actions, factory)
	if err != nil {
		panic(err)
	}

	// Distribute funds to accounts:
	unitPrices, err := cli.UnitPrices(ctx, false)
	if err != nil {
		panic(err)
	}

	feePerTx, err := fees.MulSum(unitPrices, maxUnits)
	if err != nil {
		panic(err)
	}
	witholding := feePerTx * uint64(numAccounts)
	if balance < witholding {
		panic(fmt.Errorf("insufficient funds (have=%d need=%d)", balance, witholding))
	}
	distAmount := (balance - witholding) / uint64(numAccounts)
	utils.Outf(
		"{{yellow}}distributing funds to each account:{{/}} %s %s\n",
		utils.FormatBalance(distAmount, decimals),
		"SEQ",
	)
	accounts := make([]*PrivateKey, numAccounts)
	dcli, err := rpc.NewWebSocketClient(uris[0], rpc.DefaultHandshakeTimeout, pubsub.MaxPendingMessages, pubsub.MaxReadMessageSize) // we write the max read
	if err != nil {
		panic(err)
	}
	funds := map[codec.Address]uint64{}
	var fundsL sync.Mutex
	for i := 0; i < numAccounts; i++ {
		// Create account
		pk, err := createAccount()
		if err != nil {
			panic(err)
		}
		accounts[i] = pk

		// Send funds
		_, tx, err := cli.GenerateTransactionManual(parser, generateTransfer(pk.Address, distAmount), factory, feePerTx)
		if err != nil {
			panic(err)
		}
		if err := dcli.RegisterTx(tx); err != nil {
			panic(fmt.Errorf("%w: failed to register tx", err))
		}
		funds[pk.Address] = distAmount
	}

	for i := 0; i < numAccounts; i++ {
		_, dErr, result, err := dcli.ListenTx(ctx)
		if err != nil {
			panic(err)
		}
		if dErr != nil {
			panic(dErr)
		}
		if !result.Success {
			// Should never happen
			panic(fmt.Errorf("%w: %s", ErrTxFailed, result.Error))
		}
	}
	var recipientFunc func() (*PrivateKey, error)
	utils.Outf("{{yellow}}distributed funds to %d accounts{{/}}\n", numAccounts)
	// Kickoff txs
	clients := []*txIssuer{}
	for i := 0; i < len(uris); i++ {
		for j := 0; j < numClients; j++ {
			cli := rpc.NewJSONRPCClient(uris[i])
			dcli, err := rpc.NewWebSocketClient(uris[i], rpc.DefaultHandshakeTimeout, pubsub.MaxPendingMessages, pubsub.MaxReadMessageSize) // we write the max read
			if err != nil {
				panic(err)
			}
			clients = append(clients, &txIssuer{c: cli, d: dcli, uri: i})
		}
	}
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// confirm txs (track failure rate)
	unitPrices, err = clients[0].c.UnitPrices(ctx, false)
	if err != nil {
		panic(err)
	}
	PrintUnitPrices(unitPrices)
	cctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, client := range clients {
		startIssuer(cctx, client)
	}

	// log stats
	t := time.NewTicker(1 * time.Second) // ensure no duplicates created
	defer t.Stop()
	var psent int64
	go func() {
		for {
			select {
			case <-t.C:
				current := sent.Load()
				l.Lock()
				if totalTxs > 0 {
					unitPrices, err = clients[0].c.UnitPrices(ctx, false)
					if err != nil {
						continue
					}
					utils.Outf(
						"{{yellow}}txs seen:{{/}} %d {{yellow}}success rate:{{/}} %.2f%% {{yellow}}inflight:{{/}} %d {{yellow}}issued/s:{{/}} %d {{yellow}}unit prices:{{/}} [%s]\n", //nolint:lll
						totalTxs,
						float64(confirmedTxs)/float64(totalTxs)*100,
						inflight.Load(),
						current-psent,
						ParseDimensions(unitPrices),
					)
				}
				l.Unlock()
				psent = current
			case <-cctx.Done():
				return
			}
		}
	}()

	// broadcast txs
	g, gctx := errgroup.WithContext(ctx)
	for ri := 0; ri < numAccounts; ri++ {
		i := ri
		g.Go(func() error {
			t := time.NewTimer(0) // ensure no duplicates created
			defer t.Stop()
			source := rand.NewSource(time.Now().UnixNano())
			issuerIndex, issuer := getRandomIssuer(clients)
			factory, err := getFactory(accounts[i])
			if err != nil {
				return err
			}
			fundsL.Lock()
			balance := funds[accounts[i].Address]
			fundsL.Unlock()
			defer func() {
				fundsL.Lock()
				funds[accounts[i].Address] = balance
				fundsL.Unlock()
			}()
			ut := time.Now().Unix()
			for {
				select {
				case <-t.C:
					// Ensure we aren't too backlogged
					if inflight.Load() > int64(maxTxBacklog) {
						t.Reset(1 * time.Second)
						continue
					}

					// Select tx time
					//
					// Needed to prevent duplicates if called within the same
					// unix second.
					nextTime := time.Now().Unix()
					if nextTime <= ut {
						nextTime = ut + 1
					}
					ut = nextTime
					tm := &timeModifier{nextTime*MillisecondsPerSecond + parser.Rules(nextTime).GetValidityWindow() - 5*MillisecondsPerSecond}

					// Send transaction
					start := time.Now()
					selected := map[codec.Address]int{}
					for k := 0; k < numTxsPerAccount; k++ {
						recipient, err := getNextRecipient(i, recipientFunc, accounts)
						if err != nil {
							utils.Outf("{{orange}}failed to get next recipient:{{/}} %v\n", err)
							return err
						}
						v := selected[recipient] + 1
						selected[recipient] = v
						actions := getSequencerMessage(recipient, chainIDs[rand.Int()%len(chainIDs)], source.Int63()/1_200_000)
						fee, err := fees.MulSum(unitPrices, maxUnits)
						if err != nil {
							utils.Outf("{{orange}}failed to estimate max fee:{{/}} %v\n", err)
							return err
						}

						_, tx, err := issuer.c.GenerateTransactionManual(parser, actions, factory, fee, tm)
						if err != nil {
							utils.Outf("{{orange}}failed to generate tx:{{/}} %v\n", err)
							continue
						}
						if err := issuer.d.RegisterTx(tx); err != nil {
							issuer.l.Lock()
							if issuer.d.Closed() {
								if issuer.abandoned != nil {
									issuer.l.Unlock()
									return issuer.abandoned
								}
								// recreate issuer
								utils.Outf("{{orange}}re-creating issuer:{{/}} %d {{orange}}uri:{{/}} %d\n", issuerIndex, issuer.uri)
								dcli, err := rpc.NewWebSocketClient(uris[issuer.uri], rpc.DefaultHandshakeTimeout, pubsub.MaxPendingMessages, pubsub.MaxReadMessageSize) // we write the max read
								if err != nil {
									issuer.abandoned = err
									utils.Outf("{{orange}}could not re-create closed issuer:{{/}} %v\n", err)
									issuer.l.Unlock()
									return err
								}
								issuer.d = dcli
								startIssuer(cctx, issuer)
								issuer.l.Unlock()
								utils.Outf("{{green}}re-created closed issuer:{{/}} %d\n", issuerIndex)
							}
							continue
						}
						balance -= (fee + uint64(v))
						issuer.l.Lock()
						issuer.outstandingTxs++
						issuer.l.Unlock()
						inflight.Add(1)
						sent.Add(1)
					}

					// Determine how long to sleep
					dur := time.Since(start)
					sleep := max(float64(MillisecondsPerSecond-dur.Milliseconds()), 0)
					t.Reset(time.Duration(sleep) * time.Millisecond)
				case <-gctx.Done():
					return gctx.Err()
				case <-cctx.Done():
					return nil
				case <-signals:
					exiting.Do(func() {
						utils.Outf("{{yellow}}exiting broadcast loop{{/}}\n")
						cancel()
					})
					return nil
				}
			}
		})
	}
	if err := g.Wait(); err != nil {
		panic(err)
	}
}

func createClient(uri string, networkID uint32, chainID ids.ID) (*trpc.JSONRPCClient, *rpc.WebSocketClient, error) {
	tclient := trpc.NewJSONRPCClient(uri, networkID, chainID)
	sc, err := rpc.NewWebSocketClient(uri, rpc.DefaultHandshakeTimeout, pubsub.MaxPendingMessages, pubsub.MaxReadMessageSize)
	if err != nil {
		return nil, nil, err
	}
	sclient := sc
	return tclient, sclient, nil
}

func getFactory(priv *PrivateKey) (chain.AuthFactory, error) {
	return auth.NewED25519Factory(ed25519.PrivateKey(priv.Bytes)), nil
}

func createAccount() (*PrivateKey, error) { // createAccount
	p, err := ed25519.GeneratePrivateKey()
	if err != nil {
		return nil, err
	}
	return &PrivateKey{
		Address: auth.NewED25519Address(p.PublicKey()),
		Bytes:   p[:],
	}, nil
}

func lookupBalance(tclient *trpc.JSONRPCClient, address string) (uint64, error) {
	balance, err := tclient.Balance(context.TODO(), address, ids.Empty)
	if err != nil {
		return 0, err
	}
	utils.Outf(
		"{{cyan}}address:{{/}} %s {{cyan}}balance:{{/}} %s %s\n",
		address,
		utils.FormatBalance(balance, consts.Decimals),
		consts.Symbol,
	)
	return balance, err
}

// @todo
func getSequencerMessage(addr codec.Address, chainID []byte, dataLen int64) []chain.Action {
	data := make([]byte, dataLen)
	_, err := rand.Read(data)
	if err != nil {
		fmt.Println("error getting random data", err)
	}
	return []chain.Action{&actions.SequencerMsg{
		ChainId:     chainID,
		Data:        data,
		FromAddress: addr,
		RelayerID:   int(dataLen % 10),
	}}
}

func PrintUnitPrices(d fees.Dimensions) {
	utils.Outf(
		"{{cyan}}unit prices{{/}} {{yellow}}bandwidth:{{/}} %d {{yellow}}compute:{{/}} %d {{yellow}}storage(read):{{/}} %d {{yellow}}storage(allocate):{{/}} %d {{yellow}}storage(write):{{/}} %d\n",
		d[fees.Bandwidth],
		d[fees.Compute],
		d[fees.StorageRead],
		d[fees.StorageAllocate],
		d[fees.StorageWrite],
	)
}

func ParseDimensions(d fees.Dimensions) string {
	return fmt.Sprintf(
		"bandwidth=%d compute=%d storage(read)=%d storage(allocate)=%d storage(write)=%d",
		d[fees.Bandwidth],
		d[fees.Compute],
		d[fees.StorageRead],
		d[fees.StorageAllocate],
		d[fees.StorageWrite],
	)
}

type timeModifier struct {
	Timestamp int64
}

func (t *timeModifier) Base(b *chain.Base) {
	b.Timestamp = t.Timestamp
}

func startIssuer(cctx context.Context, issuer *txIssuer) {
	issuerWg.Add(1)
	go func() {
		for {
			_, dErr, result, err := issuer.d.ListenTx(context.TODO())
			if err != nil {
				return
			}
			inflight.Add(-1)
			issuer.l.Lock()
			issuer.outstandingTxs--
			issuer.l.Unlock()
			l.Lock()
			if result != nil {
				if result.Success {
					confirmedTxs++
				} else {
					utils.Outf("{{orange}}on-chain tx failure:{{/}} %s %t\n", string(result.Error), result.Success)
				}
			} else {
				// We can't error match here because we receive it over the wire.
				if !strings.Contains(dErr.Error(), rpc.ErrExpired.Error()) {
					utils.Outf("{{orange}}pre-execute tx failure:{{/}} %v\n", dErr)
				}
			}
			totalTxs++
			l.Unlock()
		}
	}()
	go func() {
		defer func() {
			_ = issuer.d.Close()
			issuerWg.Done()
		}()

		<-cctx.Done()
		start := time.Now()
		for time.Since(start) < issuerShutdownTimeout {
			if issuer.d.Closed() {
				return
			}
			issuer.l.Lock()
			outstanding := issuer.outstandingTxs
			issuer.l.Unlock()
			if outstanding == 0 {
				return
			}
			time.Sleep(500 * time.Millisecond)
		}
		utils.Outf("{{orange}}issuer shutdown timeout{{/}}\n")
	}()
}

func getRandomIssuer(issuers []*txIssuer) (int, *txIssuer) {
	index := rand.Int() % len(issuers)
	return index, issuers[index]
}

func getNextRecipient(self int, createAccount func() (*PrivateKey, error), keys []*PrivateKey) (codec.Address, error) {
	// Send to a random, new account
	// if createAccount != nil {
	// 	priv, err := createAccount()
	// 	if err != nil {
	// 		return codec.EmptyAddress, err
	// 	}
	// 	return priv.Address, nil
	// }

	// Select item from array
	index := rand.Int() % len(keys)
	if index == self {
		index++
		if index == len(keys) {
			index = 0
		}
	}
	return keys[index].Address, nil
}

func generateTransfer(addr codec.Address, amount uint64) []chain.Action { // getTransfer
	return []chain.Action{&actions.Transfer{
		To:    addr,
		Asset: ids.Empty,
		Value: amount,
	}}
}
