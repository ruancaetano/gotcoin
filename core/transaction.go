package core

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/asn1"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"
)

type Transaction struct {
	FromAddr  string  `json:"from-addr"`
	ToAddr    string  `json:"to-addr"`
	Amount    float64 `json:"amount"`
	Timestamp int64   `json:"timestamp"`
	Hash      string  `json:"hash"`
	Signature string  `json:"signature"`
}

func NewTransaction(fromAddr, toAddr string, amount float64) *Transaction {
	t := &Transaction{
		FromAddr:  fromAddr,
		ToAddr:    toAddr,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
	}
	hashBytes := t.CalculateHash()
	t.Hash = hex.EncodeToString(hashBytes[:])
	return t
}

func (t *Transaction) IsCoinbase() bool {
	return t.FromAddr == ""
}

func (t *Transaction) CalculateHash() [32]byte {
	value := fmt.Sprintf("%s-%s-%f-%d", t.FromAddr, t.ToAddr, t.Amount, t.Timestamp)
	return sha256.Sum256([]byte(value))
}

func (t *Transaction) Sign(privateKeyString string) error {
	privateKeyBytes, _ := hex.DecodeString(privateKeyString)
	publicKeyBytes, _ := hex.DecodeString(t.FromAddr)
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKeyBytes)

	// Create an ECDSA private key
	privateKey := &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     x,
			Y:     y,
		},
		D: new(big.Int).SetBytes(privateKeyBytes),
	}

	hashBytes := t.CalculateHash()
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashBytes[:])
	if err != nil {
		return errors.New("failed to sign transaction")
	}

	signature, err := asn1.Marshal(struct{ R, S *big.Int }{r, s})
	if err != nil {
		return errors.New("failed to marshal signature")
	}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)

	t.Signature = encodedSignature
	return nil
}

func (t *Transaction) IsValid() bool {
	if t.IsCoinbase() {
		return true
	}

	publicKeyBytes, _ := hex.DecodeString(t.FromAddr)
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), publicKeyBytes)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(t.Signature)
	if err != nil {
		return false
	}

	hashBytes := t.CalculateHash()
	return ecdsa.VerifyASN1(&publicKey, hashBytes[:], signatureBytes)
}
