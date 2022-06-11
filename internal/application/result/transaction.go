package result

type Transaction struct {
  SenderBlockchainAddress    string  `json:"sender_blockchain_address"`
  RecipientBlockchainAddress string  `json:"recipient_blockchain_address"`
  Value                      float64 `json:"value"`
}
