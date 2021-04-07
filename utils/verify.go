package utils

import (
	"github.com/btcsuite/btcd/btcec"
)

// VerifyMessage takes message hash and verifies it's signature using public key
func VerifyMessage(data, sig []byte, publicKey *btcec.PublicKey) (bool, error) {
	signature, err := btcec.ParseSignature(sig, btcec.S256())
	if err != nil {
		return false, err
	}

	hash := Reverse(HashMessage(data))
	return signature.Verify(hash, publicKey), nil
}
