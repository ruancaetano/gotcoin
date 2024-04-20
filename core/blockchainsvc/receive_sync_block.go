package blockchainsvc

import (
	"fmt"

	"go.uber.org/zap"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) ReceiveSyncBlock(block *core.Block, totalBlocksCount int, newDifficulty int) {
	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	for _, b := range bs.Blockchain.Blocks {
		if b.Hash == block.Hash {
			return
		}
	}

	if !block.IsValid() {
		bs.Logger.Debug("Block skipped because it is not valid", zap.String("hash", block.Hash))
		return
	}

	bs.Blockchain.Blocks = append(bs.Blockchain.Blocks, block)

	var newPendingTransactions []*core.Transaction
	for _, transaction := range bs.Blockchain.PendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}
	bs.Blockchain.PendingTransactions = newPendingTransactions
	bs.Blockchain.Difficulty = newDifficulty
	bs.verifySyncStatus(totalBlocksCount)
}

func (bs *blockchainServiceImpl) verifySyncStatus(totalBlocksCount int) {
	if len(bs.Blockchain.Blocks) < totalBlocksCount {
		bs.Logger.Debug(fmt.Sprintf("Sync percent complete: %.2f", float64(len(bs.Blockchain.Blocks))/float64(totalBlocksCount)*100.0))
		return
	}

	bs.Blockchain.SortBlocks()
	bs.Blockchain.Synced = true
	bs.Logger.Debug("Block chain sync completed")
	bs.Logger.Debug(fmt.Sprintf("Blocks count: %d", len(bs.Blockchain.Blocks)))
	bs.Logger.Debug(fmt.Sprintf("Block chain valid: %+v", bs.Blockchain.IsChainValid()))
	bs.Logger.Debug(fmt.Sprintf("Wallet 1: %.2f", bs.GetBalance(util.Wallet1.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet 2: %.2f", bs.GetBalance(util.Wallet2.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet 3: %.2f", bs.GetBalance(util.Wallet3.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet Miner: %.2f", bs.GetBalance(util.Miner1.PublicKey)))
	return
}
