package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	hrpc "github.com/AnomalyFi/hypersdk/rpc"
	"github.com/ava-labs/avalanche-network-runner/client"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils/logging"
)

func main() {
	ctx := context.Background()

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

	fmt.Println("chain id", filledChainID)
	chainID := string([]byte("nkit"))
	cli0 := hrpc.NewJSONRPCClient(uris[0])
	cli1 := hrpc.NewJSONRPCClient(uris[1])
	cli2 := hrpc.NewJSONRPCClient(uris[2])
	cli3 := hrpc.NewJSONRPCClient(uris[3])
	cli4 := hrpc.NewJSONRPCClient(uris[4])

	for {

		go func() {
			price, err := cli0.NameSpacePrice(ctx, chainID)
			if err != nil {
				fmt.Println("error getting namespace price", err)
			}
			fmt.Println("namespace price 0", price)
		}()
		go func() {
			price, err := cli1.NameSpacePrice(ctx, chainID)
			if err != nil {
				fmt.Println("error getting namespace price", err)
			}
			fmt.Println("namespace price 1", price)
		}()
		go func() {
			price, err := cli2.NameSpacePrice(ctx, chainID)
			if err != nil {
				fmt.Println("error getting namespace price", err)
			}
			fmt.Println("namespace price 2", price)
		}()
		go func() {
			price, err := cli3.NameSpacePrice(ctx, chainID)
			if err != nil {
				fmt.Println("error getting namespace price", err)
			}
			fmt.Println("namespace price 3", price)
		}()
		go func() {
			price, err := cli4.NameSpacePrice(ctx, chainID)
			if err != nil {
				fmt.Println("error getting namespace price", err)
			}
			fmt.Println("namespace price 4", price)
		}()
		time.Sleep(500 * time.Millisecond)
	}
}
