package model

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"

	"gariwallet/pkg/address"
	"gariwallet/pkg/converter"
)

type Wallet struct {
	idpId             string
	privateKey        *ecdsa.PrivateKey
	publicKey         *ecdsa.PublicKey
	blockchainAddress string
}

func NewWallet(idpId string) *Wallet {
	w := new(Wallet)
	privateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	w.idpId = idpId
	w.privateKey = privateKey
	w.publicKey = &w.privateKey.PublicKey
	w.blockchainAddress = address.New(w.publicKey.X.Bytes(), w.privateKey.X.Bytes())

	return w
}

func (w *Wallet) IdpId() string {
	return w.idpId
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
	return json.Marshal(struct {
		IdpId             string `json:"idp_id"`
		PrivateKey        string `json:"private_key"`
		PublicKey         string `json:"public_key"`
		BlockchainAddress string `json:"blockchain_address"`
	}{
		IdpId:             w.IdpId(),
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

// setters
func (w *Wallet) SetIdpId(idpId string) {
	w.idpId = idpId
}

func (w *Wallet) SetPrivateKey(privateKey string, publicKey *ecdsa.PublicKey) {
	w.privateKey = converter.PrivateKeyFromString(privateKey, publicKey)
}

func (w *Wallet) SetPublicKey(publicKey string) {
	w.publicKey = converter.PublicKeyFromString(publicKey)
}

func (w *Wallet) SetBlockchainAddress(blockchainAddress string) {
	w.blockchainAddress = blockchainAddress
}
