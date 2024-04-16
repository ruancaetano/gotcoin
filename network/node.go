package network

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	libhost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/patrickmn/go-cache"

	"github.com/ruancaetano/gotcoin/core"
)

type Node struct {
	ID           peer.ID
	Addr         string
	Host         libhost.Host
	EventHandler *core.EventHandler
	Streams      map[string]network.Stream
	eventsCache  *cache.Cache
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
		Streams:      map[string]network.Stream{},
		eventsCache:  cache.New(5*time.Minute, 10*time.Minute),
	}

	host.SetStreamHandler("/p2p/1.0.0", node.HandleNewStream)
	node.registerCallbacks()

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
		Streams: map[string]network.Stream{
			genesisPeerInfo.ID.String(): s,
		},
		eventsCache: cache.New(5*time.Minute, 10*time.Minute),
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	host.SetStreamHandler("/p2p/1.0.0", node.HandleNewStream)
	node.setupDiscovery(ctx, host, kdht)
	node.registerCallbacks()
	go node.ReadEvent(rw)

	if err = node.InitSync(genesisPeerInfo.ID); err != nil {
		log.Fatal(err)
	}

	return node
}

func (n *Node) HandleNewStream(s network.Stream) {
	n.Streams[s.Conn().RemotePeer().String()] = s
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go n.ReadEvent(rw)
}

func (n *Node) InitSync(peerID peer.ID) error {
	s := n.Streams[peerID.String()]
	if s == nil {
		return fmt.Errorf("no stream to peer %s", peerID)
	}

	return n.SendEvent(s, core.RequestBlockChainSyncEvent())
}

func (n *Node) setupDiscovery(ctx context.Context, host libhost.Host, dht *dht.IpfsDHT) {
	discoveryAddrChan := make(chan string, 10)
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
			peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
			if err != nil {
				log.Fatal("failed to parse peer:", err)
			}

			s, err := ConnectToHost(ctx, host, addr)
			if err != nil {
				log.Println(err)
				continue
			}
			rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
			go n.ReadEvent(rw)

			n.Streams[peerInfo.ID.String()] = s
		}
	}(ctx, host, discoveryAddrChan)
}

func (n *Node) registerCallbacks() {
	n.Host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(net network.Network, c network.Conn) {
			fmt.Printf("Peer disconnected: %s\n", c.RemotePeer())
			delete(n.Streams, c.RemotePeer().String())
		},
	})

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

			if _, found := n.eventsCache.Get(eventData.ID); found {
				return
			}

			n.eventsCache.Set(eventData.ID, struct{}{}, cache.DefaultExpiration)
			go n.EventHandler.HandleEvent(rw, eventData)
			go n.PropagateEvent(eventData)
		}
	}
}

func (n *Node) SendEvent(s network.Stream, event core.EventData) error {
	event.Metadata.OriginPeerID = n.ID.String()
	event.Metadata.FromPeerID = n.ID.String()

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

func (n *Node) SendBroadcastEvent(event core.EventData) {
	event.Metadata.OriginPeerID = n.ID.String()
	event.Metadata.FromPeerID = n.ID.String()

	for _, s := range n.Streams {
		err := n.SendEvent(s, event)
		if err != nil {
			log.Println("Failed to send event", err)
		}
	}
}

func (n *Node) PropagateEvent(eventData core.EventData) {
	if !eventData.Metadata.MustPropagate {
		return
	}

	oldFromPeerID := eventData.Metadata.FromPeerID
	newFromPeerID := n.ID.String()
	originFromPeerID := eventData.Metadata.OriginPeerID

	eventData.Metadata.FromPeerID = newFromPeerID

	excludePeers := map[string]bool{
		oldFromPeerID:    true,
		newFromPeerID:    true,
		originFromPeerID: true,
	}

	for _, s := range n.Streams {
		if _, found := excludePeers[s.Conn().RemotePeer().String()]; found {
			continue
		}

		err := n.SendEvent(s, eventData)
		if err != nil {
			log.Println("Failed to send event", err)
		}
	}
}
