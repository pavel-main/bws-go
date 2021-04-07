package credentials

import (
	"testing"

	"github.com/pavel-main/bws-go/config"
	"github.com/pavel-main/bws-go/utils"
	"github.com/stretchr/testify/assert"
)

func TestNewInvalidBitSize(t *testing.T) {
	credentials, err := New(config.NewPublic(), 127)
	assert.Error(t, err, "should return invalid bitsize error")
	assert.Nil(t, credentials, "should return nil credentials")
}

func TestNewInvalidMnemonic(t *testing.T) {
	credentials, err := NewFromMnemonic(config.NewPublicTestnet(), "hello world")
	assert.Error(t, err, "should return invalid mnemonic")
	assert.Nil(t, credentials, "should return nil credentials")
}

func TestNewFromPrivateKey(t *testing.T) {
	privateKey := "tprv8ZgxMBicQKsPetcGAZY273DFjDSopBXJNEwFtK7nfCAnAficDoYmTGBRMLHxNoNdpxawo11wnfPoERHbqAcbbn7svZxunP55HPJeNSKoRUZ"
	credentials, err := NewFromPrivateKey(config.NewPublicTestnet(), privateKey)
	assert.NoError(t, err, "should create new credentials from private key string")

	rootPrvKey := utils.ToHex(credentials.RootPrvKey.Serialize())
	rootPubKey := utils.ToHex(credentials.RootPubKey.SerializeCompressed())
	reqPrvKey := utils.ToHex(credentials.ReqPrvKey.Serialize())
	reqPubKey := utils.ToHex(credentials.ReqPubKey.SerializeCompressed())
	accExtPubKey := credentials.AccExtPubKey.String()

	assert.Equal(t, "40a1b9b4645418203057571e8f41f4aa44a3eb435e7b026a45684386c29ac168", rootPrvKey, "should create root private key")
	assert.Equal(t, "033b0be659cb32418403e13c4a8fb8ce382eeed6b8391fc6dc5bfb5ce9c4c3c51b", rootPubKey, "should create root public key")
	assert.Equal(t, "2b07ec22f254a9b2695f57e28375f316a6cf59c6ee556dc2c45432bb45369a87", reqPrvKey, "should create request private key")
	assert.Equal(t, "0357449b15b27543d586856455ca8272e75b6e14fa9bf5e62c1b49cc25b416afe3", reqPubKey, "should create request public key")
	assert.Equal(t, "tpubDCuK338XJGn4aWo4WWaBKK2HwW7gSPNY3nvdhz22XdKi3xqRbPBtKRs98at7MatX3drTjjhhe12AZjVV6QcRcssTLRXa774Mo7pVj15BYqt", accExtPubKey, "should create account extended public key")

}

func TestNewFromMnemonic(t *testing.T) {
	mnemonic := "cause panel agent rare face frog dune congress thought assault urban impose"
	credentials, err := NewFromMnemonic(config.NewPublicTestnet(), mnemonic)
	assert.NoError(t, err, "should create new credentials from mnemonic")

	rootPrvKey := utils.ToHex(credentials.RootPrvKey.Serialize())
	rootPubKey := utils.ToHex(credentials.RootPubKey.SerializeCompressed())
	reqPrvKey := utils.ToHex(credentials.ReqPrvKey.Serialize())
	reqPubKey := utils.ToHex(credentials.ReqPubKey.SerializeCompressed())
	accExtPubKey := credentials.AccExtPubKey.String()

	assert.Equal(t, "69b140483788e80c0d8276c50dea8cf5bb8da422176a4f7da967a64a9730bee4", rootPrvKey, "should create root private key")
	assert.Equal(t, "033fe1dec3372dd059514d8f648b8f4c5c090f8e29987cc2636eae7a3a0d4a5ce7", rootPubKey, "should create root public key")
	assert.Equal(t, "05f1c94c6f40da65e35b60521e247bfbce7f2600870dba63d6b7b1fcb98b3169", reqPrvKey, "should create request private key")
	assert.Equal(t, "0270e9e8c9f9471d75f27aafd0e3a2aea2ac6ff7160086fabcf09428e695b01adb", reqPubKey, "should create request public key")
	assert.Equal(t, "tpubDDtcJNdS3crzf8oQrhaPFPBw2UY58tSBGFi8XNgmJKMcgkrsuaN7reXkFxFH3Kb3ru5faZkGnBuWBEQwxBTAnBxUkY2nYbr5Vet4hYcbkJB", accExtPubKey, "should create account extended public key")
}

