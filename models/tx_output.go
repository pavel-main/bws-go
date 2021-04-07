package models

// TxOutput represents transaction output
type TxOutput struct {
	Amount    int64   `json:"amount"`
	Address   string  `json:"address,omitempty"` // For received transactions
	ToAddress string  `json:"toAddress"`
	Message   *string `json:"message"`
}

// NewTxOutput creates new tx output without a message
func NewTxOutput(amount int64, toAddress string) *TxOutput {
	out := new(TxOutput)
	out.Amount = amount
	out.ToAddress = toAddress
	return out
}

// NewTxOutputSingle creates tx output list with a single output
func NewTxOutputSingle(amount int64, toAddress string) []*TxOutput {
	result := []*TxOutput{}
	result = append(result, NewTxOutput(amount, toAddress))
	return result
}
