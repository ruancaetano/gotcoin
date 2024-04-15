package blockchain

import (
	"github.com/ruancaetano/gotcoin/core"
)

func (bc *BlockChain) CreateGenesisBlock() *core.Block {
	block, _ := core.NewBlock(0, "", []*core.Transaction{})
	return block
}
