package models

// Transaction represents transaction from history
type Transaction struct {
	TxID                 string      `json:"txid"`
	BlockHeight          uint64      `json:"blockheight"`
	ProposalID           string      `json:"proposalId"`
	CreatedOn            uint        `json:"createdOn"`
	CreatorName          string      `json:"creatorName"`
	Action               string      `json:"action"`
	Actions              []*Action   `json:"actions"`
	Amount               int64       `json:"amount"`
	AddressTo            string      `json:"addressTo"`
	Fees                 int64       `json:"fees"`
	Time                 uint        `json:"time"`
	Confirmations        uint        `json:"confirmations"`
	FeePerKB             uint        `json:"feePerKb"`
	Outputs              []*TxOutput `json:"outputs"`
	HasUnconfirmedInputs bool        `json:"hasUnconfirmedInputs"`
	LowFees              bool        `json:"lowFees"`
}
