package models

// FeeLevel represents network-specific fee level
type FeeLevel struct {
	Level     string `json:"level"`
	FeePerKb  uint   `json:"feePerKb"`
	NumBlocks uint   `json:"nbBlocks"`
}
