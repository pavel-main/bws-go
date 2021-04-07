package models

import (
	"bytes"
	"errors"
	"sort"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/pavel-main/bws-go/utils"
)

// TxProposal represents transaction proposal
type TxProposal struct {
	ID                      string      `json:"id"`
	TxID                    string      `json:"txId"`
	WalletID                string      `json:"walletId"`
	CreatorID               string      `json:"creatorId"`
	Version                 int32       `json:"version"`
	CreatedOn               uint        `json:"createdOn"`
	BroadcastedOn           uint        `json:"broadcastedOn"`
	Coin                    string      `json:"coin"`
	Network                 string      `json:"network"`
	Message                 *string     `json:"message"`
	PayProURL               *string     `json:"payProUrl"`
	WalletM                 int         `json:"walletM"`
	WalletN                 int         `json:"walletN"`
	RequiredSignatures      uint        `json:"requiredSignatures"`
	RequiredRejections      uint        `json:"requiredRejections"`
	Status                  string      `json:"status"`
	FeeLevel                string      `json:"feeLevel"`
	FeePerKB                uint        `json:"feePerKb"`
	ExcludeUnconfirmedUtxos bool        `json:"excludeUnconfrimedUtxos"`
	AddressType             string      `json:"addressType"`
	Amount                  int64       `json:"amount"`
	Fee                     int64       `json:"fee"`
	CreatorName             string      `json:"creatorName"`
	HasUnconfirmedInputs    bool        `json:"hasUnconfirmedInputs"`
	InputPaths              []string    `json:"inputPaths"`
	OutputOrder             []int       `json:"outputOrder"`
	ChangeAddress           *Address    `json:"changeAddress"`
	Inputs                  []*TxInput  `json:"inputs"`
	Outputs                 []*TxOutput `json:"outputs"`
	Actions                 []*TxAction `json:"actions"`
	redeemScript            []byte
}

// Validate performs basic validation before serialization
func (txp *TxProposal) Validate() error {
	// Validate outputs (second one is implied by ChangeAddress)
	if len(txp.Outputs) != 1 {
		return errors.New("Only transactions with two outputs are supported")
	}

	// Validate output order
	if len(txp.OutputOrder) != 2 {
		return errors.New("Invalid output order")
	}

	// Validate change address
	if txp.ChangeAddress == nil {
		return errors.New("Change address not specified")
	}

	return nil
}

// BuildRedeemScript builds multisig redeem script
func (txp *TxProposal) BuildRedeemScript(input *TxInput, net *chaincfg.Params) ([]byte, error) {
	sort.Strings(input.PublicKeys)

	multiSigBuilder := txscript.NewScriptBuilder().AddInt64(int64(txp.WalletM))
	for _, key := range input.PublicKeys {
		bytes, err := utils.ToBytes(key)
		if err != nil {
			return nil, err
		}

		multiSigBuilder.AddData(bytes)
	}

	multiSigBuilder.AddInt64(int64(len(input.PublicKeys)))
	multiSigBuilder.AddOp(txscript.OP_CHECKMULTISIG)

	multiSigScript, err := multiSigBuilder.Script()
	if err != nil {
		return nil, err
	}

	return multiSigScript, nil
}

// ToTransaction converts tx proposal to btcd transaction type
func (txp *TxProposal) ToTransaction(net *chaincfg.Params) (*wire.MsgTx, error) {
	// Basic validation
	if err := txp.Validate(); err != nil {
		return nil, err
	}

	// Build tx
	tx := wire.NewMsgTx(wire.TxVersion)
	to := txp.Outputs[0].ToAddress // Only single destination is now supported

	// Decode destination address
	toDest, err := btcutil.DecodeAddress(to, net)
	if err != nil {
		return nil, err
	}

	// Build destination output
	toScript, err := txscript.PayToAddrScript(toDest)
	if err != nil {
		return nil, err
	}

	// Add inputs
	var change int64
	var inputErr error
	if txp.WalletM >= 2 {
		change, inputErr = txp.AddMultisigInputs(tx, net)
	} else {
		change, inputErr = txp.AddInputs(tx, net)
	}

	if inputErr != nil {
		return nil, inputErr
	}

	// Decode change address
	toChange, err := btcutil.DecodeAddress(txp.ChangeAddress.Address, net)
	if err != nil {
		return nil, err
	}

	// Build change output
	changeScript, err := txscript.PayToAddrScript(toChange)
	if err != nil {
		return nil, err
	}

	// Default output order
	if txp.OutputOrder[0] == 0 {
		tx.AddTxOut(wire.NewTxOut(txp.Amount, toScript))

		if change > 0 {
			tx.AddTxOut(wire.NewTxOut(change, changeScript))
		}
	} else {
		if change > 0 {
			tx.AddTxOut(wire.NewTxOut(change, changeScript))
		}
		tx.AddTxOut(wire.NewTxOut(txp.Amount, toScript))
	}

	return tx, nil
}

