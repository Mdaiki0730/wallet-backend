package command

type TransactionCreate struct {
	IdpId                      string  `json:"idp_id"`
	RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
	Value                      float64 `json:"value"`
}
