package model

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math/big"
)

type Transaction struct {
	senderBlockchainAddress    string
	recipientBlockchainAddress string
	value                      float64
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (s *Signature) String() string {
	return fmt.Sprintf("%064x%064x", s.R, s.S)
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

func (t *Transaction) UnmarshalJSON(data []byte) error {
	v := &struct {
		SenderBlockchainAddress    *string  `json:"sender_blockchain_address"`
		RecipientBlockchainAddress *string  `json:"recipient_blockchain_address"`
		Value                      *float64 `json:"value"`
	}{
		SenderBlockchainAddress:    &t.senderBlockchainAddress,
		RecipientBlockchainAddress: &t.recipientBlockchainAddress,
		Value:                      &t.value,
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	return nil
}