func (txp *TxProposal) AddInputs(tx *wire.MsgTx, net *chaincfg.Params) (int64, error) {
	var total int64
	for _, input := range txp.Inputs {
		hash, err := chainhash.NewHashFromStr(input.TxID)
		if err != nil {
			return 0, err
		}

		outPoint := wire.NewOutPoint(hash, input.Vout)
		txInput := wire.NewTxIn(outPoint, nil, nil)
		tx.AddTxIn(txInput)
		total += input.Satoshis
	}

	// Wire input and outputs together
	change := total - txp.Amount - txp.Fee
	return change, nil
}

func (txp *TxProposal) AddMultisigInputs(tx *wire.MsgTx, net *chaincfg.Params) (int64, error) {
	var total int64
	for _, input := range txp.Inputs {
		if len(txp.redeemScript) == 0 {
			redeemScript, err := txp.BuildRedeemScript(input, net)
			if err != nil {
				return 0, err
			}

			txp.redeemScript = redeemScript
		}

		// Build multi-sig input script
		builder := txscript.NewScriptBuilder().
			AddInt64(txscript.OP_0).
			AddData(txp.redeemScript)

		script, err := builder.Script()
		if err != nil {
			return 0, err
		}

		hash, err := chainhash.NewHashFromStr(input.TxID)
		if err != nil {
			return 0, err
		}

		outPoint := wire.NewOutPoint(hash, input.Vout)
		txInput := wire.NewTxIn(outPoint, script, nil)
		tx.AddTxIn(txInput)
		total += input.Satoshis
	}

	// // Wire input and outputs together
	change := total - txp.Amount - txp.Fee
	return change, nil
}

// Serialize returns raw transaction from proposal data
func (txp *TxProposal) Serialize(net *chaincfg.Params) ([]byte, error) {
	// Convert
	tx, err := txp.ToTransaction(net)
	if err != nil {
		return nil, err
	}

	// Serialize
	buffer := bytes.NewBuffer([]byte{})
	if err := tx.Serialize(buffer); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// ProposalSignature serializes transaction and signs with private key
func (txp *TxProposal) ProposalSignature(privKey *btcec.PrivateKey, net *chaincfg.Params) ([]byte, error) {
	txBytes, err := txp.Serialize(net)
	if err != nil {
		return nil, err
	}

	hash := []byte(utils.ToHex(txBytes))
	proposalSignature, err := utils.SignMessage(hash, privKey)
	if err != nil {
		return nil, err
	}

	return proposalSignature, nil
}

// InputSignature signs transaction input
func (txp *TxProposal) InputSignature(privKey *btcec.PrivateKey, net *chaincfg.Params, idx int) ([]byte, error) {
	if len(txp.Inputs) == 0 {
		return nil, errors.New("Not enough inputs in transaction proposal")
	}

	tx, err := txp.ToTransaction(net)
	if err != nil {
		return nil, err
	}

	pkScript := txp.redeemScript
	if txp.WalletM == 1 {
		outputScript, err := utils.ToBytes(txp.Inputs[idx].ScriptPubKey)
		if err != nil {
			return nil, err
		}

		pkScript = outputScript
	}

	hash, err := txscript.CalcSignatureHash(pkScript, txscript.SigHashAll, tx, idx)
	if err != nil {
		return nil, err
	}

	signature, err := privKey.Sign(hash)
	if err != nil {
		return nil, err
	}

	return signature.Serialize(), nil
}
