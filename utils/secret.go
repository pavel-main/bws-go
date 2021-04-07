package utils

import (
	"strings"

	"github.com/btcsuite/btcd/btcec"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcutil"
	"github.com/btcsuite/btcutil/base58"
)

// BuildSecret derives invitation secret for multi-signature wallets
func BuildSecret(privateKey *btcec.PrivateKey, walletID, coin, netShort string) (string, error) {
	walletFmtID, err := ToBytes(strings.Replace(walletID, "-", "", -1))
	if err != nil {
		return "", err
	}

	privKeyWIF, err := btcutil.NewWIF(privateKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return "", err
	}

	parts := []string{
		padEnd(base58.Encode(walletFmtID), "0", 22),
		privKeyWIF.String(),
		netShort,
		coin,
	}

	return strings.Join(parts, ""), nil
}

// ParseSecret parses shared wallet secret
func ParseSecret(secret string) (*btcec.PrivateKey, string, string, string, error) {
	secretSplit := split(secret, []int{22, 74, 75})
	widBase58 := strings.Replace(secretSplit[0], "0", "", -1)
	widHex := ToHex(base58.Decode(widBase58))
	walletID := strings.Join(split(widHex, []int{8, 12, 16, 20}), "-")

	wif, err := btcutil.DecodeWIF(secretSplit[1])
	if err != nil {
		return nil, "", "", "", err
	}

	// Detect network
	var network string
	if secretSplit[2] == "T" {
		network = "testnet"
	} else {
		network = "livenet"
	}

	// Detect coin
	coin := secretSplit[3]
	if len(coin) == 0 {
		coin = "btc"
	}

	return wif.PrivKey, walletID, coin, network, nil
}

func padEnd(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}

func split(input string, indexes []int) []string {
	parts := []string{}
	indexes = append(indexes, len(input))

	i := 0
	for i < len(indexes) {
		start := 0
		if i != 0 {
			start = indexes[i-1]
		}

		part := input[start:indexes[i]]
		parts = append(parts, part)
		i++
	}

	return parts
}
