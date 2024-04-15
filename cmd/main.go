package main

import (
	"context"
	"flag"
	"fmt"
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
	host, err := network.InitHost(*config.Genesis, 0)
	if err != nil {
		log.Fatal(err)
	}

	var node *network.Node
	if *config.Genesis {
		bc := core.NewBlockChain()
		eh := core.NewEventHandler(bc)

		node = network.NewGenesisNode(ctx, host, bc, eh)
	} else {
		bc := core.NewEmptyBlockChain()
		eh := core.NewEventHandler(bc)
		node = network.NewNode(ctx, host, bc, eh)
	}

	fmt.Println("Node ID: ", node.ID)
	fmt.Println("Node Addr: ", node.Addr)
	// hang forever
	select {}
}
