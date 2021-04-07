package models

// Wallet represents generic wallet data structure
type Wallet struct {
	ID                 string `json:"id"`
	Version            string `json:"version"`
	CreatedOn          uint   `json:"createdOn"`
	M                  uint   `json:"m"`
	N                  uint   `json:"n"`
	SingleAddress      bool   `json:"singleAddress"`
	Status             string `json:"status"`
	PubKey             string `json:"pubKey"`
	Coin               string `json:"coin"`
	Network            string `json:"network"`
	DerivationStrategy string `json:"derivationStrategy"`
	AddressType        string `json:"addressType"`
}
