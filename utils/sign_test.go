package utils

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
)

func TestSignMessage(t *testing.T) {
	key, err := hex.DecodeString("09458c090a69a38368975fb68115df2f4b0ab7d1bc463fc60c67aa1730641d6c")
	assert.NoError(t, err, "should decode private key from string")

	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), key)
	assert.NotNil(t, pk, "should create private key from bytes")

	signature, err := SignMessage([]byte("hola"), pk)
	assert.NoError(t, err, "should sign message with private key")

	result := "3045022100f2e3369dd4813d4d42aa2ed74b5cf8e364a8fa13d43ec541e4bc29525e0564c302205b37a7d1ca73f684f91256806cdad4b320b4ed3000bee2e388bcec106e0280e0"
	assert.Equal(t, result, ToHex(signature), "should create a signed message")
}

func TestSignVerify(t *testing.T) {
	key, err := hex.DecodeString("09458c090a69a38368975fb68115df2f4b0ab7d1bc463fc60c67aa1730641d6c")
	assert.NoError(t, err, "should decode private key from string")

	pk, _ := btcec.PrivKeyFromBytes(btcec.S256(), key)
	assert.NotNil(t, pk, "should create private key from bytes")

	msg := []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.")
	signature, err := SignMessage(msg, pk)
	assert.NoError(t, err, "should sign message with private key")

	verify, err := VerifyMessage(msg, signature, pk.PubKey())
	assert.NoError(t, err, "should parse signature bytes")
	assert.Equal(t, true, verify, "should verify signed message")
}
