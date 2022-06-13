package usecase

import (
	"context"
	"errors"

	"gariwallet/internal/application/command"
	"gariwallet/internal/application/result"
	"gariwallet/internal/domain/bcsif"
	"gariwallet/internal/domain/model"
	"gariwallet/internal/domain/repoif"
)

type WalletAppIF interface {
	Create(ctx context.Context, cmd command.WalletCreate) (*result.Wallet, error)
	Delete(ctx context.Context, cmd command.WalletDelete) error
	Get(ctx context.Context, cmd command.WalletGet) (*result.Wallet, float64, error)
}

type walletApp struct {
	repository       repoif.Wallet
	blockchainServer bcsif.BlockchainServer
}

func NewWalletApp(r repoif.Wallet, bcs bcsif.BlockchainServer) WalletAppIF {
	return &walletApp{r, bcs}
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

func (wa *walletApp) Delete(ctx context.Context, cmd command.WalletDelete) error {
	if err := wa.repository.DeleteById(ctx, cmd.IdpId); err != nil {
		return err
	}
	return nil
}

func (wa *walletApp) Get(ctx context.Context, cmd command.WalletGet) (*result.Wallet, float64, error) {
	wallet, err := wa.repository.FindById(ctx, cmd.IdpId)
	if err != nil {
		return nil, 0, err
	}

	// get amount
	value, err := wa.blockchainServer.Amount(ctx, cmd.AccessToken, wallet.BlockchainAddress())
	if err != nil {
		return nil, 0, err
	}

	// data transfer dto
	result := result.Wallet{wallet.BlockchainAddress()}
	return &result, value, nil
}
