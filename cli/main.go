package main

import (
	"context"
	"fmt"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/network"
)

func main() {
	ctx := context.Background()
	node, err := network.InitNode(false, 0)
	if err != nil {
		panic(err)
	}

	bc := core.NewEmptyBlockChain()
	eh := core.NewEventHandler(bc)

	network.SetupPeerNode(ctx, node, bc, eh)
	for !bc.Synced {
	}
	fmt.Println("Blockchain synced")

	// hang forever
	select {}
}
