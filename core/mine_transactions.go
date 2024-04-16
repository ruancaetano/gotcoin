package core

import (
	"fmt"
)

func (bc *BlockChain) MinePendingTransactions(mineRewardAddress string) {
	if len(bc.PendingTransactions) == 0 {
		return
	}
	lastBlock := bc.GetLastBlock()

	transactionsToMine := append(bc.PendingTransactions, NewTransaction("", mineRewardAddress, MineReward))
	block, _ := NewBlock(lastBlock.Index+1, lastBlock.Hash, transactionsToMine)
	block.MineBlock(bc.Difficulty)
	fmt.Println("Block mined: ", block.Hash)

	bc.lock.Lock()
	defer bc.lock.Unlock()

	bc.Blocks = append(bc.Blocks, block)

	var newPendingTransactions []*Transaction
	for _, transaction := range bc.PendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}

	bc.PendingTransactions = newPendingTransactions
}
