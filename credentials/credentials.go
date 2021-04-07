// Package credentials is related to HD keys management
package credentials

import (
	"strconv"
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/pavel-main/bws-go/config"
	bip39 "github.com/tyler-smith/go-bip39"
)

const hardening = "'"

// Credentials contains is BIP-39 Root Key and derivatives
type Credentials struct {
	RootKey      *hdkeychain.ExtendedKey
	RootPrvKey   *btcec.PrivateKey
	RootPubKey   *btcec.PublicKey
	ReqPrvKey    *btcec.PrivateKey
	ReqPubKey    *btcec.PublicKey
	AccExtKey    *hdkeychain.ExtendedKey
	AccExtPubKey *hdkeychain.ExtendedKey
}

// New creates new Credentials for livenet and random mnemonic
func New(cfg *config.Config, bitSize int) (*Credentials, error) {
	// Create entropy
	entropy, err := bip39.NewEntropy(bitSize)
	if err != nil {
		return nil, err
	}

	// Create mnemonic
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	return newFromMnemonic(mnemonic, "", cfg.CoinType(), cfg.NetParams())
}

// NewFromMnemonic creates new Credentials based on existing mnemonic
func NewFromMnemonic(cfg *config.Config, mnemonic string) (*Credentials, error) {
	return newFromMnemonic(mnemonic, "", cfg.CoinType(), cfg.NetParams())
}

// NewFromMnemonicWithPasshprase creates new Credentials based on existing mnemonic and passphrase
func NewFromMnemonicWithPasshprase(cfg *config.Config, mnemonic, passphrase string) (*Credentials, error) {
	return newFromMnemonic(mnemonic, passphrase, cfg.CoinType(), cfg.NetParams())
}

// NewFromPrivateKey creates new Credentials based on existing mnemonic
func NewFromPrivateKey(cfg *config.Config, privateKey string) (*Credentials, error) {
	return newFromPrivateKey(privateKey, cfg.CoinType(), cfg.NetParams())
}

func newFromPrivateKey(privateKey string, coinType uint32, net *chaincfg.Params) (*Credentials, error) {
	rootKey, err := hdkeychain.NewKeyFromString(privateKey)
	if err != nil {
		return nil, err
	}

	rootKey.SetNet(net)
	return deriveChildren(rootKey, coinType)
}

func newFromMnemonic(mnemonic, passphrase string, coinType uint32, net *chaincfg.Params) (*Credentials, error) {
	// Create seed
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, passphrase)
	if err != nil {
		return nil, err
	}

	// Create master key
	rootKey, err := hdkeychain.NewMaster(seed, net)
	if err != nil {
		return nil, err
	}

	return deriveChildren(rootKey, coinType)
}

// DeriveFromAccount derives child key pair from account extended key by provided BIP44-compliant path
func (c *Credentials) DeriveFromAccount(path string) (*btcec.PrivateKey, *btcec.PublicKey, error) {
	parts := strings.Split(path, "/")

	current := c.AccExtKey
	for _, part := range parts {
		if part == "m" {
			continue
		}

		idx := part
		hardened := false
		if strings.Contains(part, hardening) {
			idx = strings.Replace(part, hardening, "", -1)
			hardened = true
		}

		id, err := strconv.ParseUint(idx, 10, 64)
		if err != nil {
			return nil, nil, err
		}

		index := uint32(id)
		if hardened {
			index += hdkeychain.HardenedKeyStart
		}

		child, err := current.Child(index)
		if err != nil {
			return nil, nil, err
		}

		current = child
	}

	return toElliptic(current)
}

func deriveChildren(rootKey *hdkeychain.ExtendedKey, coinType uint32) (*Credentials, error) {
	// Derive request path (m/1')
	requestBaseKey, err := rootKey.Child(1 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Derive request key (m/1'/0)
	requestKey, err := requestBaseKey.Child(0)
	if err != nil {
		return nil, err
	}

	// Derive purpose (m/44')
	purposeKey, err := rootKey.Child(44 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Derive Ethereum coin type (m/44'/coin_type')
	coinTypeKey, err := purposeKey.Child(coinType + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Derive first account (m/44'/coin_type'/0')
	accExtKey, err := coinTypeKey.Child(0 + hdkeychain.HardenedKeyStart)
	if err != nil {
		return nil, err
	}

	// Neuter account public key
	accExtPubKey, err := accExtKey.Neuter()
	if err != nil {
		return nil, err
	}

	// Convert root keypair to ECDSA
	rootPrvKey, rootPubKey, err := toElliptic(rootKey)
	if err != nil {
		return nil, err
	}

	// Convert request keypair to ECDSA
	reqPrvKey, reqPubKey, err := toElliptic(requestKey)
	if err != nil {
		return nil, err
	}

	// Save
	k := new(Credentials)
	k.RootKey = rootKey
	k.RootPrvKey = rootPrvKey
	k.RootPubKey = rootPubKey
	k.ReqPrvKey = reqPrvKey
	k.ReqPubKey = reqPubKey
	k.AccExtKey = accExtKey
	k.AccExtPubKey = accExtPubKey
	return k, nil
}

// toElliptic converts ExtendedKey to private and public elliptic curve keys
func toElliptic(key *hdkeychain.ExtendedKey) (*btcec.PrivateKey, *btcec.PublicKey, error) {
	private, err := key.ECPrivKey()
	if err != nil {
		return nil, nil, err
	}

	public, err := key.ECPubKey()
	if err != nil {
		return nil, nil, err
	}

	return private, public, nil
}
