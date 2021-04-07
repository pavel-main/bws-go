package utils

import (
	"testing"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

type SecretTestFixture struct {
	PrivateKey string
	WalletID   string
	Coin       string
	NetShort   string
	Network    string
	Params     *chaincfg.Params
	Secret     string
}

var fixtures = []SecretTestFixture{
	{
		PrivateKey: "76c5f9a1c1549a762437e9ba4755c79375aea5bdccf7dc954e40cac92e92c657",
		WalletID:   "42037ba2-3d32-4b19-af2b-f8be6d3fe3bd",
		Coin:       "btc",
		NetShort:   "T",
		Network:    "testnet",
		Secret:     "99o84P9EjxEkfbjB3bfJ7nL1CbDUzrZgq1ifSHpR8xcDLoBA616A7djHhgnfZD6JfUL8rdxKWMTbtc",
	},
	{
		PrivateKey: "969ca4125efd4334a8d42c38d04c49638e7ce7a7179350e24f95b05d74c1b1e2",
		WalletID:   "04f6a3eb-2482-4c43-86f4-2ec36a442bfa",
		Coin:       "btc",
		NetShort:   "T",
		Network:    "testnet",
		Secret:     "cYpRkrKAsedQbvPMwFRoX0L2GUqJSgW5WkyqbVBByZR8sbp9gGrbsfzmHvkUyrv4ip6bNFHiwXTbtc",
	},
	{
		PrivateKey: "3931a7153692c11492bc3275ca09afe603ea52ae2f67b62a52c5dae3c6887ca1",
		WalletID:   "c1c3b65f-8f38-4fd7-a310-b95d7aa5b672",
		Coin:       "btc",
		NetShort:   "T",
		Network:    "testnet",
		Secret:     "QvkypTGW6gQ4HsEWbfnMFFKy8tV2N63MRjaYqJYehQ4RM3rrESe9GrsaiSDGdRKkf1oJypByFiTbtc",
	},
}

func TestBuildSecret(t *testing.T) {
	for _, fix := range fixtures {
		pkBytes, err := ToBytes(fix.PrivateKey)
		assert.NoError(t, err, "should convert hex-encoded private key to bytes")

		privateKey, _ := btcec.PrivKeyFromBytes(btcec.S256(), pkBytes)
		secret, err := BuildSecret(privateKey, fix.WalletID, fix.Coin, fix.NetShort)
		assert.NoError(t, err, "should generate secret from private key and walletID")
		assert.Equal(t, fix.Secret, secret, "should generate secret")
	}
}

func TestParseSecret(t *testing.T) {
	for _, fix := range fixtures {
		privKey, walletID, coin, network, err := ParseSecret(fix.Secret)
		assert.NoError(t, err, "should parse secret")
		assert.Equal(t, fix.PrivateKey, ToHex(privKey.Serialize()), "private key should match")
		assert.Equal(t, fix.WalletID, walletID, "wallet ID should match")
		assert.Equal(t, fix.Coin, coin, "coin should match")
		assert.Equal(t, fix.Network, network, "network should match")
		assert.Equal(t, fix.PrivateKey, ToHex(privKey.Serialize()), "private key should match")
	}
}
