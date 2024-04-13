package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type Wallet struct {
	PrivateKey string `json:"private-key"`
	PublicKey  string `json:"public-key"`
}

func NewWallet() (*Wallet, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.New("failed to generate private key")
	}

	publicKey := &privateKey.PublicKey
	publicKeyBytes := elliptic.MarshalCompressed(publicKey.Curve, publicKey.X, publicKey.Y)
	return &Wallet{
		PrivateKey: hex.EncodeToString(privateKey.D.Bytes()),
		PublicKey:  hex.EncodeToString(publicKeyBytes),
	}, nil
}
