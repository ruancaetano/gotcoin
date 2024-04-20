package network

import (
	"context"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	libhost "github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"

	"github.com/ruancaetano/gotcoin/infra"
)

func GetGenesisIdentity() (crypto.PrivKey, error) {
	// Open the PEM file for reading
	pemData, err := os.ReadFile("genesis_node_key.pem")
	if err != nil {
		panic(err)
	}

	// Decode the PEM data
	block, _ := pem.Decode(pemData)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the PEM-encoded private key
	privateKeyRsa, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Fatal(err)
	}

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKeyRsa)

	return crypto.UnmarshalRsaPrivateKey(privKeyBytes)
}

func InitHost(genesis bool, listenPort int) (libhost.Host, error) {
	var priv crypto.PrivKey
	var err error

	if genesis {
		listenPort = infra.GenesisPort
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

	host, err := libp2p.New(opts...)
	if err != nil {
		return nil, err
	}

	return host, nil
}

func ConnectHostToAddr(ctx context.Context, host libhost.Host, addr multiaddr.Multiaddr) (network.Stream, error) {
	peerInfo, err := peer.AddrInfoFromP2pAddr(addr)
	if err != nil {
		return nil, errors.New("failed to parse peer address")
	}

	host.Peerstore().AddAddrs(peerInfo.ID, peerInfo.Addrs, peerstore.PermanentAddrTTL)
	s, err := host.NewStream(ctx, peerInfo.ID, "/p2p/1.0.0")
	if err != nil {
		log.Println(err)
		return nil, errors.New("failed to open stream")
	}

	return s, nil
}
