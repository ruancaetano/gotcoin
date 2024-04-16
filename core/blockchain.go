package core

import (
	"sort"
	"sync"
)

type BlockChain struct {
	Synced              bool           `json:"synced"`
	PendingTransactions []*Transaction `json:"pendingTransactions"`
	Blocks              []*Block       `json:"blocks"`
	Difficulty          int            `json:"difficulty"`

	mining bool
	lock   sync.Mutex
}

func NewBlockChain() *BlockChain {
	bc := &BlockChain{
		Difficulty: InitialMineDifficulty,
		lock:       sync.Mutex{},
	}
	bc.Blocks = append(bc.Blocks, bc.CreateGenesisBlock())
	return bc
}

func NewEmptyBlockChain() *BlockChain {
	bc := &BlockChain{
		Difficulty: InitialMineDifficulty,
		lock:       sync.Mutex{},
	}
	return bc
}

func (bc *BlockChain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *BlockChain) IsChainValid() bool {
	for idx, block := range bc.Blocks {
		if idx != 0 && block.PrevHash != bc.Blocks[idx-1].Hash {
			return false
		}

		if !block.IsValid() {
			return false
		}
	}

	return true
}

func (bc *BlockChain) SortBlocks() {
	// TODO: explore insertion sort
	less := func(i, j int) bool {
		return bc.Blocks[i].Index < bc.Blocks[j].Index
	}
	sort.Slice(bc.Blocks, less)
}
