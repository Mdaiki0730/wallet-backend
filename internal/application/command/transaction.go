package command

type TransactionCreate struct {
	IdpId                      string  `json:"idp_id"`
	AccessToken                string  `json:"access_token"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value"`
}
