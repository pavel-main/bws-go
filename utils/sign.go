package utils

import (
	"github.com/btcsuite/btcd/btcec"
)

// SignMessage takes message hash and generates ECDSA signature with private key
func SignMessage(data []byte, privateKey *btcec.PrivateKey) ([]byte, error) {
	signature, err := privateKey.Sign(Reverse(HashMessage(data)))
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}
