package main

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/flashbots/go-boost-utils/bls"
)

func main() {
	sk, pk, err := bls.GenerateNewKeypair()
	if err != nil {
		panic(err)
	}

	skBytes := sk.Bytes()
	pkBytes := pk.Bytes()

	skStr := hexutil.Encode(skBytes[:])
	pkStr := hexutil.Encode(pkBytes[:])

	fmt.Printf("sk: %s\n", skStr)
	fmt.Printf("pk: %s\n", pkStr)

	_, err = bls.SecretKeyFromBytes(skBytes[:])
	if err != nil {
		panic(err)
	}
}
