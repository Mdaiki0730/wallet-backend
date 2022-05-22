package repoif

import (
	"context"

	"gariwallet/internal/domain/model"
)

type Wallet interface {
	InsertOne(ctx context.Context, wallet *model.Wallet) error
}
