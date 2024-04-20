package blockchainsvc

import (
	"log"
	"strings"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) AddBlock(newBlock *core.Block, newDifficulty int) {
	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	lastBlock := bs.Blockchain.GetLastBlock()

	duplicatedBlock := newBlock.PrevHash == lastBlock.PrevHash
	if duplicatedBlock {
		log.Println("Block skipped because it is duplicated: ", newBlock.Hash)
		return
	}

	// verify if all transaction in new block are in pending transactions
	newPendingTransactions := util.CloneSlice(bs.Blockchain.PendingTransactions)
	for _, transaction := range newBlock.Transactions {
		found := false
		for _, pendingTransaction := range bs.Blockchain.PendingTransactions {
			if core.CompareTransactionFunc(pendingTransaction, transaction) {
				newPendingTransactions = util.RemoveFromSlice(bs.Blockchain.PendingTransactions, pendingTransaction, core.CompareTransactionFunc)
				found = true
				break
			}
		}

		if !found {
			log.Println("Block skipped because transaction not found in pending transactions: ", transaction.Signature)
			return
		}
	}

	if !bs.isNewValidBlock(newBlock) {
		log.Println("Block skipped because it is not valid: ", newBlock.Hash)
		return
	}

	bs.Blockchain.PendingTransactions = newPendingTransactions
	bs.Blockchain.Blocks = append(bs.Blockchain.Blocks, newBlock)
	bs.Blockchain.Difficulty = newDifficulty
}

func (bs *blockchainServiceImpl) isNewValidBlock(newBlock *core.Block) bool {
	lastBlock := bs.Blockchain.GetLastBlock()

	expectedNewBlockHash := core.CalculateBlockHash(
		newBlock.Index,
		newBlock.Timestamp,
		lastBlock.PrevHash,
		core.JoinBlockTransactionsSignatures(newBlock),
		newBlock.Nonce)

	difficultyMatch := strings.Repeat("0", bs.Blockchain.Difficulty)
	hashMatchesWithDifficulty := strings.HasPrefix(expectedNewBlockHash, difficultyMatch)

	validTimestamp := newBlock.Timestamp > lastBlock.Timestamp

	return newBlock.Hash == expectedNewBlockHash &&
		hashMatchesWithDifficulty &&
		validTimestamp
}
