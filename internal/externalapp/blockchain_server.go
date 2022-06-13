package externalapp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"gariwallet/internal/domain/bcsif"
	"gariwallet/pkg/config"
	"gariwallet/pkg/rest"
)

type blockchainServer struct{}

func NewBlockchainServer() bcsif.BlockchainServer {
	return &blockchainServer{}
}

type CreateTransactionBody struct {
	SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	SenderPublicKey            string  `json:"sender_public_key"`
	Value                      float64 `json:"value"`
	Signature                  string  `json:"signature"`
}

func (bcs *blockchainServer) CreateTransaction(ctx context.Context, token, sender, recipient, senderPubKey, signature string, value float64) error {
	reqBody, _ := json.Marshal(&CreateTransactionBody{
		SenderBlockchainAddress:    sender,
		RecipientBlockchainAddress: recipient,
		SenderPublicKey:            senderPubKey,
		Value:                      value,
		Signature:                  signature,
	})
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token),
	}
	endpoint := fmt.Sprintf("%s/v1/transactions", config.Global.BlockchainServerDomain)
	status, _, err := rest.Request("POST", endpoint, headers, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}
	if status != 200 {
		return errors.New("failed creating transaction")
	}
	return nil
}

func (bcs *blockchainServer) Amount(ctx context.Context, blockchainAddress string) (float64, error) {
	return 0, nil
}
