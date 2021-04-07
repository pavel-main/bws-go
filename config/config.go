// Package config contains data structures related to Client configuration
package config

import (
	"errors"
	"net/url"

	"github.com/btcsuite/btcd/chaincfg"
)

const (
	publicAPI = "https://bws.bitpay.com/bws/api"
	localAPI  = "http://localhost:3232/bws/api"
)

// List of supported coins
const (
	CoinBTC = "btc"
	CoinBCH = "bch"
)

// BIP-44 Coin types
const (
	CoinTypeBTC  uint32 = 0
	CoinTypeTEST uint32 = 1
	CoinTypeBCH  uint32 = 145
)

// List of supported networks
const (
	NetworkLive = "livenet"
	NetworkTest = "testnet"
)

// Config contains Client configuration
type Config struct {
	Debug    bool
	BaseURL  string
	Coin     string
	Network  string
	Timeout  int
	Deadline int
}

// NewPublic creates new instance of Config for public API, BTC and livenet
func NewPublic() *Config {
	return newConfig(publicAPI, CoinBTC, NetworkLive)
}

// NewLocal creates new instance of Config for localhost API, BTC and livenet
func NewLocal() *Config {
	return newConfig(localAPI, CoinBTC, NetworkLive)
}

// NewPublicTestnet creates new instance of Config for public API, BTC and testnet
func NewPublicTestnet() *Config {
	return newConfig(publicAPI, CoinBTC, NetworkTest)
}

// NewLocalTestnet creates new instance of Config for localhost API, BTC and testnet
func NewLocalTestnet() *Config {
	return newConfig(localAPI, CoinBTC, NetworkTest)
}

// NewCashPublic creates new instance of Config for public API, BCH and livenet
func NewCashPublic() *Config {
	return newConfig(publicAPI, CoinBCH, NetworkLive)
}

// NewCashLocal creates new instance of Config for localhost API, BCH and livenet
func NewCashLocal() *Config {
	return newConfig(localAPI, CoinBCH, NetworkLive)
}

// NewCashPublicTestnet creates new instance of Config for public API, BCH and livenet
func NewCashPublicTestnet() *Config {
	return newConfig(publicAPI, CoinBCH, NetworkTest)
}

// NewCustom creates new instance of Config with custom parameters
func NewCustom(baseURL, coin, network string) (*Config, error) {
	// Validate input
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, err
	}

	if coin != CoinBTC && coin != CoinBCH {
		return nil, errors.New("Invalid coin name")
	}

	if network != NetworkLive && network != NetworkTest {
		return nil, errors.New("Invalid network name")
	}

	return newConfig(baseURL, coin, network), nil
}

func newConfig(baseURL, coin, network string) *Config {
	// Create and return config
	c := new(Config)
	c.Debug = false
	c.BaseURL = baseURL
	c.Coin = coin
	c.Network = network
	c.Timeout = 10000
	c.Deadline = 10000
	return c
}

// CoinType returns BIP-44 coin type
func (cfg *Config) CoinType() uint32 {
	if cfg.Network == NetworkTest {
		return CoinTypeTEST
	}

	if cfg.Coin == CoinBCH {
		return CoinTypeBCH
	}

	return CoinTypeBTC
}

// NetParams return network params for btcd
func (cfg *Config) NetParams() *chaincfg.Params {
	if cfg.Network == NetworkTest {
		return &chaincfg.TestNet3Params
	}

	return &chaincfg.MainNetParams
}

// NetShort returns network ID in one symbol
func (cfg *Config) NetShort() string {
	if cfg.Network == NetworkTest {
		return "T"
	}

	return "L"
}
