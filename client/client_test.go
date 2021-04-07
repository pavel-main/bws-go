package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/AlekSi/pointer"
	"github.com/pavel-main/bws-go/config"
	"github.com/pavel-main/bws-go/credentials"
	"github.com/pavel-main/bws-go/models"
	"github.com/stretchr/testify/assert"
)

var (
	rootKey = "tprv8ZgxMBicQKsPf1Zu9VstrFcmfHVRBibGLcTKn4ZxEYZkxR8fzUQsj1B49LRze1JpL2GAkL5GbqingWSqcW3cNNngt736xpeLJbYE6mHjaRr"
	secret  = "3FfGK34vwMvVFDee1Ff1sVL1mQi4VbfaHAZ6fQm9PzVcUmLWotY3XLTQfPhUAbEDv3rkp8uNFXTbtc"
	rawTx   = "0100000001206dad82a8b1e1f2b3a854292d317467559d1ccedc3739cf5bdaaee97347a36301000000490047522102b5d05a90c5fe9a302a41acd67538ec9a507d893a706ea242eca0479ccd46254b21038f2f0c0d2fd79d166bc5968b0d87b8b8e827205eb8b75c5172ce8d0a5f6ef55952aeffffffff02184e12000000000017a9143b60a37a82414fd4cd15efe51093e68a8f665f0d87e0930400000000001976a914512c17fd86d596deb44c230ae4a98efae01f013c88ac00000000"
)

var mockAddress = &models.Address{
	Version:   "3",
	CreatedOn: 1538469762,
	Coin:      config.CoinBTC,
	Address:   "2NG1nqt5gieNfAULwAVbzdmDK3MnvAmQueE",
	WalletID:  "123e4567-e89b-12d3-a456-426655440000",
}

var mockTxProposal = &models.TxProposal{
	ID:          "123e4567-e89b-12d3-a456-426655440000",
	WalletID:    "123e4567-e89b-12d3-a456-426655440000",
	TxID:        "b53d75d4b45b574d8200c2539b0af761dddc94aa2005047643e2a5a71a695d52",
	Version:     3,
	CreatedOn:   1538469762,
	Coin:        config.CoinBTC,
	Network:     config.NetworkTest,
	Amount:      1958820,
	Fee:         246,
	OutputOrder: []int{1, 0},
	Inputs: []*models.TxInput{
		{
			TxID:         "0d5e1687d8f3dc24532798f25dcd9719d7148766b4516ac81e8e33bda54979b4",
			Path:         "m/1/4",
			Satoshis:     14110412,
			ScriptPubKey: "76a9143874e8eb8a2e018c46721fcdd8e9049f619af35488ac",
			Vout:         1,
		},
	},
	Outputs: []*models.TxOutput{
		{
			Amount:    1958820,
			ToAddress: "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY",
		},
	},
	ChangeAddress: &models.Address{
		Address: "mj4exG7YrSTxvpvXyFapoVRjNn9hMvYG1C",
	},
}

var mockTx = &models.Transaction{
	TxID:          "cd123d19e7b85f672213064e78cd868c53e84ef8b5ab079ecead251025b0e170",
	Action:        "sent",
	Amount:        12345,
	Fees:          100,
	Time:          1538558571,
	AddressTo:     "mnv9rH2VfAUX9YZzFkoRysGFtggvz1wRnY",
	Confirmations: 0,
	FeePerKB:      1088,
}

var mockWallet = &models.Wallet{
	ID:                 "123e4567-e89b-12d3-a456-426655440000",
	Version:            "3",
	CreatedOn:          1538469762,
	M:                  1,
	N:                  1,
	SingleAddress:      false,
	Status:             "pending",
	PubKey:             "xpub6DizKHgPEv7pwABAi7ZFe1k3Z1ds6ShjiXfYEvwWrj4KNXc1UBr14TVmy3WKGbVm8TQiCMpWQYzEekdQpFSNtnXTBZDrtftH8CjN8p4bXzh",
	Coin:               config.CoinBTC,
	Network:            config.NetworkTest,
	DerivationStrategy: "BIP44",
	AddressType:        "P2PKH",
}

type ScenarioCallback func(t *testing.T, expected, res interface{}, err error, msg string)

type Scenario struct {
	Status   int
	Expected interface{}
	Callback ScenarioCallback
	Message  string
}

