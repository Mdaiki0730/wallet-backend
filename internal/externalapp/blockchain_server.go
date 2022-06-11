package externalapp

import (
	"context"

	"gariwallet/internal/domain/bcsif"
	"gariwallet/internal/domain/model"
)

type blockchainServer struct{}

func NewBlockchainServer() bcsif.BlockchainServer {
	return &blockchainServer{}
}

func (bcs *blockchainServer) CreateTransaction(ctx context.Context, transaction model.Transaction) error {
	return nil
}

func (bcs *blockchainServer) Amount(ctx context.Context, blockchainAddress string) error {
	return nil
}
