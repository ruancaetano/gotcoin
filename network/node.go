package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	libhost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/core"
)

type Node struct {
	ID           peer.ID
	Addr         string
	Host         libhost.Host
	EventHandler *core.EventHandler
	Streams      map[peer.ID]network.Stream
}

func NewGenesisNode(ctx context.Context, host libhost.Host, bc *core.BlockChain, eh *core.EventHandler) *Node {
	core.SetupInitialBlocks(bc)

	_, err := NewKDHT(ctx, host, nil)
	if err != nil {
		log.Fatal(err)
	}

	node := &Node{
		ID:           host.ID(),
		Addr:         host.Addrs()[0].String(),
		Host:         host,
		EventHandler: eh,
		Streams:      map[peer.ID]network.Stream{},
	}

	host.SetStreamHandler("/p2p/1.0.0", node.HandleNewStream)
	return node
}

func NewNode(ctx context.Context, host libhost.Host, _ *core.BlockChain, eh *core.EventHandler) *Node {
	genesisAddr, err := multiaddr.NewMultiaddr(GenesisNodeAddr)
	if err != nil {
		log.Fatal("failed to parse multiaddr:", err)
	}
	genesisPeerInfo, _ := peer.AddrInfoFromP2pAddr(genesisAddr)

	kdht, err := NewKDHT(ctx, host, []multiaddr.Multiaddr{genesisAddr})
	if err != nil {
		log.Fatal(err)
	}

	s, err := ConnectToHost(ctx, host, genesisAddr)
	if err != nil {
		log.Fatal(err)
	}

	node := &Node{
		ID:           host.ID(),
		Addr:         host.Addrs()[0].String(),
		Host:         host,
		EventHandler: eh,
		Streams: map[peer.ID]network.Stream{
			genesisPeerInfo.ID: s,
		},
	}
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go node.ReadEvent(rw)

	node.setupDiscovery(ctx, host, kdht)

	if err = node.InitSync(genesisPeerInfo.ID); err != nil {
		log.Fatal(err)
	}

	host.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
			fmt.Printf("Peer connected: %s\n", c.RemotePeer())
		},
		DisconnectedF: func(n network.Network, c network.Conn) {
			fmt.Printf("Peer disconnected: %s\n", c.RemotePeer())
			delete(node.Streams, c.RemotePeer())
		},
	})

	return node
}

func (n *Node) HandleNewStream(s network.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go n.ReadEvent(rw)
}

func (n *Node) InitSync(peerID peer.ID) error {
	s := n.Streams[peerID]
	if s == nil {
		return fmt.Errorf("no stream to peer %s", peerID)
	}

	return n.SendEvent(s, core.RequestBlockChainSyncEvent())
}

func (n *Node) setupDiscovery(ctx context.Context, host libhost.Host, dht *dht.IpfsDHT) {
	discoveryAddrChan := make(chan string)
	go Discover(ctx, discoveryAddrChan, host, dht, "gotcoin")
	go func(ctx context.Context, host libhost.Host, discoveryAddrChan chan string) {
		fmt.Println("listening for discovery addresses")
		for {
			foundAddr := <-discoveryAddrChan

			fmt.Println("foundAddr:", foundAddr)

			addr, err := multiaddr.NewMultiaddr(foundAddr)
			if err != nil {
				log.Fatal("failed to parse multiaddr:", err)
			}

			s, err := ConnectToHost(ctx, host, addr)
			if err != nil {
				log.Println(err)
			}
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			go n.ReadEvent(rw)

			peerInfo, _ := peer.AddrInfoFromP2pAddr(addr)
			n.Streams[peerInfo.ID] = s
		}
	}(ctx, host, discoveryAddrChan)
}

func (n *Node) ReadEvent(rw *bufio.ReadWriter) {
	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			return
		}

		if str == "" {
			return
		}
		if str != "\n" {
			eventData := core.EventData{}
			if err = json.Unmarshal([]byte(str), &eventData); err != nil {
				log.Println("Failed to unmarshal event data")
				continue
			}

			go n.EventHandler.HandleEvent(rw, eventData)
		}
	}
}

func (n *Node) SendEvent(s network.Stream, event core.EventData) error {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = rw.WriteString(fmt.Sprintf("%s\n", data))
	if err != nil {
		return err
	}

	err = rw.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (n *Node) BroadcastEvent(event core.EventData) {
	for _, s := range n.Streams {
		err := n.SendEvent(s, event)
		if err != nil {
			log.Println("Failed to send event", err)
		}
	}
}
