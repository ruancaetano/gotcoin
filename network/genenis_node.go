package network

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

const (
	GenesisNodeAddr = "/ip4/127.0.0.1/tcp/50000/p2p/QmPHZR5AdqSpKtcSZUCA5AqEni9sKgwPh7R6LxjZCLbXav"
	GenesisPort     = 50000
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
