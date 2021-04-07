package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolToString(t *testing.T) {
	assert.Equal(t, "0", BoolToString(false), "should convert false to 0")
	assert.Equal(t, "1", BoolToString(true), "should convert true to 1")
}
