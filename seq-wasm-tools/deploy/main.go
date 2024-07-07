package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	"github.com/AnomalyFi/hypersdk/chain"
	"github.com/AnomalyFi/hypersdk/codec"
	"github.com/AnomalyFi/hypersdk/crypto/ed25519"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"

	"github.com/AnomalyFi/hypersdk/rpc"
	"github.com/AnomalyFi/nodekit-seq/actions"
	"github.com/AnomalyFi/nodekit-seq/auth"
	trpc "github.com/AnomalyFi/nodekit-seq/rpc"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/coreth/accounts/abi"
	"github.com/ava-labs/coreth/accounts/abi/bind"
	"github.com/ava-labs/coreth/core/types"

	"github.com/AnomalyFi/tools/state-keys/blobstream"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// BinaryMerkleProof is an auto generated low-level Go binding around an user-defined struct.
type BinaryMerkleProof struct {
	SideNodes [][32]byte
	Key       *big.Int
	NumLeaves *big.Int
}

// CommitHeaderRangeInput is an auto generated low-level Go binding around an user-defined struct.
type CommitHeaderRangeInput struct {
	Proof        []byte
	PublicValues []byte
}

// DataRootTuple is an auto generated low-level Go binding around an user-defined struct.
type DataRootTuple struct {
	Height   *big.Int
	DataRoot [32]byte
}

// InitializerInput is an auto generated low-level Go binding around an user-defined struct.
type InitializerInput struct {
	Height                    uint64
	Header                    [32]byte
	BlobstreamProgramVKeyHash []byte
	BlobstreamProgramVKey     []byte
}

// UpdateFreezeInput is an auto generated low-level Go binding around an user-defined struct.
type UpdateFreezeInput struct {
	Freeze bool
}

// UpdateGenesisStateInput is an auto generated low-level Go binding around an user-defined struct.
type UpdateGenesisStateInput struct {
	Height uint64
	Header [32]byte
}

// UpdateProgramVkeyInput is an auto generated low-level Go binding around an user-defined struct.
type UpdateProgramVkeyInput struct {
	BlobstreamProgramVKeyHash []byte
	BlobstreamProgramVKey     []byte
}

// VAInput is an auto generated low-level Go binding around an user-defined struct.
type VAInput struct {
	TupleRootNonce *big.Int
	Tuple          DataRootTuple
	Proof          BinaryMerkleProof
}

// BlobStreamInputsMetaData contains all meta data concerning the BlobStreamInputs contract.
var BlobStreamInputsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"proof\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"publicValues\",\"type\":\"bytes\"}],\"internalType\":\"structCommitHeaderRangeInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"commitHeaderRange\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"header\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"blobstreamProgramVKeyHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"blobstreamProgramVKey\",\"type\":\"bytes\"}],\"internalType\":\"structInitializerInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"initializer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bool\",\"name\":\"freeze\",\"type\":\"bool\"}],\"internalType\":\"structUpdateFreezeInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"updateFreeze\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"height\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"header\",\"type\":\"bytes32\"}],\"internalType\":\"structUpdateGenesisStateInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"updateGenesisState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"blobstreamProgramVKeyHash\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"blobstreamProgramVKey\",\"type\":\"bytes\"}],\"internalType\":\"structUpdateProgramVkeyInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"updateProgramVkey\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"tuple_root_nonce\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"dataRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structDataRootTuple\",\"name\":\"tuple\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32[]\",\"name\":\"sideNodes\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256\",\"name\":\"key\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numLeaves\",\"type\":\"uint256\"}],\"internalType\":\"structBinaryMerkleProof\",\"name\":\"proof\",\"type\":\"tuple\"}],\"internalType\":\"structVAInput\",\"name\":\"inputs\",\"type\":\"tuple\"}],\"name\":\"verifyAppend\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

var BlobStreamInputsABI, _ = BlobStreamInputsMetaData.GetAbi()

func main() {
	ctx := context.Background()
	privBytes, _ := codec.LoadHex(
		"323b1d8f4eed5f0da9da93071b034f2dce9d2d22692c172f3cb252a64ddfafd01b057de320297c29ad0c1f589ea216869cf1938d88c9fbd70d6748323dbf2fa7", //nolint:lll
		ed25519.PrivateKeyLen,
	)
	uri := "http://127.0.0.1:9658/ext/bc/tEpDFmDWyU4C7FCUYLg7YNudkJsRo6AQyvksmLAaAJ14yM1cs"
	var networkID uint32 = 1337
	chainID, _ := ids.FromString("tEpDFmDWyU4C7FCUYLg7YNudkJsRo6AQyvksmLAaAJ14yM1cs")
	cli := rpc.NewJSONRPCClient(uri)
	tcli := trpc.NewJSONRPCClient(uri, networkID, chainID)
	parser, err := tcli.Parser(ctx)
	if err != nil {
		panic(err)
	}
	factory := auth.NewED25519Factory(ed25519.PrivateKey(privBytes))

	contractCode, err := os.ReadFile("/home/ubuntu/seq-wasm/target/wasm32-unknown-unknown/release/blobstream_contracts_rust.wasm")
	if err != nil {
		panic(err)
	}
	vkey, err := os.ReadFile("/home/ubuntu/tools/seq-wasm-tools/vk.bin")
	if err != nil {
		panic(err)
	}

	// trusted_header_hash: 0x188b708bee180f43e3a252471754fd35283a6b09a6fd02f5b9130cc15604f80b
	// target_header_hash:  0x78d9f8d4d7af68e27ce06774748020b13f4df6d4f36dfd975e46614f8d941aad
	// data_commitment: 0xc1b21b6ad52a22080bfb9fa6f1bc7bdd53c73c9b1e41fed2c1d2b1eaebdcfb8e
	// trusted_block:  2202300
	// target_block: 2202310
	// validator_bit_map: 1267650595505862918627057991679
	height := uint64(2202300)
	inputPacked, err := abi.ABI.Pack(*BlobStreamInputsABI, "initializer", InitializerInput{
		Height:                    height,
		Header:                    [32]byte(common.Hex2BytesFixed("188b708bee180f43e3a252471754fd35283a6b09a6fd02f5b9130cc15604f80b", 32)),
		BlobstreamProgramVKeyHash: []byte("414456900754233403821469318749333346230962952863679230760144647782402486705"),
		BlobstreamProgramVKey:     vkey,
	})
	if err != nil {
		panic(err)
	}
	DynamicStateSlots := blobstream.StateKeysInitializer(height)
	action := []chain.Action{&actions.Deploy{
		ContractCode:            contractCode,
		InitializerFunctionName: "initializer",
		Input:                   inputPacked[4:],
		DynamicStateSlots:       DynamicStateSlots,
	}}

	_, tx, _, err := cli.GenerateTransaction(ctx, parser, action, factory)
	if err != nil {
		panic(err)
	}
	txID, err := cli.SubmitTx(ctx, tx.Bytes())
	if err != nil {
		panic(err)
	}
	fmt.Println(txID.String())
}
