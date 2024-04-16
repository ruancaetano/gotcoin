package core

func (bc *BlockChain) AddBlock(block *Block) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	for _, b := range bc.Blocks {
		if b.Hash == block.Hash {
			return
		}
	}

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
