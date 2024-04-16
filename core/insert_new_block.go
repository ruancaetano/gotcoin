package core

import (
	"log"
	"strings"

	"github.com/ruancaetano/gotcoin/util"
)

func (bc *BlockChain) InsertNewBlock(newBlock *Block, newDifficulty int) {
	bc.lock.Lock()
	defer bc.lock.Unlock()

	lastBlock := bc.GetLastBlock()

	duplicatedBlock := newBlock.PrevHash == lastBlock.PrevHash
	if duplicatedBlock {
		log.Println("Block skipped because it is duplicated: ", newBlock.Hash)
		return
	}

	// verify if all transaction in new block are in pending transactions
	newPendingTransactions := util.CloneSlice(bc.PendingTransactions)
	for _, transaction := range newBlock.Transactions {
		found := false
		for _, pendingTransaction := range bc.PendingTransactions {
			if CompareTransactionFunc(pendingTransaction, transaction) {
				newPendingTransactions = util.RemoveFromSlice(bc.PendingTransactions, pendingTransaction, CompareTransactionFunc)
				found = true
				break
			}
		}

		if !found {
			log.Println("Block skipped because transaction not found in pending transactions: ", transaction.Signature)
			return
		}
	}

	if !bc.isNewValidBlock(newBlock) {
		log.Println("Block skipped because it is not valid: ", newBlock.Hash)
		return
	}

	bc.PendingTransactions = newPendingTransactions
	bc.Blocks = append(bc.Blocks, newBlock)
	bc.Difficulty = newDifficulty
}

func (bc *BlockChain) isNewValidBlock(newBlock *Block) bool {
	lastBlock := bc.GetLastBlock()

	expectedNewBlockHash := CalculateBlockHash(
		newBlock.Index,
		newBlock.Timestamp,
		lastBlock.PrevHash,
		JoinBlockTransactionsSignatures(newBlock),
		newBlock.Nonce)

	difficultyMatch := strings.Repeat("0", bc.Difficulty)
	hashMatchesWithDifficulty := strings.HasPrefix(expectedNewBlockHash, difficultyMatch)

	validTimestamp := newBlock.Timestamp > lastBlock.Timestamp

	return newBlock.Hash == expectedNewBlockHash &&
		hashMatchesWithDifficulty &&
		validTimestamp
}
