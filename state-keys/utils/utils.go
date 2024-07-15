package utils

import (
	"fmt"
	"math/big"
)

func BigIntToBytes32(bigInt *big.Int) ([]byte, error) {
	nonceBytes := bigInt.Bytes()
	nonceBytesLen := len(nonceBytes)
	if nonceBytesLen > 32 {
		return nil, fmt.Errorf("bigInt is too large to fit in 32 bytes")
	}
	bigBuff := make([]byte, 32-nonceBytesLen)
	nonceBytes = append(bigBuff, nonceBytes...)
	return nonceBytes, nil
}
