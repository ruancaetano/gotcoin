package core

import (
	"fmt"
)

type BlockChain struct {
	pendingTransactions []*Transaction
	Blocks              []*Block `json:"blocks"`
	calculator          *BlockChainCalculator
}

func NewBlockChain() *BlockChain {
	bc := &BlockChain{}
	bc.calculator = NewBlockChainCalculator(bc)
	bc.Blocks = append(bc.Blocks, bc.CreateGenesisBlock())
	return bc
}

func (bc *BlockChain) CreateGenesisBlock() *Block {
	block, _ := NewBlock("", []*Transaction{})
	return block
}

func (bc *BlockChain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *BlockChain) AddTransaction(transaction *Transaction) error {
	addressBalance := bc.calculator.GetBalance(transaction.FromAddr)
	if !transaction.IsCoinbase() && addressBalance < transaction.Amount {
		return fmt.Errorf("Not enough balance")
	}

	if !transaction.IsValid() {
		return fmt.Errorf("Transaction is not valid")
	}
	bc.pendingTransactions = append(bc.pendingTransactions, transaction)
	return nil
}

func (bc *BlockChain) MinePendingTransactions(mineRewardAddress string) {
	if len(bc.pendingTransactions) == 0 {
		return
	}
	block, _ := NewBlock(bc.GetLastBlock().Hash, bc.pendingTransactions)
	block.MineBlock(MineDifficulty)
	fmt.Println("Block mined: ", block.Hash)
	bc.Blocks = append(bc.Blocks, block)

	bc.pendingTransactions = []*Transaction{
		NewTransaction("", mineRewardAddress, MineReward),
	}
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
