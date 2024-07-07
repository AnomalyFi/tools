package blobstream

import (
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"strconv"

	"github.com/AnomalyFi/tools/state-keys/utils"
)

func StateKeysInitializer(height uint64) []string {
	hB := binary.BigEndian.AppendUint64(nil, height)
	slot := "slot" + strconv.Itoa(1) + hex.EncodeToString(hB)
	return []string{slot}
}

func StateKeysUpdateGenesisState(height uint64) []string {
	hB := binary.BigEndian.AppendUint64(nil, height)
	slot := "slot" + strconv.Itoa(1) + hex.EncodeToString(hB)
	return []string{slot}
}

func StateKeysCommitHeaderRange(height uint64, targetHeight uint64, nonce *big.Int) []string {
	hB := binary.BigEndian.AppendUint64(nil, height)
	slot1 := "slot" + strconv.Itoa(1) + hex.EncodeToString(hB)
	hB = binary.BigEndian.AppendUint64(nil, targetHeight)
	slot2 := "slot" + strconv.Itoa(1) + hex.EncodeToString(hB)
	nonceBytes, _ := utils.BigIntToBytes32(nonce)
	slot3 := "slot" + strconv.Itoa(2) + hex.EncodeToString(nonceBytes)
	return []string{slot1, slot2, slot3}
}

func StateKeysVerifyAttestation(nonce *big.Int) []string {
	nonceBytes, _ := utils.BigIntToBytes32(nonce)
	slot := "slot" + strconv.Itoa(2) + hex.EncodeToString(nonceBytes)
	return []string{slot}
}
