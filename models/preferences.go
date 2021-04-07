package models

// Preferences represent user preferences
type Preferences struct {
	Version   string `json:"version"`
	CreatedOn uint   `json:"createdOn"`
	WalletID  string `json:"walletId"`
	CopayerID string `json:"copayerId"`
	Email     string `json:"email"`
	Language  string `json:"language"`
	Unit      string `json:"unit"`
}
