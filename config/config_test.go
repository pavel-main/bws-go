package config

import (
	"testing"

	"github.com/btcsuite/btcd/chaincfg"
	"github.com/stretchr/testify/assert"
)

func TestNewCustom(t *testing.T) {
	cfg, err := NewCustom(publicAPI, CoinBTC, NetworkLive)
	assert.NoError(t, err, "should create new custom config instance")
	assert.Equal(t, publicAPI, cfg.BaseURL, "should create public API targeted config")
	assert.Equal(t, CoinBTC, cfg.Coin, "should create BTC-targeted config")
	assert.Equal(t, NetworkLive, cfg.Network, "should create livenet-targeted config")
	assert.Equal(t, CoinTypeBTC, cfg.CoinType(), "should return BTC coin type")
	assert.Equal(t, &chaincfg.MainNetParams, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "L", cfg.NetShort(), "should return short network name")
}

func TestNewCustomErrors(t *testing.T) {
	_, err := NewCustom("", CoinBTC, NetworkLive)
	assert.Error(t, err, "should fail on URL validation")

	_, err = NewCustom(publicAPI, "", NetworkLive)
	assert.Error(t, err, "should fail on coin validation")

	_, err = NewCustom(publicAPI, CoinBTC, "")
	assert.Error(t, err, "should fail on network validation")
}

func TestNewPublic(t *testing.T) {
	cfg := NewPublic()
	assert.Equal(t, publicAPI, cfg.BaseURL, "should create public API targeted config")
	assert.Equal(t, CoinBTC, cfg.Coin, "should create BTC-targeted config")
	assert.Equal(t, NetworkLive, cfg.Network, "should create livenet-targeted config")
	assert.Equal(t, CoinTypeBTC, cfg.CoinType(), "should return BTC coin type")
	assert.Equal(t, &chaincfg.MainNetParams, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "L", cfg.NetShort(), "should return short network name")
}

func TestNewLocal(t *testing.T) {
	cfg := NewLocal()
	assert.Equal(t, localAPI, cfg.BaseURL, "should create localhost API targeted config")
	assert.Equal(t, CoinBTC, cfg.Coin, "should create BTC-targeted config")
	assert.Equal(t, NetworkLive, cfg.Network, "should create livenet-targeted config")
	assert.Equal(t, CoinTypeBTC, cfg.CoinType(), "should return BTC coin type")
	assert.Equal(t, &chaincfg.MainNetParams, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "L", cfg.NetShort(), "should return short network name")
}

func TestNewPublicTestnet(t *testing.T) {
	cfg := NewPublicTestnet()
	assert.Equal(t, publicAPI, cfg.BaseURL, "should create public API targeted config")
	assert.Equal(t, CoinBTC, cfg.Coin, "should create BTC-targeted config")
	assert.Equal(t, NetworkTest, cfg.Network, "should create testnet-targeted config")
	assert.Equal(t, CoinTypeTEST, cfg.CoinType(), "should return testnet coin type")
	assert.Equal(t, &chaincfg.TestNet3Params, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "T", cfg.NetShort(), "should return short network name")
}

func TestNewLocalTestnet(t *testing.T) {
	cfg := NewLocalTestnet()
	assert.Equal(t, localAPI, cfg.BaseURL, "should create localhost API targeted config")
	assert.Equal(t, CoinBTC, cfg.Coin, "should create BTC-targeted config")
	assert.Equal(t, NetworkTest, cfg.Network, "should create testnet-targeted config")
	assert.Equal(t, CoinTypeTEST, cfg.CoinType(), "should return testnet coin type")
	assert.Equal(t, &chaincfg.TestNet3Params, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "T", cfg.NetShort(), "should return short network name")
}

func TestNewCashPublic(t *testing.T) {
	cfg := NewCashPublic()
	assert.Equal(t, publicAPI, cfg.BaseURL, "should create public API targeted config")
	assert.Equal(t, CoinBCH, cfg.Coin, "should create BCH-targeted config")
	assert.Equal(t, NetworkLive, cfg.Network, "should create livenet-targeted config")
	assert.Equal(t, CoinTypeBCH, cfg.CoinType(), "should return BCH coin type")
	assert.Equal(t, &chaincfg.MainNetParams, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "L", cfg.NetShort(), "should return short network name")
}

func TestNewCashLocal(t *testing.T) {
	cfg := NewCashLocal()
	assert.Equal(t, localAPI, cfg.BaseURL, "should create localhost API targeted config")
	assert.Equal(t, CoinBCH, cfg.Coin, "should create BCH-targeted config")
	assert.Equal(t, NetworkLive, cfg.Network, "should create livenet-targeted config")
	assert.Equal(t, CoinTypeBCH, cfg.CoinType(), "should return BCH coin type")
	assert.Equal(t, &chaincfg.MainNetParams, cfg.NetParams(), "should return mainnet params")
	assert.Equal(t, "L", cfg.NetShort(), "should return short network name")
}

func TestNewCashPublicTestnet(t *testing.T) {
	cfg := NewCashPublicTestnet()
	assert.Equal(t, publicAPI, cfg.BaseURL, "should create public API targeted config")
	assert.Equal(t, CoinBCH, cfg.Coin, "should create BCH-targeted config")
	assert.Equal(t, NetworkTest, cfg.Network, "should create testnet-targeted config")
	assert.Equal(t, CoinTypeTEST, cfg.CoinType(), "should return testnet coin type")
	assert.Equal(t, &chaincfg.TestNet3Params, cfg.NetParams(), "should return testnet params")
	assert.Equal(t, "T", cfg.NetShort(), "should return short network name")
}
