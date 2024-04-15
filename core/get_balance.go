package core

func (bc *BlockChain) GetBalance(address string) float64 {
	balance := 0.0
	if address == "" {
		return balance
	}

	for _, block := range bc.Blocks {
		for _, transaction := range block.Transactions {
			if transaction.FromAddr == address {
				balance -= transaction.Amount
			}
			if transaction.ToAddr == address {
				balance += transaction.Amount
			}
		}
	}

	return balance
}