func TestNewFromMnemonicWithPassphrase(t *testing.T) {
	mnemonic := "cause panel agent rare face frog dune congress thought assault urban impose"
	credentials, err := NewFromMnemonicWithPasshprase(config.NewPublicTestnet(), mnemonic, "helloworld")
	assert.NoError(t, err, "should create new credentials from mnemonic")

	rootPrvKey := utils.ToHex(credentials.RootPrvKey.Serialize())
	rootPubKey := utils.ToHex(credentials.RootPubKey.SerializeCompressed())
	reqPrvKey := utils.ToHex(credentials.ReqPrvKey.Serialize())
	reqPubKey := utils.ToHex(credentials.ReqPubKey.SerializeCompressed())
	accExtPubKey := credentials.AccExtPubKey.String()

	assert.Equal(t, "39b3d4ce6c010a26e177709ffee596716c70fd0ab4f4f4884f4ea0d46cc76c25", rootPrvKey, "should create root private key")
	assert.Equal(t, "02e355fdafbfe516126db777cfb01eeb233c443e772505ecc74629468c6ba74e50", rootPubKey, "should create root public key")
	assert.Equal(t, "7fe157924b4897b7777ba44a3f46a107091cc1ad4478da50d11f5a373a1bb40b", reqPrvKey, "should create request private key")
	assert.Equal(t, "035a53e0d345c777007510e7f8ecf4ba1f2c23fc8044393b3101aa1d7a1524ab87", reqPubKey, "should create request public key")
	assert.Equal(t, "tpubDCaeU9nacMnL1CLpGTYKgkKzUw9TEKPhYjJNomdCetw8WjCPbAv7wWSD2t3usFQ1t4h3rJAMgR7Fs24icCUVCuiCA3NiPRNvAUHT59voput", accExtPubKey, "should create account extended public key")
}

func TestDerivePath(t *testing.T) {
	hdkey := "tprv8ZgxMBicQKsPf1Zu9VstrFcmfHVRBibGLcTKn4ZxEYZkxR8fzUQsj1B49LRze1JpL2GAkL5GbqingWSqcW3cNNngt736xpeLJbYE6mHjaRr"
	credentials, err := NewFromPrivateKey(config.NewPublicTestnet(), hdkey)
	assert.NoError(t, err, "should create new credentials from import")

	expected := "d52aeb6fd6d9db37f1ae326eab69a45f1d026b600db4db9d69a4e465b8411ff4"
	privKey, _, err := credentials.DeriveFromAccount("m/1/4")
	assert.NoError(t, err, "should derive private key from account key")
	assert.Equal(t, expected, utils.ToHex(privKey.Serialize()), "should return valid derived key")

	privKey, _, err = credentials.DeriveFromAccount("m/abc/4")
	assert.Error(t, err, "should not derive private key if invalid path provided")
	assert.Nil(t, privKey, "should not derive private key if invalid path provided")

	expected = "962434db3cabd2e36443fd01a19273576aa6f50aa5a7dc6b94fc629e60aa39ca"
	privKey, _, err = credentials.DeriveFromAccount("m/1'/4")
	assert.NoError(t, err, "should derive private key from hardened account key")
	assert.Equal(t, expected, utils.ToHex(privKey.Serialize()), "should return valid derived key")
}
