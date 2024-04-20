package network

import (
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/infra"
)

func GetGenesisPeerID() string {
	genesisAddr, _ := multiaddr.NewMultiaddr(infra.GenesisNodeAddr)
	genesisPeerInfo, _ := peer.AddrInfoFromP2pAddr(genesisAddr)
	return genesisPeerInfo.ID.String()
}
