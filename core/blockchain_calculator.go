package core

type BlockChainCalculator struct {
	BlockChain *BlockChain
}

func NewBlockChainCalculator(bc *BlockChain) *BlockChainCalculator {
	return &BlockChainCalculator{
		BlockChain: bc,
	}
}

func (bce *BlockChainCalculator) GetBalance(address string) float64 {
	balance := 0.0
	if address == "" {
		return balance
	}

	for _, block := range bce.BlockChain.Blocks {
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
