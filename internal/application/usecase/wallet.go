package usecase

import (
	"context"
	"errors"

	"gariwallet/internal/application/command"
	"gariwallet/internal/application/result"
	"gariwallet/internal/domain/model"
	"gariwallet/internal/domain/repoif"
)

type WalletAppIF interface {
	Create(ctx context.Context, cmd command.WalletCreate) (*result.Wallet, error)
}

type walletApp struct {
	repository repoif.Wallet
}

func NewWalletApp(r repoif.Wallet) WalletAppIF {
	return &walletApp{r}
}

func (wa *walletApp) Create(ctx context.Context, cmd command.WalletCreate) (*result.Wallet, error) {
	// create wallet instance
	wallet := model.NewWallet(cmd.IdpId)

	// check wallet existance
	obj, _ := wa.repository.FindById(ctx, wallet.IdpId())
	if obj != nil {
		return nil, errors.New("your wallet is already exist")
	}

	// insert database
	if err := wa.repository.InsertOne(ctx, wallet); err != nil {
		return nil, err
	}

	// data transfer dto
	result := result.Wallet{wallet.BlockchainAddress()}
	return &result, nil
}
