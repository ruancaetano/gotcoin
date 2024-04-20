package protocols

import (
	"github.com/ruancaetano/gotcoin/core"
)

type BlockchainService interface {
	AddTransaction(transaction *core.Transaction) error
	GetBalance(address string) float64
	AddBlock(newBlock *core.Block, newDifficulty int) error
	MinePendingTransactions(mineRewardAddress string)
	ReceiveSyncBlock(block *core.Block, totalBlocksCount int, newDifficulty int)
	SendSyncBlocks(targetPeerID string)
	GetBlocks() []*core.Block
	SetupInitialBlocks()
}
