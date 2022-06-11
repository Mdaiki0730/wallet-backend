package usecase

import (
	"context"
	"errors"
	"fmt"

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
	fmt.Println(wallet.PrivateKeyStr())
	signature := transaction.GenerateSignature(wallet.PrivateKey()).String()
	fmt.Println(signature, transaction)

	// data transfer dto
	result := result.Transaction{wallet.BlockchainAddress(), cmd.RecipientBlockchainAddress, cmd.Value}
	return &result, nil
}
