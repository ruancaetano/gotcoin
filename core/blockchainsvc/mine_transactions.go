package blockchainsvc

import (
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
)

func (bs *blockchainServiceImpl) MinePendingTransactions(mineRewardAddress string) {
	if len(bs.Blockchain.PendingTransactions) == 0 {
		return
	}
	lastBlock := bs.Blockchain.GetLastBlock()

	transactionsToMine := append(bs.Blockchain.PendingTransactions, core.NewTransaction("", mineRewardAddress, core.MineReward))
	block, _ := core.NewBlock(lastBlock.Index+1, lastBlock.Hash, transactionsToMine)
	block.MineBlock(bs.Blockchain.Difficulty)
	fmt.Println("Block mined: ", block.Hash)

	bs.Blockchain.Mutex.Lock()
	defer bs.Blockchain.Mutex.Unlock()

	bs.Blockchain.Blocks = append(bs.Blockchain.Blocks, block)

	var newPendingTransactions []*core.Transaction
	for _, transaction := range bs.Blockchain.PendingTransactions {
		if !block.HasTransaction(transaction) {
			newPendingTransactions = append(newPendingTransactions, transaction)
		}
	}

	bs.Blockchain.PendingTransactions = newPendingTransactions
}
