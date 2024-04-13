package network

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/core"
)

func InitNode(genesis bool, listenPort int) (host.Host, error) {
	var priv crypto.PrivKey
	var err error

	if genesis {
		listenPort = GenesisPort
		priv, err = GetGenesisIdentity()
	} else {
		priv, _, err = crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, rand.Reader)
	}
	if err != nil {
		return nil, err
	}

	opts := []libp2p.Option{
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	}

	node, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	peerInfo := peer.AddrInfo{
		ID:    node.ID(),
		Addrs: node.Addrs(),
	}
	addrs, err := peer.AddrInfoToP2pAddrs(&peerInfo)
	log.Printf("Running: %s\n", addrs)

	node.Network().Notify(&network.NotifyBundle{
		ConnectedF: func(n network.Network, c network.Conn) {
			fmt.Printf("Peer connected: %s\n", c.RemotePeer())
		},
		DisconnectedF: func(n network.Network, c network.Conn) {
			fmt.Printf("Peer disconnected: %s\n", c.RemotePeer())
		},
	})

	return node, nil
}

func HandleNewStream(s network.Stream, eventHandler *core.EventHandler) {
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go ReadEvent(rw, eventHandler)
}

func ReadEvent(rw *bufio.ReadWriter, eventHandler *core.EventHandler) {
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

			go eventHandler.HandleEvent(rw, eventData)
		}
	}
}

func SendEvent(rw *bufio.ReadWriter, event core.EventData) error {
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

func ConnectToNode(ctx context.Context, node host.Host, addr multiaddr.Multiaddr, eh *core.EventHandler) (*network.Stream, *bufio.ReadWriter, error) {
	peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil, nil, errors.New("failed to parse peer address")
	}

	node.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
	s, err := node.NewStream(ctx, peerInfo.ID, "/p2p/1.0.0")
	if err != nil {
		log.Println(err)
		return nil, nil, errors.New("failed to open stream")
	}

	if err != nil {
		log.Fatal(err)
	}

	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))
	go ReadEvent(rw, eh)
	return &s, rw, nil
}
