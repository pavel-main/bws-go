package models

// WalletCreate represents wallet creation response
type WalletCreate struct {
	WalletID string `json:"walletId"`
	Secret   string `json:"secret"`
}
