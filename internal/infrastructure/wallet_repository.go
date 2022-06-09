package infrastructure

import (
	"context"

	"gariwallet/internal/domain/model"
	"gariwallet/internal/domain/repoif"
	"gariwallet/internal/infrastructure/dbmodel"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type walletRepository struct {
	collection *mongo.Collection
}

func NewWalletRepository(client *mongo.Client) repoif.Wallet {
	return &walletRepository{client.Database("BLOCKCHAIN").Collection("wallets")}
}

func (wr *walletRepository) InsertOne(ctx context.Context, wallet *model.Wallet) error {
	dbModel := &dbmodel.Wallet{*wallet}
	err := insertOne(ctx, wr.collection, dbModel)
	if err != nil {
		return err
	}
	return nil
}

func (wr *walletRepository) FindById(ctx context.Context, idpId string) (*model.Wallet, error) {
	obj, err := findOne[dbmodel.Wallet](ctx, wr.collection, bson.D{{"idp_id", idpId}})
	if err != nil {
		return nil, err
	}
	return obj.ConvertDomainModel(), nil
}

func (wr *walletRepository) DeleteById(ctx context.Context, idpId string) error {
	err := deleteOne(ctx, wr.collection, bson.D{{"idp_id", idpId}})
	if err != nil {
		return err
	}
	return nil
}
