package models

// TxInput represents transaction input
type TxInput struct {
	TxID          string   `json:"txid"`
	Vout          uint32   `json:"vout"`
	Address       string   `json:"address"`
	ScriptPubKey  string   `json:"scriptPubKey"`
	Satoshis      int64    `json:"satoshis"`
	Confirmations uint     `json:"confirmations"`
	Locked        bool     `json:"locked"`
	Path          string   `json:"path"`
	PublicKeys    []string `json:"publicKeys"`
}
