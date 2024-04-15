package core

import (
	"fmt"
	"sort"
)

type BlockChain struct {
	Synced              bool
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

func NewEmptyBlockChain() *BlockChain {
	bc := &BlockChain{}
	bc.calculator = NewBlockChainCalculator(bc)
	return bc
}

func (bc *BlockChain) CreateGenesisBlock() *Block {
	block, _ := NewBlock(0, "", []*Transaction{})
	return block
}

func (bc *BlockChain) GetLastBlock() *Block {
	return bc.Blocks[len(bc.Blocks)-1]
}

func (bc *BlockChain) AddBlock(block *Block) {
	if !block.IsValid() {
		return
	}

	bc.Blocks = append(bc.Blocks, block)

	var newPendingTransactions []*Transaction
	for _, transaction := range bc.pendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}
	bc.pendingTransactions = newPendingTransactions
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
	lastBlock := bc.GetLastBlock()
	block, _ := NewBlock(lastBlock.Index+1, lastBlock.Hash, bc.pendingTransactions)
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

func (bc *BlockChain) SortBlocks() {
	// TODO: explore insertion sort
	less := func(i, j int) bool {
		return bc.Blocks[i].Index < bc.Blocks[j].Index
	}
	sort.Slice(bc.Blocks, less)
}
