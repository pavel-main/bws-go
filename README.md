# bws-go

[![Build Status](https://travis-ci.com/pavel-main/bws-go.svg?branch=main)](https://travis-ci.com/pavel-main/bws-go) 
[![codecov](https://codecov.io/gh/pavel-main/bws-go/branch/main/graph/badge.svg?token=v356EeCTBE)](https://codecov.io/gh/pavel-main/bws-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/pavel-main/bws-go)](https://goreportcard.com/report/github.com/pavel-main/bws-go)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/pavel-main/bws-go)](https://pkg.go.dev/github.com/pavel-main/bws-go)

[Bitcore Wallet Service](https://github.com/bitpay/bitcore-wallet-service) API client implementation in Go.

# Limitations

Messages (e.g. in transaction / tx outputs) are **NOT** being encrypted, because Go doesn't support [AES-CCM](https://en.wikipedia.org/wiki/CCM_mode), which is used by default in [SJCL](https://github.com/bitwiseshiftleft/sjcl), original implementation's cryptography dependency. More info:

* [proposal: crypto/tls: add support for AES-CCM](https://github.com/golang/go/issues/27484)
* [SJCL: In Other Languages](https://github.com/bitwiseshiftleft/sjcl/wiki/In-Other-Languages)

# Methods

Implemented API [methods](https://github.com/bitpay/bitcore-wallet-client#class-api):

- [x] `getFeeLevels`
- [x] `getVersion`
- [x] `createWallet`
- [x] `joinWallet`
- [x] `getNotifications`
- [x] `getStatus`
- [x] `getPreferences`
- [x] `savePreferences`
- [x] `getUtxos`
- [x] `createTxProposal`
- [x] `publishTxProposal`
- [x] `createAddress`
- [x] `getMainAddresses`
- [x] `getBalance`
- [x] `getTxProposals`
- [x] `signTxProposal`
- [x] `rejectTxProposal`
- [x] `broadcastRawTx`
- [x] `broadcastTxProposal`
- [x] `removeTxProposal`
- [x] `getTxHistory`
- [x] `getTx`
- [x] `startScan`
- [x] `getFiatRate`
- [x] `pushNotificationsSubscribe`
- [x] `pushNotificationsUnsubscribe`
- [x] `getSendMaxInfo`
- [ ] `recreateWallet`
- [ ] `fetchPayPro`
- [ ] `signTxProposalAirGapped`
- [ ] `createWalletFromOldCopay`

# Examples

* [examples/simple](examples/simple/main.go) - open existing wallet by a single copayer and send transaction
* [examples/multisig](examples/multisig/main.go) - create & join multi-signature wallet and send transaction
