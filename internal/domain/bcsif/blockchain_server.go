package bcsif

import (
	"context"

	"gariwallet/internal/domain/model"
)

type BlockchainServer interface {
	CreateTransaction(ctx context.Context, transaction model.Transaction) error
	Amount(ctx context.Context, blockchainAddress string) error
}
