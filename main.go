package main

import (
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
)

func main() {
	bc := core.NewBlockChain()

	wallet, _ := core.NewWallet()
	wallet2, _ := core.NewWallet()
	wallet3, _ := core.NewWallet()
	minerWallter, _ := core.NewWallet()

	t0 := core.NewTransaction("", wallet.PublicKey, 100)
	t01 := core.NewTransaction("", wallet2.PublicKey, 100)
	t02 := core.NewTransaction("", wallet3.PublicKey, 100)
	bc.AddTransaction(t0)
	bc.AddTransaction(t01)
	bc.AddTransaction(t02)
	bc.MinePendingTransactions(minerWallter.PublicKey)

	t := core.NewTransaction(wallet.PublicKey, wallet2.PublicKey, 50)
	t.Sign(wallet.PrivateKey)

	t2 := core.NewTransaction(wallet2.PublicKey, wallet3.PublicKey, 75)
	t2.Sign(wallet2.PrivateKey)

	t3 := core.NewTransaction(wallet3.PublicKey, wallet.PublicKey, 95)
	t3.Sign(wallet3.PrivateKey)

	bc.AddTransaction(t)
	bc.AddTransaction(t2)

	bc.MinePendingTransactions(minerWallter.PublicKey)

	bc.AddTransaction(t3)
	bc.MinePendingTransactions(minerWallter.PublicKey)

	c := core.NewBlockChainCalculator(bc)
	fmt.Printf("Wallet 1: %.2f\n", c.GetBalance(wallet.PublicKey))
	fmt.Printf("Wallet 2: %.2f\n", c.GetBalance(wallet2.PublicKey))
	fmt.Printf("Wallet 3: %.2f\n", c.GetBalance(wallet3.PublicKey))
	fmt.Printf("Wallet Miner: %.2f\n", c.GetBalance(minerWallter.PublicKey))
}
