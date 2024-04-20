package blockchainsvc

import (
	"time"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/events"
)

func (bs *blockchainServiceImpl) MinePendingTransactions(mineRewardAddress string) {
	if len(bs.Blockchain.PendingTransactions) == 0 {
		return
	}
	lastBlock := bs.Blockchain.GetLastBlock()

	bs.Blockchain.Mining = true
	transactionsToMine := append(bs.Blockchain.PendingTransactions, core.NewTransaction("", mineRewardAddress, core.MineReward))
	block, _ := core.NewBlock(lastBlock.Index+1, lastBlock.Hash, transactionsToMine)
	startTime := time.Now()
	block.MineBlock(bs.Blockchain.Difficulty)
	endTime := time.Now()
	bs.Logger.Debug("Block mined: " + block.Hash)

	err := bs.AddBlock(block, bs.calculateNewDifficulty(endTime.Sub(startTime)))
	if err != nil {
		bs.Logger.Debug("Block not added: " + err.Error())
		return
	}
	bs.Node.SendBroadcastEvent(events.SendNewBlockEvent(block, &bs.Blockchain.Difficulty))
	bs.Blockchain.Mining = false
}

func (bs *blockchainServiceImpl) calculateNewDifficulty(timeToMineBlock time.Duration) int {
	//if timeToMineBlock < core.ExpectedTimeToMine {
	//	return bs.Blockchain.Difficulty + 1
	//}
	//return bs.Blockchain.Difficulty - 1
	return core.InitialMineDifficulty
}
