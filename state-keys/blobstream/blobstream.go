package blobstream

import (
	"encoding/binary"
	"math/big"

	"github.com/AnomalyFi/tools/state-keys/utils"
)

func StateKeysInitializer(height uint64) [][]byte {
	key := binary.BigEndian.AppendUint64(nil, height)
	slot := append(binary.BigEndian.AppendUint32(nil, 1), key...)
	return [][]byte{slot}
}

func StateKeysUpdateGenesisState(height uint64) [][]byte {
	key := binary.BigEndian.AppendUint64(nil, height)
	slot := append(binary.BigEndian.AppendUint32(nil, 1), key...)
	return [][]byte{slot}
}

func StateKeysCommitHeaderRange(height uint64, targetHeight uint64, nonce *big.Int) [][]byte {
	key1 := binary.BigEndian.AppendUint64(nil, height)
	slot1 := append(binary.BigEndian.AppendUint32(nil, 1), key1...)
	key2 := binary.BigEndian.AppendUint64(nil, targetHeight)
	slot2 := append(binary.BigEndian.AppendUint32(nil, 1), key2...)
	nonceBytes, _ := utils.BigIntToBytes32(nonce)
	slot3 := append(binary.BigEndian.AppendUint32(nil, 2), nonceBytes...)
	return [][]byte{slot1, slot2, slot3}
}

func StateKeysVerifyAttestation(nonce *big.Int) [][]byte {
	nonceBytes, _ := utils.BigIntToBytes32(nonce)
	slot := append(binary.BigEndian.AppendUint32(nil, 2), nonceBytes...)
	return [][]byte{slot}
}
