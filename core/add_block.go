package core

func (bc *BlockChain) AddBlock(block *Block) {
	if !block.IsValid() {
		return
	}

	bc.Blocks = append(bc.Blocks, block)

	var newPendingTransactions []*Transaction
	for _, transaction := range bc.PendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}
	bc.PendingTransactions = newPendingTransactions
}