func happyCallback(t *testing.T, expected, res interface{}, err error, msg string) {
	assert.NoError(t, err, msg)
	assert.Equal(t, expected, res, msg)
}

func errorCallback(t *testing.T, expected, res interface{}, err error, msg string) {
	assert.Nil(t, res, msg)
	assert.Error(t, err, msg)
}

func newClientServer(t *testing.T, status int, expected interface{}) (*httptest.Server, *Client) {
	// Init handler
	handler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(expected)
	})

	// Init server
	server := httptest.NewServer(handler)

	// Init config
	cfg, err := config.NewCustom(server.URL, config.CoinBTC, config.NetworkTest)
	if err != nil {
		assert.FailNow(t, "Error loading test credentials", err)
	}

	// Init credentials from private key string
	credentials, err := credentials.NewFromPrivateKey(cfg, rootKey)
	if err != nil {
		assert.FailNow(t, "Error loading test credentials", err)
	}

	// Init BWS client
	client, err := New(cfg, credentials)
	if err != nil {
		assert.FailNow(t, "Error initializing BWS client", err)
	}

	return server, client
}

func TestInvalidUrl(t *testing.T) {
	server, client := newClientServer(t, 200, nil)
	client.cfg.BaseURL = "httpz://example.com"
	defer server.Close()

	response, err := client.GetVersion()
	assert.Nil(t, response, "should not return response")
	assert.Error(t, err, "should return error")
}

func TestVersion(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.Version{ServiceVersion: "bws-2.4.0"},
			Callback: happyCallback,
			Message:  "version should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetVersion()
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestCreateWallet(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status: http.StatusOK,
			Expected: &models.WalletCreate{
				WalletID: mockWallet.ID,
				Secret:   secret,
			},
			Callback: happyCallback,
			Message:  "wallet ID should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.CreateWallet("test", 1, 1, false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestJoinWallet(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.WalletJoin{Wallet: mockWallet},
			Callback: happyCallback,
			Message:  "wallet should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.JoinWallet("test", secret)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetStatus(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.WalletStatus{Wallet: mockWallet},
			Callback: happyCallback,
			Message:  "wallet should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetStatus(false, false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetMaxInfo(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.MaxInfo{Amount: 10000, Fee: 1000},
			Callback: happyCallback,
			Message:  "send max info should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response1, err := client.GetMaxInfo("", 0)
		scenario.Callback(t, scenario.Expected, response1, err, scenario.Message)

		response2, err := client.GetMaxInfo("urgent", 1)
		scenario.Callback(t, scenario.Expected, response2, err, scenario.Message)
	}
}

func TestGetUtxos(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status: http.StatusOK,
			Expected: []*models.TxInput{
				{TxID: "002d489a24bf1275a68514dbc7c6d0636ed8fc2ebcd5e86b05d674e31f74bb78", Vout: 0},
				{TxID: "7d737d8a75fbd073d4afeb4bbd520a8f661bb884622b18c8fef3406a105d8501", Vout: 1},
			},
			Callback: happyCallback,
			Message:  "tx inputs should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response1, err := client.GetUtxos([]string{})
		scenario.Callback(t, scenario.Expected, response1, err, scenario.Message)

		response2, err := client.GetUtxos([]string{mockAddress.Address})
		scenario.Callback(t, scenario.Expected, response2, err, scenario.Message)
	}
}

func TestGetPreferences(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.Preferences{Language: "en", Unit: "btc"},
			Callback: happyCallback,
			Message:  "prefs should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetPreferences()
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestSavePreferences(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: nil,
			Callback: happyCallback,
			Message:  "response should match",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		err := client.SavePreferences(map[string]interface{}{"language": "en"})
		scenario.Callback(t, scenario.Expected, nil, err, scenario.Message)
	}
}

