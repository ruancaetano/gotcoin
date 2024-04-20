package adapters

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	libhost "github.com/libp2p/go-libp2p/core/host"
	libnetwork "github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/core/events"
	"github.com/ruancaetano/gotcoin/core/protocols"
	"github.com/ruancaetano/gotcoin/infra"
	network2 "github.com/ruancaetano/gotcoin/infra/network"

	"github.com/patrickmn/go-cache"
)

type NodeAdapter struct {
	Genesis      bool
	ID           peer.ID
	Addr         string
	Host         libhost.Host
	EventHandler protocols.EventHandler
	Streams      map[string]libnetwork.Stream
	eventsCache  *cache.Cache
}

func NewNodeAdapterAsGenesis(ctx context.Context, host libhost.Host) protocols.Node {
	_, err := network2.NewKDHT(ctx, host, nil)
	if err != nil {
		log.Fatal(err)
	}

	return &NodeAdapter{
		Genesis:     true,
		ID:          host.ID(),
		Addr:        host.Addrs()[0].String(),
		Host:        host,
		Streams:     map[string]libnetwork.Stream{},
		eventsCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func NewNodeAdapter(ctx context.Context, host libhost.Host) protocols.Node {
	genesisAddr, err := multiaddr.NewMultiaddr(infra.GenesisNodeAddr)
	if err != nil {
		log.Fatal("failed to parse multiaddr:", err)
	}
	genesisPeerInfo, _ := peer.AddrInfoFromP2pAddr(genesisAddr)

	s, err := network2.ConnectHostToAddr(ctx, host, genesisAddr)
	if err != nil {
		log.Fatal(err)
	}

	return &NodeAdapter{
		ID:   host.ID(),
		Addr: host.Addrs()[0].String(),
		Host: host,
		Streams: map[string]libnetwork.Stream{
			genesisPeerInfo.ID.String(): s,
		},
		eventsCache: cache.New(5*time.Minute, 10*time.Minute),
	}
}

func (n *NodeAdapter) GetID() string {
	return n.ID.String()
}

func (n *NodeAdapter) GetAddr() string {
	return n.Addr
}

func (n *NodeAdapter) Setup(ctx context.Context, eh protocols.EventHandler) {
	fmt.Println("NodeID:", n.ID.String())
	fmt.Println("NodeAddr:", n.Addr)

	n.EventHandler = eh

	if n.Genesis {
		n.Host.SetStreamHandler("/p2p/1.0.0", n.HandleNewStream)
		n.registerNetworkListeners()
		return
	}

	genesisPeerID := network2.GetGenesisPeerID()
	genesisStream := n.Streams[genesisPeerID]

	n.Host.SetStreamHandler("/p2p/1.0.0", n.HandleNewStream)

	n.setupNetworkDiscovery(ctx)
	n.registerNetworkListeners()
	go n.ReadEvent(genesisStream)

	if err := n.InitSync(); err != nil {
		log.Fatal(err)
	}
}

func (n *NodeAdapter) HandleNewStream(s libnetwork.Stream) {
	n.Streams[s.Conn().RemotePeer().String()] = s
	go n.ReadEvent(s)
}

func (n *NodeAdapter) InitSync() error {
	peerID := network2.GetGenesisPeerID()
	s := n.Streams[peerID]
	if s == nil {
		return fmt.Errorf("no stream to peer %s", peerID)
	}

	return n.SendEvent(s, events.RequestBlockChainSyncEvent())
}

func (n *NodeAdapter) ReadEvent(s libnetwork.Stream) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	for {
		str, err := rw.ReadString('\n')
		if err != nil {
			return
		}

		if str == "" {
			return
		}
		if str != "\n" {
			eventData := events.EventData{}
			if err = json.Unmarshal([]byte(str), &eventData); err != nil {
				log.Println("Failed to unmarshal event data")
				continue
			}

			if _, found := n.eventsCache.Get(eventData.ID); found {
				return
			}

			n.eventsCache.Set(eventData.ID, struct{}{}, cache.DefaultExpiration)
			go n.EventHandler.Handle(eventData)
			go n.PropagateEvent(eventData)
		}
	}
}

func (n *NodeAdapter) SendEvent(s libnetwork.Stream, event events.EventData) error {
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

func (n *NodeAdapter) SendBroadcastEvent(event events.EventData) {
	event.Metadata.OriginPeerID = n.ID.String()
	event.Metadata.FromPeerID = n.ID.String()

	for _, s := range n.Streams {
		err := n.SendEvent(s, event)
		if err != nil {
			log.Println("Failed to send event", err)
		}
	}
}

func (n *NodeAdapter) PropagateEvent(eventData events.EventData) {
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

func (n *NodeAdapter) GetPeerStream(peerID string) libnetwork.Stream {
	return n.Streams[peerID]
}

func (n *NodeAdapter) registerNetworkListeners() {
	n.Host.Network().Notify(&libnetwork.NotifyBundle{
		DisconnectedF: func(net libnetwork.Network, c libnetwork.Conn) {
			fmt.Printf("Peer disconnected: %s\n", c.RemotePeer())
			delete(n.Streams, c.RemotePeer().String())
		},
	})
}

func (n *NodeAdapter) setupNetworkDiscovery(ctx context.Context) {
	genesisAddr, err := multiaddr.NewMultiaddr(infra.GenesisNodeAddr)
	if err != nil {
		log.Fatal(err)
	}

	kdht, err := network2.NewKDHT(ctx, n.Host, []multiaddr.Multiaddr{genesisAddr})
	if err != nil {
		log.Fatal(err)
	}

	discoveryAddrChan := make(chan string, 10)
	go network2.Discover(ctx, discoveryAddrChan, n.Host, kdht, "gotcoin")
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

			s, err := network2.ConnectHostToAddr(ctx, host, addr)
			if err != nil {
				log.Println(err)
				continue
			}
			go n.ReadEvent(s)

			n.Streams[peerInfo.ID.String()] = s
		}
	}(ctx, n.Host, discoveryAddrChan)
}
