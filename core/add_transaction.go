package core

import (
	"fmt"
)

func (bc *BlockChain) AddTransaction(transaction *Transaction) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	for _, t := range bc.PendingTransactions {
		if t.Signature == transaction.Signature {
			return nil
		}
	}

	addressBalance := bc.GetBalance(transaction.FromAddr)
	if !transaction.IsCoinbase() && addressBalance < transaction.Amount {
		return fmt.Errorf("Not enough balance")
	}

	if !transaction.IsValid() {
		return fmt.Errorf("Transaction is not valid")
	}
	bc.PendingTransactions = append(bc.PendingTransactions, transaction)

	if len(bc.PendingTransactions) >= MinTransactionsToMine && !bc.mining {
		go bc.MinePendingTransactions(Miner1.PublicKey)
	}
	return nil
}
