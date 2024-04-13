package main

import (
	"context"
	"log"
	"os"

	"github.com/ruancaetano/gotcoin/infra"
)

func main() {
	ctx := context.Background()
	node, err := infra.InitNode(0)
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 1 {
		SetupPeerNode(ctx, node)
	} else {
		SetupGenesisNode(ctx, node)
	}
}
