package utils

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/stretchr/testify/assert"
)

func TestVerifyMessage(t *testing.T) {
	publicKey, err := ToBytes("03bec86ad4a8a91fe7c11ec06af27246ec55094db3d86098b7d8b2f12afe47627f")
	assert.NoError(t, err, "should convert hex-encoded public key to bytes")

	signature, err := ToBytes("3045022100d6186930e4cd9984e3168e15535e2297988555838ad10126d6c20d4ac0e74eb502201095a6319ea0a0de1f1e5fb50f7bf10b8069de10e0083e23dbbf8de9b8e02785")
	assert.NoError(t, err, "should convert hex-encoded signature to bytes")

	key, err := btcec.ParsePubKey(publicKey, btcec.S256())
	assert.NoError(t, err, "should parse public key bytes")

	verify, err := VerifyMessage([]byte("hola"), signature, key)
	assert.NoError(t, err, "should parse signature bytes")
	assert.Equal(t, true, verify, "should verify signed message")
}

func TestVerifyMessageFailure(t *testing.T) {
	publicKey, err := ToBytes("02555a2d45e309c00cc8c5090b6ec533c6880ab2d3bc970b3943def989b3373f16")
	assert.NoError(t, err, "should convert hex-encoded public key to bytes")

	signature, err := ToBytes("3045022100d6186930e4cd9984e3168e15535e2297988555838ad10126d6c20d4ac0e74eb502201095a6319ea0a0de1f1e5fb50f7bf10b8069de10e0083e23dbbf8de9b8e02785")
	assert.NoError(t, err, "should convert hex-encoded signature to bytes")

	key, err := btcec.ParsePubKey(publicKey, btcec.S256())
	assert.NoError(t, err, "should parse public key bytes")

	verify, err := VerifyMessage([]byte("hola"), signature, key)
	assert.NoError(t, err, "should parse signature bytes")
	assert.Equal(t, false, verify, "should fail to verify signed message")
}

func TestVerifyMessageSignatureFailure(t *testing.T) {
	publicKey, err := ToBytes("02555a2d45e309c00cc8c5090b6ec533c6880ab2d3bc970b3943def989b3373f16")
	assert.NoError(t, err, "should convert hex-encoded public key to bytes")

	signature, err := ToBytes("3045022100d6186930e4cd9984e3168e15535e2297")
	assert.NoError(t, err, "should convert hex-encoded signature to bytes")

	key, err := btcec.ParsePubKey(publicKey, btcec.S256())
	assert.NoError(t, err, "should parse public key bytes")

	verify, err := VerifyMessage([]byte("hola"), signature, key)
	assert.Error(t, err, "should fail to parse signature bytes")
	assert.Equal(t, false, verify, "should fail to verify signed message")
}
