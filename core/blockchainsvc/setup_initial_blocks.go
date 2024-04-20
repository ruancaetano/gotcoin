package blockchainsvc

import (
	"fmt"
	"time"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) SetupInitialBlocks() {
	time.Sleep(1 * time.Second)
	t0 := core.NewTransaction("", util.Wallet1.PublicKey, 100)
	t01 := core.NewTransaction("", util.Wallet2.PublicKey, 100)
	t02 := core.NewTransaction("", util.Wallet3.PublicKey, 100)
	bs.AddTransaction(t0)
	bs.AddTransaction(t01)
	bs.AddTransaction(t02)
	bs.MinePendingTransactions(util.Miner1.PublicKey)
	//t := core.NewTransaction(util.Wallet1.PublicKey, util.Wallet2.PublicKey, 50)
	//t.Sign(util.Wallet1.PrivateKey)
	//
	//t2 := core.NewTransaction(util.Wallet2.PublicKey, util.Wallet3.PublicKey, 75)
	//t2.Sign(util.Wallet2.PrivateKey)
	//
	//t3 := core.NewTransaction(util.Wallet3.PublicKey, util.Wallet1.PublicKey, 95)
	//t3.Sign(util.Wallet3.PrivateKey)
	//
	//bs.AddTransaction(t)
	//bs.AddTransaction(t2)
	//
	//bs.MinePendingTransactions(util.Miner1.PublicKey)
	//
	//bs.AddTransaction(t3)
	//bs.MinePendingTransactions(util.Miner1.PublicKey)

	bs.Logger.Debug(fmt.Sprintf("Wallet 1: %.2f", bs.GetBalance(util.Wallet1.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet 2: %.2f", bs.GetBalance(util.Wallet2.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet 3: %.2f", bs.GetBalance(util.Wallet3.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Wallet Miner: %.2f", bs.GetBalance(util.Miner1.PublicKey)))
	bs.Logger.Debug(fmt.Sprintf("Blockchain is valid: %t", bs.Blockchain.IsChainValid()))
	bs.Logger.Debug("Initial blocks setup completed")
}
