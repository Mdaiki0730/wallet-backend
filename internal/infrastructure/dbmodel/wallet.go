package dbmodel

import (
	"gariwallet/internal/domain/model"
	"gariwallet/pkg/converter"

	"go.mongodb.org/mongo-driver/bson"
)

type Wallet struct {
	model.Wallet
}

func (w *Wallet) ConvertDomainModel() *model.Wallet {
	wallet := model.NewWallet("")
	wallet.SetIdpId(w.IdpId())
	wallet.SetPrivateKey(w.PrivateKeyStr(), converter.PublicKeyFromString(w.PublicKeyStr()))
	wallet.SetPublicKey(w.PublicKeyStr())
	wallet.SetBlockchainAddress(w.BlockchainAddress())
	return wallet
}

func (w *Wallet) MarshalBSON() ([]byte, error) {
	return bson.Marshal(struct {
		IdpId             string `bson:"idp_id"`
		PrivateKey        string `bson:"private_key"`
		PublicKey         string `bson:"public_key"`
		BlockchainAddress string `bson:"blockchain_address"`
	}{
		IdpId:             w.IdpId(),
		PrivateKey:        w.PrivateKeyStr(),
		PublicKey:         w.PublicKeyStr(),
		BlockchainAddress: w.BlockchainAddress(),
	})
}

func (w *Wallet) UnmarshalBSON(data []byte) error {
	v := &struct {
		IdpId             string `bson:"idp_id"`
		PrivateKey        string `bson:"private_key"`
		PublicKey         string `bson:"public_key"`
		BlockchainAddress string `bson:"blockchain_address"`
	}{}
	if err := bson.Unmarshal(data, &v); err != nil {
		return err
	}
	w.SetIdpId(v.IdpId)
	w.SetPrivateKey(v.PrivateKey, converter.PublicKeyFromString(v.PublicKey))
	w.SetPublicKey(v.PublicKey)
	w.SetBlockchainAddress(v.BlockchainAddress)
	return nil
}
