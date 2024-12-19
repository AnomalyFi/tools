package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"

	"github.com/AnomalyFi/hypersdk/codec"
	"github.com/AnomalyFi/hypersdk/crypto/bls"
	"github.com/AnomalyFi/nodekit-seq/auth"
)

func main() {
	argsWithoutProg := os.Args[1:]
	blsKeyHex := argsWithoutProg[0]

	keyBytes, err := hex.DecodeString(strings.Trim(blsKeyHex, "0x"))
	if err != nil {
		fmt.Println("unable to decode bls key hex")
		os.Exit(1)
	}
	sk, err := bls.PrivateKeyFromBytes(keyBytes)
	if err != nil {
		fmt.Println("unable to parse bls key")
		os.Exit(1)
	}
	pk := bls.PublicFromPrivateKey(sk)
	addr := auth.NewBLSAddress(pk)

	benAddr, err := codec.AddressBech32("seq", addr)
	if err != nil {
		fmt.Println("unable to create ben address")
		os.Exit(1)
	}

	fmt.Printf("%s", benAddr)
}
