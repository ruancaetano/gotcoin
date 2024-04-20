package blockchainsvc

func (bs *blockchainServiceImpl) GetBalance(address string) float64 {
	balance := 0.0
	if address == "" {
		return balance
	}

	for _, block := range bs.Blockchain.Blocks {
		for _, transaction := range block.Transactions {
			if transaction.FromAddr == address {
				balance -= transaction.Amount
			}
			if transaction.ToAddr == address {
				balance += transaction.Amount
			}
		}
	}

	for _, transaction := range bs.Blockchain.PendingTransactions {
		if transaction.FromAddr == address {
			balance -= transaction.Amount
		}
		if transaction.ToAddr == address {
			balance += transaction.Amount
		}
	}

	return balance
}
