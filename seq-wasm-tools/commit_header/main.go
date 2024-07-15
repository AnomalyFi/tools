package main

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"

	"github.com/AnomalyFi/hypersdk/chain"
	"github.com/AnomalyFi/hypersdk/codec"
	"github.com/AnomalyFi/hypersdk/crypto/ed25519"
	"github.com/AnomalyFi/hypersdk/rpc"
	"github.com/AnomalyFi/nodekit-seq/actions"
	"github.com/AnomalyFi/nodekit-seq/auth"
	trpc "github.com/AnomalyFi/nodekit-seq/rpc"
	"github.com/AnomalyFi/tools/state-keys/blobstream"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
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

	contractAddress, _ := ids.FromString("KAPGGtG1HMyEwSE4mj16FrPYyiboiUayxNMtVzJ9jHaV8bBoP")
	contractAddress = chain.CreateActionID(contractAddress, 0)
	// contractCode, err := os.ReadFile("../blobstream.wasm")
	// if err != nil {
	// 	panic(err)
	// }
	// fsc, err := tcli.GetContract(ctx, contractAddress)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(bytes.Compare(contractCode, fsc))
	// trusted_header_hash: 0x188b708bee180f43e3a252471754fd35283a6b09a6fd02f5b9130cc15604f80b
	// target_header_hash:  0x78d9f8d4d7af68e27ce06774748020b13f4df6d4f36dfd975e46614f8d941aad
	// data_commitment: 0xc1b21b6ad52a22080bfb9fa6f1bc7bdd53c73c9b1e41fed2c1d2b1eaebdcfb8e
	// trusted_block:  2202300
	// target_block: 2202310
	// validator_bit_map: 1267650595505862918627057991679
	height := uint64(2202300)
	proof := "244e7b9370d3380deeff6340beadf03a8f584235d135e8e91095ac512fbd623b0951e0bdc6960f9115632aeb2715ac4bd39c16af03668159657270b8a0e01fe01c80f3578780b4fabeb831f51e9a2fc13dd966396d24b4bf5ed776df252d254b2b049e67c7f48eba905f332c5c6864ad5963ab20fc7ce27be4665d9c508b73c92274409ea6382d1b2e7db12f5c6274e71a085105d0feb0b1e7b76e5b0ecd6751054c8573b81e5381886b1a59d84626501942b2fd27c0380cc1a072fa0e89013b00dcdb5030d760ad2813e3cff52b5b63289f61c793f067bdfeef27df560e50df286d1d3c6d4966d1e4e1ba6cd811b6056ea80dbd650624f338addbce6f5c04ac1166db59d30b59e37812d63e219533f9a42ead2e5c633723987f1fd8241dd12826ecb0cc7ae72af7e90cceb5b2d01f8cafb16a2c603272f363e088a94ab5b74c05cd243c82d17e6a9aeceb7c5ff41600d2adeb83105e576d731eee02da58e95116dfa8f982f0c94c448028425a7e896082c68c6b2174856306192acadf1557fa0bf78b5b4ef07b5176d3ec45ee40304ea754a1b32951b08454c6b5e4d07e196703543dae9e5f5b2e9a08451bd01adf7cf3a6c35784c53f56fb0bb8ac368afc0a1e1b3a2bcf7aeeb32c4021cc7543c0bd2b4ad181c90b0442172db7d24dbe6be712542de410c22c1c2bf2593caeb7517f2c8986b8ea7463bf87848970b9dff7b400000007267459b6e97a3ee95dbac22ee24444bc22a433e52c20f7847bd253261a43965c1afda52230d7885632811c9608817c694efc563f07828394306a13795fb9559c1f7e71266f3ebe99dbec2f31118eab3c5959cdec5a9af80ae26c27896a2eee3412933ef1ca07e738604ad2520b1e47a85f5be64974dffb74ea27ff8e1ea85e2c12f39a1b85d95e7f0210ee4bed6c6bab5f9b730496b067f39336c6aad9097b3e1e7f976c0859eec6fb8e94220cae0ce8d01774057a6f1316c0c312453d6ec2bd255fe341459509a3867ee1bc6ecbb58de487141ed90e63b42a2164e8825092e92ffc75199966d6774a499c65f7e654d10e7a1dea2a056086e8d2e1e28cface1c2db9608c69874f762c7f3f1fd45dd4f719b77a3f56b55eab22e334ca9384911b115659d31d6a744e994ed42141a239e9b740137a4fd5dcd27d0d04b5fd0cf305000000012f6315f6219fc990b0accef92e45f47e7e26654a10cc5c71267384bd089309e8299de2e8cd06931596485f24160415f0ccf1da3e7430722629122102dfc21710"
	publicValues := []byte{24, 139, 112, 139, 238, 24, 15, 67, 227, 162, 82, 71, 23, 84, 253, 53, 40, 58, 107, 9, 166, 253, 2, 245, 185, 19, 12, 193, 86, 4, 248, 11, 120, 217, 248, 212, 215, 175, 104, 226, 124, 224, 103, 116, 116, 128, 32, 177, 63, 77, 246, 212, 243, 109, 253, 151, 94, 70, 97, 79, 141, 148, 26, 173, 193, 178, 27, 106, 213, 42, 34, 8, 11, 251, 159, 166, 241, 188, 123, 221, 83, 199, 60, 155, 30, 65, 254, 210, 193, 210, 177, 234, 235, 220, 251, 142, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 33, 154, 188, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 33, 154, 198, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 15, 255, 255, 254, 255, 255, 255, 255, 255, 255, 255, 255, 255}
	inputPacked, err := abi.ABI.Pack(*BlobStreamInputsABI, "commitHeaderRange", CommitHeaderRangeInput{
		Proof:        []byte(proof),
		PublicValues: publicValues,
	})
	if err != nil {
		panic(err)
	}

	DynamicStateSlots := blobstream.StateKeysInitializer(height)
	action := []chain.Action{&actions.Transact{
		ContractAddress:   contractAddress,
		FunctionName:      "commit_header_range",
		Input:             inputPacked[4:],
		DynamicStateSlots: DynamicStateSlots,
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
