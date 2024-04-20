package blockchainsvc

import (
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) AddTransaction(transaction *core.Transaction) error {
	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	for _, t := range bs.Blockchain.PendingTransactions {
		if t.Signature == transaction.Signature && transaction.FromAddr != "" {
			return nil
		}
	}

	addressBalance := bs.GetBalance(transaction.FromAddr)
	if !transaction.IsCoinbase() && addressBalance < transaction.Amount {
		return fmt.Errorf("not enough balance")
	}

	if !transaction.IsValid() {
		return fmt.Errorf("transaction is not valid")
	}
	bs.Blockchain.PendingTransactions = append(bs.Blockchain.PendingTransactions, transaction)

	if len(bs.Blockchain.PendingTransactions) >= core.MinTransactionsToMine && !bs.Blockchain.Mining {
		go bs.MinePendingTransactions(util.Miner1.PublicKey)
	}
	return nil
}
