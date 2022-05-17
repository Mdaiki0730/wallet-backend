package model

import (
  "fmt"
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/rand"
  "encoding/json"

  "gariwallet/pkg/address"
)

type Wallet struct {
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet() *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	w.blockchainAddress = address.New(w.publicKey.X.Bytes(), w.privateKey.X.Bytes())

	return w
}

func (w *Wallet) PrivateKey() *ecdsa.PrivateKey {
	return w.privateKey
}

func (w *Wallet) PrivateKeyStr() string {
	return fmt.Sprintf("%x", w.privateKey.D.Bytes())
}

func (w *Wallet) PublicKey() *ecdsa.PublicKey {
	return w.publicKey
}

func (w *Wallet) PublicKeyStr() string {
	return fmt.Sprintf("%064x%064x", w.publicKey.X.Bytes(), w.publicKey.Y.Bytes())
}

func (w *Wallet) BlockchainAddress() string {
	return w.blockchainAddress
}

func (w *Wallet) MarshalJSON() ([]byte, error) {
  return json.Marshal(struct{
    PrivateKey string `json:"private_key"`
    PublicKey string `json:"public_key"`
    BlockchainAddress string `json:"blockchain_address"`
  }{
    PrivateKey: w.PrivateKeyStr(),
    PublicKey: w.PublicKeyStr(),
    BlockchainAddress: w.BlockchainAddress(),
  })
}
