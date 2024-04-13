package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Block struct {
	Index        int            `json:"index"`
	Timestamp    int64          `json:"timestamp"`
	Transactions []*Transaction `json:"transactions"`
	Nonce        int            `json:"nonce"`
	Hash         string         `json:"hash"`
	PrevHash     string         `json:"prev-hash"`
}

func NewBlock(index int, prevHash string, transactions []*Transaction) (*Block, error) {
	block := &Block{
		Index:        index,
		Timestamp:    time.Now().Unix(),
		Transactions: transactions,
		PrevHash:     prevHash,
	}

	hash, err := block.CalculateHash()
	if err != nil {
		return nil, err
	}
	block.Hash = hash

	return block, nil
}

func (block *Block) CalculateHash() (string, error) {
	value := fmt.Sprintf("%d-%d-%s-%s-%d", block.Index, block.Timestamp, block.PrevHash, block.joinTransactionsSignatures(), block.Nonce)
	hash := sha256.Sum256([]byte(value))
	return hex.EncodeToString(hash[:]), nil
}

func (block *Block) MineBlock(difficulty int) {
	difficultyMatch := strings.Repeat("0", difficulty)
	for {
		hash, _ := block.CalculateHash()

		if hash[:difficulty] == difficultyMatch {
			block.Hash = hash
			return
		}
		block.Nonce += 1
	}
}

func (block *Block) joinTransactionsSignatures() string {
	var joinedTransactionsSignatures []string
	for _, t := range block.Transactions {
		joinedTransactionsSignatures = append(joinedTransactionsSignatures, t.Signature)
	}
	sort.Strings(joinedTransactionsSignatures)
	return strings.Join(joinedTransactionsSignatures, "")
}

func (block *Block) IsValid() bool {
	calculatedHash, _ := block.CalculateHash()
	if block.Hash != calculatedHash {
		return false
	}

	for _, t := range block.Transactions {
		if !t.IsValid() {
			return false
		}
	}

	return true
}
