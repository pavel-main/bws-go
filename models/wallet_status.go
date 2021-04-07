package models

// WalletStatus represents wallet status response
type WalletStatus struct {
	Wallet *Wallet `json:"wallet,omitempty"`
}
