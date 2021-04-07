package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewTxOutput(t *testing.T) {
	expected := int64(10000)
	txOut := NewTxOutput(expected, "mueUUsavi1NaYQKqvtW9ANQtZscrcBt19j")
	assert.NotNil(t, txOut, "should create valid tx output")
	assert.Equal(t, expected, txOut.Amount, "amounts should match")
	assert.Equal(t, "mueUUsavi1NaYQKqvtW9ANQtZscrcBt19j", txOut.ToAddress, "addresses should match")
}

func TestNewTxOutputSingle(t *testing.T) {
	expected := int64(10000)
	txOuts := NewTxOutputSingle(expected, "mueUUsavi1NaYQKqvtW9ANQtZscrcBt19j")
	assert.Len(t, txOuts, 1, "should create tx out list with single item")
	assert.NotNil(t, txOuts[0], "should create valid tx output")
	assert.Equal(t, expected, txOuts[0].Amount, "amounts should match")
	assert.Equal(t, "mueUUsavi1NaYQKqvtW9ANQtZscrcBt19j", txOuts[0].ToAddress, "addresses should match")
}
