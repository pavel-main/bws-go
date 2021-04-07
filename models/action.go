package models

// Action represents transaction update
type Action struct {
	CreatedOn   uint   `json:"createdOn"`
	Type        string `json:"type"`
	CopayerID   string `json:"copayerId"`
	CopayerName string `json:"copayerName"`
	Comment     string `json:"comment"`
}
