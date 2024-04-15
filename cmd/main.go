package main

import (
	"context"
	"flag"
	"log"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/infra"
	"github.com/ruancaetano/gotcoin/network"
)

func main() {
	config := infra.Config{}

	config.Genesis = flag.Bool("genesis", false, "Launch as genesis node")
	flag.IntVar(&config.Port, "port", 0, "")
	flag.Parse()

	log.Println("Config: ", config)
	ctx := context.Background()
	node, err := network.InitNode(*config.Genesis, 0)
	if err != nil {
		log.Fatal(err)
	}

	if *config.Genesis {
		bc := core.NewBlockChain()
		eh := core.NewEventHandler(bc)
		network.SetupGenesisNode(ctx, node, bc, eh)
	} else {
		bc := core.NewEmptyBlockChain()
		eh := core.NewEventHandler(bc)
		network.SetupPeerNode(ctx, node, bc, eh)
	}

	// hang forever
	select {}
}
