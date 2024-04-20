package blockchainsvc

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) AddBlock(newBlock *core.Block, newDifficulty int) error {
	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	lastBlock := bs.Blockchain.GetLastBlock()

	duplicatedBlock := newBlock.PrevHash == lastBlock.PrevHash
	if duplicatedBlock {
		bs.Logger.Debug(fmt.Sprintf("block skipped because it is duplicated: %s", newBlock.Hash))
		return errors.New("block is duplicated")
	}

	// verify if all transaction in new block are in pending transactions
	newPendingTransactions := util.CloneSlice(bs.Blockchain.PendingTransactions)
	for _, transaction := range newBlock.Transactions {
		// skip reward transaction
		if transaction.FromAddr == "" {
			newPendingTransactions = util.RemoveFromSlice(newPendingTransactions, transaction, core.CompareTransactionFunc)
			continue
		}

		found := false
		for _, pendingTransaction := range bs.Blockchain.PendingTransactions {
			if core.CompareTransactionFunc(pendingTransaction, transaction) {
				newPendingTransactions = util.RemoveFromSlice(newPendingTransactions, pendingTransaction, core.CompareTransactionFunc)
				found = true
				break
			}
		}

		if !found {
			bs.Logger.Debug(fmt.Sprintf("Block skipped because transaction not found in pending transactions: %s", transaction.Signature))
			return errors.New("transaction not found in pending transactions")
		}
	}

	if !bs.isNewValidBlock(newBlock) {
		bs.Logger.Debug(fmt.Sprintf("block skipped because it is not valid: %s", newBlock.Hash))
		return errors.New("block is not valid")
	}

	bs.Blockchain.PendingTransactions = newPendingTransactions
	bs.Blockchain.Blocks = append(bs.Blockchain.Blocks, newBlock)
	bs.Blockchain.Difficulty = newDifficulty
	return nil
}

func (bs *blockchainServiceImpl) isNewValidBlock(newBlock *core.Block) bool {
	lastBlock := bs.Blockchain.GetLastBlock()

	difficultyMatch := strings.Repeat("0", bs.Blockchain.Difficulty)
	hashMatchesWithDifficulty := strings.HasPrefix(newBlock.Hash, difficultyMatch)
	validTimestamp := newBlock.Timestamp > lastBlock.Timestamp

	return hashMatchesWithDifficulty &&
		validTimestamp
}
