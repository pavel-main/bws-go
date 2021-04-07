package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashMessage(t *testing.T) {
	hash := HashMessage([]byte("hola"))
	expected := "4102b8a140ec642feaa1c645345f714bc7132d4fd2f7f6202db8db305a96172f"
	assert.Equal(t, expected, ToHex(hash), "should create a reverse double sha256")
}

func TestHash160(t *testing.T) {
	hash, err := Hash160([]byte("hola"))
	expected := "77c8c5d8228355835c70b1e5be2bf7de906cd4b1"
	assert.NoError(t, err, "should create a sha256 + ripemd160")
	assert.Equal(t, expected, ToHex(hash), "should create a sha256 + ripemd160")
}

func TestHash160Error(t *testing.T) {
	hash, err := Hash160(nil)
	assert.Error(t, err, "should not hash empty data")
	assert.Nil(t, hash, "should not hash empty data")
}
