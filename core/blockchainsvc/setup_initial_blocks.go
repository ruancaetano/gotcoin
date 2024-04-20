package blockchainsvc

import (
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/util"
)

func (bs *blockchainServiceImpl) SetupInitialBlocks() {
	t0 := core.NewTransaction("", util.Wallet1.PublicKey, 100)
	t01 := core.NewTransaction("", util.Wallet2.PublicKey, 100)
	t02 := core.NewTransaction("", util.Wallet3.PublicKey, 100)
	err := bs.AddTransaction(t0)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = bs.AddTransaction(t01)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = bs.AddTransaction(t02)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	fmt.Printf("Wallet 1: %.2f\n", bs.GetBalance(util.Wallet1.PublicKey))
	fmt.Printf("Wallet 2: %.2f\n", bs.GetBalance(util.Wallet2.PublicKey))
	fmt.Printf("Wallet 3: %.2f\n", bs.GetBalance(util.Wallet3.PublicKey))
	fmt.Printf("Wallet Miner: %.2f\n", bs.GetBalance(util.Miner1.PublicKey))
}
