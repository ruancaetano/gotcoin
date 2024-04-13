package core

import (
	"fmt"
)

func SetupInitialBlocks(bc *BlockChain) {
	t0 := NewTransaction("", Wallet1.PublicKey, 100)
	t01 := NewTransaction("", Wallet2.PublicKey, 100)
	t02 := NewTransaction("", Wallet3.PublicKey, 100)
	bc.AddTransaction(t0)
	bc.AddTransaction(t01)
	bc.AddTransaction(t02)
	bc.MinePendingTransactions(Miner1.PublicKey)

	t := NewTransaction(Wallet1.PublicKey, Wallet2.PublicKey, 50)
	t.Sign(Wallet1.PrivateKey)

	t2 := NewTransaction(Wallet2.PublicKey, Wallet3.PublicKey, 75)
	t2.Sign(Wallet2.PrivateKey)

	t3 := NewTransaction(Wallet3.PublicKey, Wallet1.PublicKey, 95)
	t3.Sign(Wallet3.PrivateKey)

	bc.AddTransaction(t)
	bc.AddTransaction(t2)

	bc.MinePendingTransactions(Miner1.PublicKey)

	bc.AddTransaction(t3)
	bc.MinePendingTransactions(Miner1.PublicKey)

	c := NewBlockChainCalculator(bc)
	fmt.Printf("Wallet 1: %.2f\n", c.GetBalance(Wallet1.PublicKey))
	fmt.Printf("Wallet 2: %.2f\n", c.GetBalance(Wallet2.PublicKey))
	fmt.Printf("Wallet 3: %.2f\n", c.GetBalance(Wallet3.PublicKey))
	fmt.Printf("Wallet Miner: %.2f\n", c.GetBalance(Miner1.PublicKey))
}
