package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	hrpc "github.com/AnomalyFi/hypersdk/rpc"
	"github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
)

type Config struct {
	SEQConfig SEQConfig `json:"seq_config"`
	LogConfig LogConfig `json:"log_config"`
}

type SEQConfig struct {
	SEQNodeUri           string `json:"seq_node_uri"`
	ChainIDStr           string `json:"chain_id_str"`
	NetworkID            uint32 `json:"network_id"`
	Ed25519PrivateKeyHex string `json:"ed25519_private_key_hex"`
}

type LogConfig struct {
	LoggingLevel string `json:"log_level"`
	ToConsole    bool   `json:"to_console"`
	LogFile      string `json:"log_file"`
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

	var config Config
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
				filledChainID = chainID
			}
		}
	}
	// fetch network id
	hcli := hrpc.NewJSONRPCClient(uris[2])
	networkID, _, _, err := hcli.Network(ctx)
	if err != nil {
		panic(err)
	}
	// update config
	config.SEQConfig.ChainIDStr = filledChainID.String()
	config.SEQConfig.NetworkID = networkID
	config.SEQConfig.SEQNodeUri = uris[2]
	// write to config file
	d, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(os.WriteFile("config"+".json", d, 0644))
}
