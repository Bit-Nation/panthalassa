package dapp

type DApp struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Code        string `json:"code"`
	Signature   string `json:"signature"`
	Icon        string `json:"icon"`
	PublicKey   string `json:"app_public_key"`
	Development bool   `json:"development"`
}
