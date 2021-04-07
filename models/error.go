package models

// Error represents HTTP response embedded error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
