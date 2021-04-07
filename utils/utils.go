// Package utils contains functions for hashing, signing and verifying data
package utils

// Reverse reverses byte array
func Reverse(data []byte) []byte {
	for i := 0; i < len(data)/2; i++ {
		j := len(data) - i - 1
		data[i], data[j] = data[j], data[i]
	}
	return data
}

// BoolToString converts bool to int represented in string
func BoolToString(input bool) string {
	result := "0"
	if input {
		result = "1"
	}

	return result
}
