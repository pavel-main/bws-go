package models

// Address represents receive address
type Address struct {
	Version     string   `json:"version"`
	CreatedOn   uint     `json:"createdOn"`
	Address     string   `json:"address"`
	WalletID    string   `json:"walletId"`
	IsChange    bool     `json:"isChange"`
	Path        string   `json:"path"`
	PublicKeys  []string `json:"publicKeys"`
	Coin        string   `json:"coin"`
	Network     string   `json:"network"`
	Type        string   `json:"type"`
	HasActivity *bool    `json:"hasActivity,omitempty"`
}
