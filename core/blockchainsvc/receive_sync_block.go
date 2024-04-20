package blockchainsvc

import (
	"fmt"
	"log"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) ReceiveSyncBlock(block *core.Block, totalBlocksCount int) {
	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	for _, b := range bs.Blockchain.Blocks {
		if b.Hash == block.Hash {
			return
		}
	}

	if !block.IsValid() {
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

	bs.verifySyncStatus(totalBlocksCount)
}

func (bs *blockchainServiceImpl) verifySyncStatus(totalBlocksCount int) {
	if len(bs.Blockchain.Blocks) < totalBlocksCount {
		log.Println("Sync percent complete: ", float64(len(bs.Blockchain.Blocks))/float64(totalBlocksCount)*100.0)
		return
	}

	bs.Blockchain.SortBlocks()
	bs.Blockchain.Synced = true
	log.Println("Block chain sync completed")
	log.Println("Blocks count: ", len(bs.Blockchain.Blocks))
	log.Println("Block chain valid: ", bs.Blockchain.IsChainValid())
	fmt.Printf("Wallet 1: %.2f\n", bs.GetBalance(util.Wallet1.PublicKey))
	fmt.Printf("Wallet 2: %.2f\n", bs.GetBalance(util.Wallet2.PublicKey))
	fmt.Printf("Wallet 3: %.2f\n", bs.GetBalance(util.Wallet3.PublicKey))
	fmt.Printf("Wallet Miner: %.2f\n", bs.GetBalance(util.Miner1.PublicKey))
	return
}