func TestGetBalance(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status: http.StatusOK,
			Expected: &models.Balance{
				TotalAmount:              10000,
				LockedAmount:             5000,
				TotalConfirmedAmount:     10000,
				LockedConfirmedAmount:    5000,
				AvailableAmount:          10000,
				AvailableConfirmedAmount: 5000,
			},
			Callback: happyCallback,
			Message:  "balance should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetBalance(false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestFeeLevels(t *testing.T) {
	expected := []*models.FeeLevel{
		{Level: "urgent", NumBlocks: 2},
		{Level: "priority", NumBlocks: 2},
		{Level: "normal", NumBlocks: 3},
		{Level: "economy", NumBlocks: 6},
		{Level: "superEconomy", NumBlocks: 24},
	}

	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: expected,
			Callback: happyCallback,
			Message:  "fee levels should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetFeeLevels()
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestCreateAddress(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockAddress,
			Callback: happyCallback,
			Message:  "address should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.CreateAddress(false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetMainAddresses(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: []*models.Address{mockAddress},
			Callback: happyCallback,
			Message:  "address should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetMainAddresses(0, false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetFiatRate(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.FiatRate{Timestamp: pointer.ToInt(1539155386351), Rate: 50000, FetchedOn: 1539155386351},
			Callback: happyCallback,
			Message:  "fiat rate should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetFiatRate("EUR", "", time.Now())
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestStartScan(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: &models.AddressScan{Started: true},
			Callback: happyCallback,
			Message:  "address should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.StartScan(false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetTx(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetTx(mockTxProposal.ID)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetTxProposals(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: []*models.TxProposal{mockTxProposal},
			Callback: happyCallback,
			Message:  "tx proposals should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetTxProposals()
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetTxHistory(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: []*models.Transaction{mockTx},
			Callback: happyCallback,
			Message:  "tx proposals should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.GetTxHistory(0, 0, false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestCreateTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.CreateTxProposal(nil, "", false)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestPublishTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.PublishTxProposal(mockTxProposal)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestSignTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.SignTxProposal(mockTxProposal)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestRejectTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.RejectTxProposal(mockTxProposal.ID, "reasons")
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestBroadcastRawTx(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: pointer.ToString(mockTxProposal.TxID),
			Callback: happyCallback,
			Message:  "txID should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response1, err := client.BroadcastRawTx(rawTx, "")
		scenario.Callback(t, scenario.Expected, response1, err, scenario.Message)

		response2, err := client.BroadcastRawTx(rawTx, config.NetworkTest)
		scenario.Callback(t, scenario.Expected, response2, err, scenario.Message)
	}
}

func TestBroadcastTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.BroadcastTxProposal(mockTxProposal.ID)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestRemoveTxProposal(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: mockTxProposal,
			Callback: happyCallback,
			Message:  "tx proposal should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response, err := client.RemoveTxProposal(mockTxProposal.ID)
		scenario.Callback(t, scenario.Expected, response, err, scenario.Message)
	}
}

func TestGetNotifications(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status: http.StatusOK,
			Expected: []*models.Notification{
				{
					ID: "1", Type: "test", Version: "3",
				},
				{
					ID: "2", Type: "test", Version: "3",
				},
			},
			Callback: happyCallback,
			Message:  "notifications should match",
		},
		{
			Status:   http.StatusOK,
			Expected: []int{1, 2, 3},
			Callback: errorCallback,
			Message:  "should return error if bad response",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		response1, err := client.GetNotifications("", 0, false)
		scenario.Callback(t, scenario.Expected, response1, err, scenario.Message)

		response2, err := client.GetNotifications("1", 1, true)
		scenario.Callback(t, scenario.Expected, response2, err, scenario.Message)
	}
}

func TestPushNotificationsSubscribe(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: nil,
			Callback: happyCallback,
			Message:  "notifications should match",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		err := client.PushNotificationsSubscribe("ios", "test")
		scenario.Callback(t, scenario.Expected, nil, err, scenario.Message)
	}
}

func TestPushNotificationsUnsubscribe(t *testing.T) {
	scenarios := []*Scenario{
		{
			Status:   http.StatusOK,
			Expected: nil,
			Callback: happyCallback,
			Message:  "notifications should match",
		},
		{
			Status:   http.StatusNotFound,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error if 404 not found",
		},
		{
			Status:   http.StatusInternalServerError,
			Expected: nil,
			Callback: errorCallback,
			Message:  "should return error on status is >= 4xx",
		},
	}

	for _, scenario := range scenarios {
		server, client := newClientServer(t, scenario.Status, scenario.Expected)
		defer server.Close()

		err := client.PushNotificationsUnsubscribe("test")
		scenario.Callback(t, scenario.Expected, nil, err, scenario.Message)
	}
}
