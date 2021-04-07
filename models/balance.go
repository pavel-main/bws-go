package models

// Balance represents wallet balance
type Balance struct {
	TotalAmount              uint64 `json:"totalAmount"`
	LockedAmount             uint64 `json:"lockedAmount"`
	TotalConfirmedAmount     uint64 `json:"totalConfirmedAmount"`
	LockedConfirmedAmount    uint64 `json:"lockedConfirmedAmount"`
	AvailableAmount          uint64 `json:"availableAmount"`
	AvailableConfirmedAmount uint64 `json:"availableConfirmedAmount"`
}
