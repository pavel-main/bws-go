package models

import (
	"encoding/hex"
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/pavel-main/bws-go/utils"
	"github.com/stretchr/testify/assert"
)

var (
	mockPrivateKey = "60aa7577a41d67187303877c9c5ee90610f9e577648a9260b63068b68bd580e7"
	net            = &chaincfg.TestNet3Params
)

var mockTxp = &TxProposal{
	Amount:      1958820,
	Fee:         246,
	OutputOrder: []int{0, 1},
	WalletN:     2,
	WalletM:     2,
	Inputs: []*TxInput{
		{
			TxID: "0d5e1687d8f3dc24532798f25dcd9719d7148766b4516ac81e8e33bda54979b4",
			Path: "m/1/4",
			PublicKeys: []string{
				"037cd8c7f67d1f6a7eadd2a61ce187e770385cdbc8935db644fd8050bd0ea81d05",
				"02a3957ab731ef0e17d856f73e69eb503b32266357e458e4cd2ff27902b2967efe",
			},
			Satoshis:     14110412,
			ScriptPubKey: "76a9143874e8eb8a2e018c46721fcdd8e9049f619af35488ac",
			Vout:         1,
		},
	},
	Outputs: []*TxOutput{
		{
			Amount:    1958820,
			ToAddress: "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY",
		},
	},
	ChangeAddress: &Address{
		Address: "mykbw8QcyMq9MeonF8626ayQYq4DVtiisK",
	},
}

func TestInputSignature(t *testing.T) {
	pkBytes, _ := hex.DecodeString(mockPrivateKey)
	requestPrivKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	expected := "304402207db432c75c7d4d5ae6cd87ff4a6895d761492990c9b2c7cd048175f200510846022064a6e3ae57b87e93c873ddefe2da1a37384c7b7792e82815d69610adcd2e450c"
	signature, err := mockTxp.InputSignature(requestPrivKey, net, 0)
	assert.NoError(t, err, "should sign transaction proposal input")
	assert.Equal(t, expected, utils.ToHex(signature), "hex-encoded signatures should match")
}

func TestProposalSignature(t *testing.T) {
	pkBytes, _ := hex.DecodeString(mockPrivateKey)
	requestPrivKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)

	expected := "30440220086ccf157ecd4ba9d600e2854052a64bedc43ac7e36e3378bc644f65242b49b802200b14cc1a16450b29b0f59371d8c70c6bb601626ab36417012348d9deac5bbab6"
	signature, err := mockTxp.ProposalSignature(requestPrivKey, net)
	assert.NoError(t, err, "should sign transaction proposal")
	assert.Equal(t, expected, utils.ToHex(signature), "hex-encoded signatures should match")
}

func TestSerializeHappyPath(t *testing.T) {
	txp := &TxProposal{
		Amount:      1958820,
		Fee:         246,
		OutputOrder: []int{1, 0},
		Inputs: []*TxInput{
			{
				TxID:     "0d5e1687d8f3dc24532798f25dcd9719d7148766b4516ac81e8e33bda54979b4",
				Satoshis: 14110412,
				Vout:     1,
			},
		},
		Outputs: []*TxOutput{
			{
				Amount:    1958820,
				ToAddress: "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY",
			},
		},
		ChangeAddress: &Address{
			Address: "mj4exG7YrSTxvpvXyFapoVRjNn9hMvYG1C",
		},
	}

	expected := "0100000001b47949a5bd338e1ec86a51b4668714d71997cd5df298275324dcf3d887165e0d0100000000ffffffff02326ab900000000001976a91426e7365e8b0a0bae05e7cfc320f8dc338dfdfa6b88aca4e31d00000000001976a914512c17fd86d596deb44c230ae4a98efae01f013c88ac00000000"
	rawTx, err := txp.Serialize(net)
	assert.NoError(t, err, "should serialize transaction")
	assert.Equal(t, expected, utils.ToHex(rawTx), "hex-encoded raw txs should match")
}

func TestSerializeExceptionPath(t *testing.T) {
	txp := &TxProposal{
		Amount:      1958820,
		Fee:         246,
		OutputOrder: []int{1, 0},
		Inputs: []*TxInput{
			{
				TxID:     "INVALID_HERE",
				Satoshis: 14110412,
				Vout:     1,
			},
		},
		Outputs: []*TxOutput{
			{
				Amount:    1958820,
				ToAddress: "INVALID_ADDRESS",
			},
		},
		ChangeAddress: &Address{
			Address: "INVALID_ADDRESS",
		},
	}

	// Fails on tx hash validation
	rawTx, err := txp.Serialize(net)
	assert.Nil(t, rawTx, "should not serialize invalid transaction")
	assert.Error(t, err, "should throw serialization error")

	// Fails on change address decoding
	txp.Inputs[0].TxID = "0d5e1687d8f3dc24532798f25dcd9719d7148766b4516ac81e8e33bda54979b4"
	rawTx, err = txp.Serialize(net)
	assert.Nil(t, rawTx, "should not serialize invalid transaction")
	assert.Error(t, err, "should throw serialization error")

	// Fails on destination address decoding
	txp.ChangeAddress.Address = "mykbw8QcyMq9MeonF8626ayQYq4DVtiisK"
	rawTx, err = txp.Serialize(net)
	assert.Nil(t, rawTx, "should not serialize invalid transaction")
	assert.Error(t, err, "should throw serialization error")

	// Suceeeds
	txp.Outputs[0].ToAddress = "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY"
	rawTx, err = txp.Serialize(net)
	assert.NoError(t, err, "should succeed finally")
	assert.NotNil(t, rawTx, "should serialize invalid transaction")
}

func TestValidateExceptionPath(t *testing.T) {
	// Fails on inputs validation
	txp := &TxProposal{}
	err := txp.Validate()
	assert.Error(t, err, "should throw validation error")

	// Fails on outputs validation
	input := &TxInput{TxID: "yolo", Vout: 0}
	txp.Inputs = []*TxInput{input}
	err = txp.Validate()
	assert.Error(t, err, "should throw validation error")

	// Fails on output order validation
	txp.Outputs = NewTxOutputSingle(10, "yolo")
	err = txp.Validate()
	assert.Error(t, err, "should throw validation error")

	// Fails on change address validation
	txp.OutputOrder = []int{1, 0}
	err = txp.Validate()
	assert.Error(t, err, "should throw validation error")

	// Succeeds
	txp.ChangeAddress = &Address{Address: "test"}
	err = txp.Validate()
	assert.NoError(t, err, "should succeed finally")
}
