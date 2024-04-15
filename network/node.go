package network

import (
	"bufio"
	"context"
	"fmt"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/core"
	"github.com/ruancaetano/gotcoin/core/blockchain"
)

func SetupGenesisNode(ctx context.Context, node host.Host, bc *blockchain.BlockChain, eh *core.EventHandler) {
	_, err := NewKDHT(ctx, node, nil)
	if err != nil {
		log.Fatal(err)
	}

	node.SetStreamHandler("/p2p/1.0.0", func(s network.Stream) { HandleNewStream(s, eh) })
	core.SetupInitialBlocks(bc)

	log.Println("listening for connections")
}

func SetupPeerNode(ctx context.Context, node host.Host, _ *blockchain.BlockChain, eh *core.EventHandler) (*network.Stream, *bufio.ReadWriter) {
	node.SetStreamHandler("/p2p/1.0.0", func(s network.Stream) { HandleNewStream(s, eh) })

	addr, err := multiaddr.NewMultiaddr(GenesisNodeAddr)
	if err != nil {
		log.Fatal("failed to parse multiaddr:", err)
	}

	dht, err := NewKDHT(ctx, node, []multiaddr.Multiaddr{addr})
	if err != nil {
		log.Fatal(err)
	}
	setupDiscovery(ctx, node, dht, eh)

	s, rw, err := ConnectToNode(ctx, node, addr, eh)

	err = SendEvent(rw, core.RequestBlockChainSyncEvent())
	if err != nil {
		log.Fatal(err)
	}

	return s, rw
}

func setupDiscovery(ctx context.Context, node host.Host, dht *dht.IpfsDHT, eh *core.EventHandler) {
	discoveryAddrChan := make(chan string)
	go Discover(ctx, discoveryAddrChan, node, dht, "gotcoin")
	go func(ctx context.Context, node host.Host, discoveryAddrChan chan string) {
		fmt.Println("listening for discovery addresses")
		for {
			foundAddr := <-discoveryAddrChan

			fmt.Println("foundAddr:", foundAddr)

			addr, err := multiaddr.NewMultiaddr(foundAddr)
			if err != nil {
				log.Fatal("failed to parse multiaddr:", err)
			}

			_, _, err = ConnectToNode(ctx, node, addr, eh)
			if err != nil {
				log.Println(err)
			}
		}
	}(ctx, node, discoveryAddrChan)
}
