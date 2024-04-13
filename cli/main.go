package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Encode the private key in PEM format
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	// Create a new file for writing
	file, err := os.Create("genesis_node_key.pem")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Write the PEM-encoded private key to the file
	if err := pem.Encode(file, privateKeyPEM); err != nil {
		panic(err)
	}
}
