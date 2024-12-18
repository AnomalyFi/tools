package main

import (
	"fmt"
	"os"

	"github.com/AnomalyFi/hypersdk/utils"
	"github.com/ava-labs/avalanchego/ids"
)

func main() {
	argsWithoutProg := os.Args[1:]

	nodeIDString := argsWithoutProg[0]
	// starting val server for node(NodeID-7Xhw2mDxuDS44j42TCB6U5579esbSt3Lg) at port: 33549
	id, err := ids.NodeIDFromString(nodeIDString)
	if err != nil {
		fmt.Printf("unable to parse nodeID: %s\n", err.Error())
		os.Exit(1)
	}
	port := utils.GetPortFromNodeID(id)
	fmt.Printf("%d", port)
}
