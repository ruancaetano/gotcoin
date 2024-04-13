package main

import (
	"context"
	"flag"
	"log"

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
		network.SetupGenesisNode(ctx, node)
	} else {
		network.SetupPeerNode(ctx, node)
	}
}
