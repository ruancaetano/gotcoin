package main

import (
	"context"
	"flag"
	"log"

	"github.com/ruancaetano/gotcoin/adapters"
	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/blockchainsvc"
	"github.com/ruancaetano/gotcoin/core/protocols"
	"github.com/ruancaetano/gotcoin/infra"
	"github.com/ruancaetano/gotcoin/infra/network"
)

func main() {
	config := infra.Config{}

	config.Genesis = flag.Bool("genesis", false, "Launch as genesis node")
	flag.IntVar(&config.Port, "port", 0, "")
	flag.Parse()

	ctx := context.Background()
	host, err := network.InitHost(*config.Genesis, 0)
	if err != nil {
		log.Fatal(err)
	}

	var node protocols.Node
	if *config.Genesis {
		node = adapters.NewNodeAdapterAsGenesis(ctx, host)

		bc := core.NewBlockChain()
		bs := blockchainsvc.NewBlockchainServiceImpl(bc, node)
		eh := adapters.NewEventHandlerAdapter(bs)

		node.Setup(ctx, eh)
		bs.SetupInitialBlocks()
	} else {
		node = adapters.NewNodeAdapter(ctx, host)

		bc := core.NewEmptyBlockChain()
		bs := blockchainsvc.NewBlockchainServiceImpl(bc, node)
		eh := adapters.NewEventHandlerAdapter(bs)

		node.Setup(ctx, eh)
	}

	// hang forever
	select {}
}
