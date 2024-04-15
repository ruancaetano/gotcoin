package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func CalculateBlockHash(index int, timestamp int64, prevHash string, transactionsSignatures string, nonce int) string {
	value := fmt.Sprintf("%d-%d-%s-%s-%d", index, timestamp, prevHash, transactionsSignatures, nonce)
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:])
}

func JoinBlockTransactionsSignatures(block *Block) string {
	var joinedTransactionsSignatures []string
	for _, t := range block.Transactions {
		joinedTransactionsSignatures = append(joinedTransactionsSignatures, t.Signature)
	}
	sort.Strings(joinedTransactionsSignatures)
	return strings.Join(joinedTransactionsSignatures, "")
}
