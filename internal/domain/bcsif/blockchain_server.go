package bcsif

import (
	"context"
)

type BlockchainServer interface {
	CreateTransaction(ctx context.Context, token, sender, recipient, senderPubKey, signature string, value float64) error
	Amount(ctx context.Context, blockchainAddress string) (float64, error)
}
