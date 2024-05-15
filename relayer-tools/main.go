package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	hrpc "github.com/AnomalyFi/hypersdk/rpc"
	srpc "github.com/AnomalyFi/nodekit-seq/rpc"
	"github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
)

type ConfigData struct {
	DAConfig         EigenDAClientConfig `json:"daConfig"`
	DBConfig         DatabaseConfig      `json:"dbConfig"`
	SeqNode          SeqNodeInfo         `json:"seqNode"`
	Rules            Rules               `json:"rules"`
	Log              LogConfig           `json:"log"`
	ServerLessConfig ServerLessConfig    `json:"serverLessConfig"`
}

type EigenDAClientConfig struct {
	Target             string `json:"target"`
	PrivateKeyFilePath string `json:"privateKeyFilePath"`
}

type DatabaseConfig struct {
	File         string `json:"file"`
	RetainWindow uint64 `json:"retainWindow"`
}

type SeqNodeInfo struct {
	URI       string `json:"uri"`
	NetworkID uint32 `json:"networkID"`
	NodeUrl   string `json:"nodeUrl"`
	// need to be unmarshalled
	ChainID string `json:"chainID"`
}

type SeqWSClientConfig struct {
	URI string `json:"uri"`
}

type SeqJsonRPCConfig struct {
	URI       string `json:"uri"`
	NetworkID uint32 `json:"networkID"`
	// need to be unmarshalled
	ChainID string `json:"chainID"`
}

type Rules struct {
	UseStableSeqHeight bool   `json:"useStableSeqHeight"`
	StableSeqHeight    uint64 `json:"stableSeqHeight"`
}

type LogConfig struct {
	Level     string `json:"level"`
	ToConsole bool   `json:"toConsole"`
	LogFile   string `json:"logFile"`
}

type ServerLessConfig struct {
	Endpoint        string `json:"endpoint"`
	BlocksPerWindow uint32 `json:"blocksPerWindow"`
	RelayerID       int    `json:"relayerID"`
	WaitBlocks      uint32 `json:"waitBlocks"`
}

func main() {
	ctx := context.Background()

	// read config file tempelate
	args := os.Args[1:]
	if len(args) < 1 {
		panic("Please specify config file path")
	}
	fmt.Println(args)

	configBytes, err := os.ReadFile(args[0])
	if err != nil {
		panic(fmt.Sprintf("unable to open file %s\n", args[0]))
	}

	var config ConfigData
	if err := json.Unmarshal(configBytes, &config); err != nil {
		panic(fmt.Sprintf("unable to parse config:\n %s\n", string(configBytes)))
	}

	// Load new items from ANR
	anrCli, err := client.New(client.Config{
		Endpoint:    "0.0.0.0:12352",
		DialTimeout: 10 * time.Second,
	}, logging.NoLog{})
	if err != nil {
		panic(err)
	}
	status, err := anrCli.Status(ctx)
	if err != nil {
		panic(err)
	}
	subnets := map[ids.ID][]ids.ID{}
	uris := []string{}
	nodeUrls := []string{}
	for chain, chainInfo := range status.ClusterInfo.CustomChains {
		chainID, err := ids.FromString(chain)
		if err != nil {
			panic(err)
		}
		subnetID, err := ids.FromString(chainInfo.SubnetId)
		if err != nil {
			panic(err)
		}
		chainIDs, ok := subnets[subnetID]
		if !ok {
			chainIDs = []ids.ID{}
		}
		chainIDs = append(chainIDs, chainID)
		subnets[subnetID] = chainIDs
	}
	var filledChainID ids.ID
	for _, nodeInfo := range status.ClusterInfo.NodeInfos {
		if len(nodeInfo.WhitelistedSubnets) == 0 {
			continue
		}
		trackedSubnets := strings.Split(nodeInfo.WhitelistedSubnets, ",")
		for _, subnet := range trackedSubnets {
			subnetID, err := ids.FromString(subnet)
			if err != nil {
				panic(err)
			}
			for _, chainID := range subnets[subnetID] {
				uri := fmt.Sprintf("%s/ext/bc/%s", nodeInfo.Uri, chainID)
				uris = append(uris, uri)
				nodeUrls = append(nodeUrls, nodeInfo.Uri)
				filledChainID = chainID
			}
		}
	}
	// fetch network id
	hcli := hrpc.NewJSONRPCClient(uris[0])
	networkID, _, _, err := hcli.Network(ctx)
	if err != nil {
		panic(err)
	}
	// fetch serverless ports for every node
	for i, uri := range uris {
		cli := srpc.NewJSONRPCClient(uri, networkID, filledChainID)
		port, err := cli.ServerlessPort(ctx)
		if err != nil {
			panic(err)
		}
		// modify config file
		config.SeqNode.URI = uri
		config.SeqNode.NetworkID = networkID
		config.SeqNode.ChainID = filledChainID.String()
		config.ServerLessConfig.Endpoint = "localhost" + port
		config.SeqNode.NodeUrl = nodeUrls[i]
		config.Log.LogFile = "./relayer" + strconv.Itoa(i) + ".log"
		config.DBConfig.File = "./relayer" + strconv.Itoa(i) + ".db"
		config.DAConfig.PrivateKeyFilePath = "./demo" + strconv.Itoa(i) + ".pk"
		// create new config file(s)
		d, err := json.Marshal(config)
		if err != nil {
			panic(err)
		}
		fmt.Println(os.WriteFile("config"+strconv.Itoa(i)+".json", d, 0644))
	}
}