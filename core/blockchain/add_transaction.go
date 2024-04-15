package blockchain

import (
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
)

func (bc *BlockChain) AddTransaction(transaction *core.Transaction) error {
	addressBalance := bc.GetBalance(transaction.FromAddr)
	if !transaction.IsCoinbase() && addressBalance < transaction.Amount {
		return fmt.Errorf("Not enough balance")
	}

	if !transaction.IsValid() {
		return fmt.Errorf("Transaction is not valid")
	}
	bc.PendingTransactions = append(bc.PendingTransactions, transaction)
	return nil
}
