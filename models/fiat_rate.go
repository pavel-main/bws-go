package models

// FiatRate represents currency exchange rate
type FiatRate struct {
	Timestamp *int    `json:"ts"`
	FetchedOn int     `json:"fetchedOn"`
	Rate      float64 `json:"rate"`
}
