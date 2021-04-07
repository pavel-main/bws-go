package models

// Notification represents... well, notification
type Notification struct {
	ID        string                 `json:"id"`
	Type      string                 `json:"type"`
	Version   string                 `json:"version"`
	Data      map[string]interface{} `json:"data"`
	CreatedOn int                    `json:"createdOn"`
	CreatorID *string                `json:"creatorId"`
	WalletID  string                 `json:"walletId"`
}
