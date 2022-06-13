package command

type WalletCreate struct {
	IdpId string `json:"idp_id"`
}

type WalletDelete struct {
	IdpId string `json:"idp_id"`
}

type WalletGet struct {
	IdpId       string `json:"idp_id"`
	AccessToken string `json:"access_token"`
}
