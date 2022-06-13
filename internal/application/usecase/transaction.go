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

type TransactionAppIF interface {
	Create(ctx context.Context, cmd command.TransactionCreate) (*result.Transaction, error)
}

type transactionApp struct {
	walletrepo       repoif.Wallet
	blockchainServer bcsif.BlockchainServer
}

func NewTransactionApp(r repoif.Wallet, bcs bcsif.BlockchainServer) TransactionAppIF {
	return &transactionApp{r, bcs}
}

func (ta *transactionApp) Create(ctx context.Context, cmd command.TransactionCreate) (*result.Transaction, error) {
	// get wallet existance
	wallet, err := ta.walletrepo.FindById(ctx, cmd.IdpId)
	if err != nil {
		return nil, errors.New("your wallet is already exist")
	}

	// create transaction instance
	transaction := model.NewTransaction(wallet.BlockchainAddress(), cmd.RecipientBlockchainAddress, cmd.Value)
	signature := transaction.GenerateSignature(wallet.PrivateKey()).String()
	err = ta.blockchainServer.CreateTransaction(
		ctx, cmd.AccessToken, wallet.BlockchainAddress(), cmd.RecipientBlockchainAddress, wallet.PublicKeyStr(), signature, cmd.Value)
	if err != nil {
		return nil, err
	}

	// data transfer dto
	result := result.Transaction{wallet.BlockchainAddress(), cmd.RecipientBlockchainAddress, cmd.Value}
	return &result, nil
}
