package core

import (
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
		Nonce:        0,
	}

	hash, err := block.CalculateHash()
	if err != nil {
		return nil, err
	}
	block.Hash = hash

	return block, nil
}

func NewGenesisBlock() *Block {
	block, _ := NewBlock(0, "", []*Transaction{})
	return block
}

func (block *Block) CalculateHash() (string, error) {
	return CalculateBlockHash(
		block.Index,
		block.Timestamp,
		block.PrevHash,
		JoinBlockTransactionsSignatures(block),
		block.Nonce), nil
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

func (block *Block) HasTransaction(transaction *Transaction) bool {
	for _, t := range block.Transactions {
		if t.Signature == transaction.Signature {
			return true
		}
	}
	return false
}
