package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"golang.org/x/crypto/ripemd160"
)

// Sha256 calculates sha-256 sum
func Sha256(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

// HashMessage calculates double sha-256 and returns reverse bytes
func HashMessage(input []byte) []byte {
	return Reverse(Sha256(Sha256(input)))
}

// ToHex converts bytes to hex-encoded string
func ToHex(input []byte) string {
	return hex.EncodeToString(input)
}

// ToBytes converts hex-encoded string to bytes
func ToBytes(input string) ([]byte, error) {
	return hex.DecodeString(input)
}

// Hash160 performs the same operations as OP_HASH160 in Bitcoin Script
// It hashes the given data first with SHA256, then RIPEMD160
func Hash160(data []byte) ([]byte, error) {
	// Does identical function to Script OP_HASH160. Hash once with SHA-256, then RIPEMD-160
	if data == nil {
		return nil, errors.New("Empty bytes cannot be hashed")
	}

	hash := Sha256(data)
	ripemd160Hash := ripemd160.New()
	ripemd160Hash.Write(hash)
	hash = ripemd160Hash.Sum(nil)
	return hash, nil
}
