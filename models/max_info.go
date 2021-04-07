package models

// MaxInfo represents send maximum information
type MaxInfo struct {
	Amount             int64      `json:"amount"`
	Fee                int64      `json:"fee"`
	FeePerKB           uint       `json:"feePerKb"`
	Size               int        `json:"size"`
	Inputs             []*TxInput `json:"inputs"`
	UtxosBelowFee      uint       `json:"utxosBelowFee"`
	AmountBelowFee     int64      `json:"amountBelowFee"`
	UtxosAboveMaxSize  uint       `json:"utxosAboveMaxSize"`
	AmountAboveMaxSize int64      `json:"amountAboveMaxSize"`
}
