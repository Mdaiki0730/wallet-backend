package model

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
)

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float64
}

func NewTransaction(sender string, recipient string, value float64) *Transaction {
	return &Transaction{sender, recipient, value}
}

func (t *Transaction) GenerateSignature(privateKey *ecdsa.PrivateKey) *Signature {
	m, _ := json.Marshal(t)
	h := sha256.Sum256([]byte(m))
	r, s, _ := ecdsa.Sign(rand.Reader, privateKey, h[:])
	return &Signature{r, s}
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
		Value                      float64 `json:"value"`
	}{
		SenderBlockchainAddress:    t.senderBlockchainAddress,
		RecipientBlockchainAddress: t.recipientBlockchainAddress,
		Value:                      t.value,
	})
}
