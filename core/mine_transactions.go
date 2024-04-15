package core

import (
	"fmt"
)

func (bc *BlockChain) MinePendingTransactions(mineRewardAddress string) {
	if len(bc.PendingTransactions) == 0 {
		return
	}
	lastBlock := bc.GetLastBlock()
	block, _ := NewBlock(lastBlock.Index+1, lastBlock.Hash, bc.PendingTransactions)
	block.MineBlock(bc.Difficulty)
	fmt.Println("Block mined: ", block.Hash)
	bc.Blocks = append(bc.Blocks, block)

	bc.PendingTransactions = []*Transaction{
		NewTransaction("", mineRewardAddress, MineReward),
	}
}
