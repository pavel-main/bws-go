package models

// TxAction represents transcation proposal actions
type TxAction struct {
	Version     string   `jsons:"version"`
	CreatedOn   uint     `json:"createdOn"`
	Type        string   `json:"type"`
	CopayerID   string   `json:"copayerId"`
	Signatures  []string `json:"signatures"`
	XPub        string   `json:"xPub"`
	CopayerName string   `json:"copayerName"`
	Comment     string   `json:"comment"`
}
