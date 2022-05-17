package usecase

import (
  "context"

  "gariwallet/internal/application/command"
  "gariwallet/internal/application/dto"
)

type WalletAppIF interface {
  Create(ctx context.Context, cmd command.WalletCreate) (*dto.Wallet, error)
}

type walletApp struct {}

func NewWalletApp() WalletAppIF {
  return &walletApp{}
}

func (wa *walletApp) Create(ctx context.Context, cmd command.WalletCreate) (*dto.Wallet, error) {
  return nil, nil
}
