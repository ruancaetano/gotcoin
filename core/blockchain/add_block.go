package blockchain

import (
	"github.com/ruancaetano/gotcoin/core"
)

func (bc *BlockChain) AddBlock(block *core.Block) {
	if !block.IsValid() {
		return
	}

	bc.Blocks = append(bc.Blocks, block)

	var newPendingTransactions []*core.Transaction
	for _, transaction := range bc.PendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}
	bc.PendingTransactions = newPendingTransactions
}
