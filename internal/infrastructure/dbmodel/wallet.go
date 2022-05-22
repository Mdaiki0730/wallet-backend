package dbmodel

import (
	"gariwallet/internal/domain/model"

	"go.mongodb.org/mongo-driver/bson"
)

type Wallet struct {
	model.Wallet
}

func (w *Wallet) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		PrivateKey        string `bson:"private_key"`
		PublicKey         string `bson:"public_key"`
		BlockchainAddress string `bson:"blockchain_address"`
	}{
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}