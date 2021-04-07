package models

// WalletJoin represents wallet status response
type WalletJoin struct {
	Wallet *Wallet `json:"wallet,omitempty"`
}
