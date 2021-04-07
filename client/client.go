// Package client contains Bitcore Wallet Client
package client

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	httpclient "github.com/ddliu/go-httpclient"
	"github.com/pavel-main/bws-go/config"
	"github.com/pavel-main/bws-go/credentials"
	"github.com/pavel-main/bws-go/models"
	"github.com/pavel-main/bws-go/utils"
)

var clientVersion = "bws-go v2.0.0"

// Client is responsible for HTTP requests and signatures
type Client struct {
	cfg    *config.Config
	client *httpclient.HttpClient
	keys   *credentials.Credentials
}

// New creates new client instance based on Config, HttpClient and Credentials
func New(cfg *config.Config, chain *credentials.Credentials) (*Client, error) {
	c := new(Client)
	c.cfg = cfg
	c.keys = chain
	c.client = httpclient.Defaults(httpclient.Map{
		httpclient.OPT_DEBUG:             cfg.Debug,
		httpclient.OPT_CONNECTTIMEOUT_MS: cfg.Timeout,
		httpclient.OPT_TIMEOUT_MS:        cfg.Deadline,
		"Content-Type":                   "application/json",
		"Accept":                         "application/json",
		"User-Agent":                     clientVersion,
	})
	return c, nil
}

// GetVersion returns BWS service version
func (c *Client) GetVersion() (*models.Version, error) {
	bytes, err := c.doGetRequest("/v1/version/", map[string]string{})
	if err != nil {
		return nil, err
	}

	response := &models.Version{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateWallet creates wallet with provided name and required number of signatures
func (c *Client) CreateWallet(name string, m, n uint, singleAddress bool) (*models.WalletCreate, error) {
	payload := map[string]interface{}{
		"name":          name,
		"m":             m,
		"n":             n,
		"pubKey":        hex.EncodeToString(c.keys.RootPubKey.SerializeCompressed()),
		"coin":          c.cfg.Coin,
		"network":       c.cfg.Network,
		"singleAddress": singleAddress,
	}

	bytes, err := c.doPostRequest("/v2/wallets", payload)
	if err != nil {
		return nil, err
	}

	response := &models.WalletCreate{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	secret, err := utils.BuildSecret(c.keys.RootPrvKey, response.WalletID, c.cfg.Coin, c.cfg.NetShort())
	if err != nil {
		return nil, err
	}

	response.Secret = secret
	return response, nil
}

// JoinWallet joins existing wallet
func (c *Client) JoinWallet(name, secret string) (*models.WalletJoin, error) {
	privateKey, walletID, coin, _, err := utils.ParseSecret(secret)
	if err != nil {
		return nil, err
	}

	requestPubKey := c.keys.ReqPubKey.SerializeCompressed()
	copayerSignature, err := utils.SignMessage(c.copayerHash(name), privateKey)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"name":             name,
		"coin":             coin,
		"xPubKey":          c.keys.AccExtPubKey.String(),
		"requestPubKey":    hex.EncodeToString(requestPubKey),
		"copayerSignature": hex.EncodeToString(copayerSignature),
	}

	path := fmt.Sprintf("/v2/wallets/%s/copayers", walletID)
	bytes, err := c.doPostRequest(path, payload)
	if err != nil {
		return nil, err
	}

	response := &models.WalletJoin{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetStatus returns wallet status
func (c *Client) GetStatus(includeExtendedInfo, twoStep bool) (*models.WalletStatus, error) {
	params := map[string]string{
		"includeExtendedInfo": utils.BoolToString(includeExtendedInfo),
		"twoStep":             utils.BoolToString(twoStep),
	}

	bytes, err := c.doGetRequest("/v2/wallets/", params)
	if err != nil {
		return nil, err
	}

	response := &models.WalletStatus{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetMaxInfo returns send max information
func (c *Client) GetMaxInfo(feeLevel string, feePerKb uint64) (*models.MaxInfo, error) {
	params := map[string]string{}
	if feeLevel != "" {
		params["feeLevel"] = feeLevel
	}

	if feePerKb != 0 {
		params["feePerKb"] = strconv.FormatUint(feePerKb, 10)
	}

	bytes, err := c.doGetRequest("/v1/sendmaxinfo/", params)
	if err != nil {
		return nil, err
	}

	response := &models.MaxInfo{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetUtxos returns unspent transaction outputs
func (c *Client) GetUtxos(addresses []string) ([]*models.TxInput, error) {
	params := map[string]string{}
	if len(addresses) != 0 {
		params["addresses"] = strings.Join(addresses, ",")
	}

	bytes, err := c.doGetRequest("/v1/utxos/", params)
	if err != nil {
		return nil, err
	}

	response := []*models.TxInput{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPreferences returns copayer preferences
func (c *Client) GetPreferences() (*models.Preferences, error) {
	bytes, err := c.doGetRequest("/v1/preferences/", map[string]string{})
	if err != nil {
		return nil, err
	}

	response := &models.Preferences{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// SavePreferences updates copayer preferences
func (c *Client) SavePreferences(payload map[string]interface{}) error {
	_, err := c.doPutRequest("/v1/preferences/", payload)
	if err != nil {
		return err
	}

	return nil
}

// GetBalance returns wallet balance
func (c *Client) GetBalance(twoStep bool) (*models.Balance, error) {
	params := map[string]string{
		"twoStep": utils.BoolToString(twoStep),
	}

	bytes, err := c.doGetRequest("/v1/balance/", params)
	if err != nil {
		return nil, err
	}

	response := &models.Balance{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetFeeLevels returns current network fee levels
func (c *Client) GetFeeLevels() ([]*models.FeeLevel, error) {
	params := map[string]string{
		"coin":    c.cfg.Coin,
		"network": c.cfg.Network,
	}

	bytes, err := c.doGetRequest("/v2/feelevels/", params)
	if err != nil {
		return nil, err
	}

	response := []*models.FeeLevel{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateAddress creates new receiving address
func (c *Client) CreateAddress(ignoreMaxGap bool) (*models.Address, error) {
	payload := map[string]interface{}{
		"ignoreMaxGap": utils.BoolToString(ignoreMaxGap),
	}

	bytes, err := c.doPostRequest("/v3/addresses/", payload)
	if err != nil {
		return nil, err
	}

	response := &models.Address{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetMainAddresses returns generated addresses
func (c *Client) GetMainAddresses(limit int, reverse bool) ([]*models.Address, error) {
	payload := map[string]string{
		"limit":   fmt.Sprint(limit),
		"reverse": utils.BoolToString(reverse),
	}

	bytes, err := c.doGetRequest("/v1/addresses/", payload)
	if err != nil {
		return nil, err
	}

	response := []*models.Address{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// StartScan starts an address scanning process
func (c *Client) StartScan(includeCopayerBranches bool) (*models.AddressScan, error) {
	payload := map[string]interface{}{
		"includeCopayerBranches": utils.BoolToString(includeCopayerBranches),
	}

	bytes, err := c.doPostRequest("/v1/addresses/scan", payload)
	if err != nil {
		return nil, err
	}

	response := &models.AddressScan{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetFiatRate returns exchange rate for the specified currency & timestamp.
func (c *Client) GetFiatRate(code, provider string, ts time.Time) (*models.FiatRate, error) {
	if len(provider) == 0 {
		provider = "BitPay"
	}

	payload := map[string]string{
		"ts":       strconv.FormatInt(ts.Unix(), 10),
		"provider": provider,
	}

	path := fmt.Sprintf("/v1/fiatrates/%s/", code)
	bytes, err := c.doGetRequest(path, payload)
	if err != nil {
		return nil, err
	}

	response := &models.FiatRate{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetTx returns transaction proposal
func (c *Client) GetTx(txID string) (*models.TxProposal, error) {
	path := fmt.Sprintf("/v1/txproposals/%s", txID)
	bytes, err := c.doGetRequest(path, map[string]string{})
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetTxProposals returns pending transaction proposals
func (c *Client) GetTxProposals() ([]*models.TxProposal, error) {
	bytes, err := c.doGetRequest("/v1/txproposals/", map[string]string{})
	if err != nil {
		return nil, err
	}

	response := []*models.TxProposal{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetTxHistory returns completed transactions
func (c *Client) GetTxHistory(skip, limit uint64, includeExtendedInfo bool) ([]*models.Transaction, error) {
	params := map[string]string{
		"skip":                strconv.FormatUint(skip, 10),
		"limit":               strconv.FormatUint(limit, 10),
		"includeExtendedInfo": utils.BoolToString(includeExtendedInfo),
	}

	bytes, err := c.doGetRequest("/v1/txhistory/", params)
	if err != nil {
		return nil, err
	}

	response := []*models.Transaction{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateTxProposal creates transaction proposal
func (c *Client) CreateTxProposal(outputs []*models.TxOutput, feeLevel string, dryRun bool) (*models.TxProposal, error) {
	if len(feeLevel) == 0 {
		feeLevel = "normal"
	}

	payload := map[string]interface{}{
		"outputs":  outputs,
		"feeLevel": feeLevel,
		"dryRun":   dryRun,
	}

	bytes, err := c.doPostRequest("/v2/txproposals/", payload)
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// PublishTxProposal publishes transaction proposal
func (c *Client) PublishTxProposal(txp *models.TxProposal) (*models.TxProposal, error) {
	proposalSignature, err := txp.ProposalSignature(c.keys.ReqPrvKey, c.cfg.NetParams())
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"proposalSignature": utils.ToHex(proposalSignature),
	}

	path := fmt.Sprintf("/v1/txproposals/%s/publish/", txp.ID)
	bytes, err := c.doPostRequest(path, payload)
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// SignTxProposal signs transaction proposal
func (c *Client) SignTxProposal(txp *models.TxProposal) (*models.TxProposal, error) {
	signatures := []string{}

	for idx, input := range txp.Inputs {
		// Derive private key for input
		privKey, _, err := c.keys.DeriveFromAccount(input.Path)
		if err != nil {
			return nil, err
		}

		// Sign transaction
		signature, err := txp.InputSignature(privKey, c.cfg.NetParams(), idx)
		if err != nil {
			return nil, err
		}

		signatures = append(signatures, utils.ToHex(signature))
	}

	payload := map[string]interface{}{
		"signatures": signatures,
	}

	path := fmt.Sprintf("/v1/txproposals/%s/signatures/", txp.ID)
	bytes, err := c.doPostRequest(path, payload)
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// RejectTxProposal signs transaction proposal
func (c *Client) RejectTxProposal(txID, reason string) (*models.TxProposal, error) {
	payload := map[string]interface{}{
		"reason": reason,
	}

	path := fmt.Sprintf("/v1/txproposals/%s/rejections/", txID)
	bytes, err := c.doPostRequest(path, payload)
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// BroadcastRawTx sends raw transaction
func (c *Client) BroadcastRawTx(rawTx, network string) (*string, error) {
	if len(network) == 0 {
		network = c.cfg.Network
	}

	payload := map[string]interface{}{
		"rawTx":   rawTx,
		"network": network,
	}

	bytes, err := c.doPostRequest("/v1/broadcast_raw/", payload)
	if err != nil {
		return nil, err
	}

	var response string
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// BroadcastTxProposal signs transaction proposal
func (c *Client) BroadcastTxProposal(txID string) (*models.TxProposal, error) {
	path := fmt.Sprintf("/v1/txproposals/%s/broadcast/", txID)
	bytes, err := c.doPostRequest(path, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// RemoveTxProposal signs transaction proposal
func (c *Client) RemoveTxProposal(txID string) (*models.TxProposal, error) {
	path := fmt.Sprintf("/v1/txproposals/%s", txID)
	bytes, err := c.doDeleteRequest(path, map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	response := &models.TxProposal{}
	if err := json.Unmarshal(bytes, response); err != nil {
		return nil, err
	}

	return response, nil
}

// GetNotifications returns notifications
func (c *Client) GetNotifications(lastNotificationID string, timeSpan uint64, includeOwn bool) ([]*models.Notification, error) {
	payload := map[string]string{
		"includeOwn": utils.BoolToString(includeOwn),
	}

	if lastNotificationID != "" {
		payload["notificationId"] = lastNotificationID
	}

	if timeSpan != 0 {
		payload["timeSpan"] = strconv.FormatUint(timeSpan, 10)
	}

	bytes, err := c.doGetRequest("/v1/notifications/", payload)
	if err != nil {
		return nil, err
	}

	response := []*models.Notification{}
	if err := json.Unmarshal(bytes, &response); err != nil {
		return nil, err
	}

	return response, nil
}

// PushNotificationsSubscribe returns notifications
func (c *Client) PushNotificationsSubscribe(platform, token string) error {
	payload := map[string]interface{}{
		"type":  platform,
		"token": token,
	}

	if _, err := c.doPostRequest("/v1/pushnotifications/subscriptions/", payload); err != nil {
		return err
	}

	return nil
}

// PushNotificationsUnsubscribe returns notifications
func (c *Client) PushNotificationsUnsubscribe(token string) error {
	path := fmt.Sprintf("/v2/pushnotifications/subscriptions/%s", token)
	if _, err := c.doDeleteRequest(path, map[string]interface{}{}); err != nil {
		return err
	}

	return nil
}

// Performs GET requests
func (c *Client) doGetRequest(path string, params map[string]string) ([]byte, error) {
	// Prepare URL query
	values := url.Values{}
	for key, value := range params {
		values.Add(key, value)
	}

	// Prepare URL
	url := url.URL{}
	url.Path = path
	url.RawQuery = values.Encode()

	return c.doRequest("GET", url.String(), map[string]interface{}{})
}

// Performs POST requests
func (c *Client) doPostRequest(path string, payload map[string]interface{}) ([]byte, error) {
	return c.doRequest("POST", path, payload)
}

// Performs PUT requests
func (c *Client) doPutRequest(path string, payload map[string]interface{}) ([]byte, error) {
	return c.doRequest("PUT", path, payload)
}

// Performs DELETE requests
func (c *Client) doDeleteRequest(path string, payload map[string]interface{}) ([]byte, error) {
	return c.doRequest("DELETE", path, payload)
}

// Performs HTTP requests
func (c *Client) doRequest(method, path string, payload map[string]interface{}) ([]byte, error) {
	args, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	signature, err := c.signRequest(method, path, string(args))
	if err != nil {
		return nil, err
	}

	reqBody := &bytes.Reader{}
	if method != "GET" {
		reqBody = bytes.NewReader(args)
	}

	res, err := c.client.Do(method, c.absoluteURL(path), c.headers(signature), reqBody)
	if err != nil {
		return nil, err
	}

	return c.handleError(res)
}

// Handle errors
func (c *Client) handleError(res *httpclient.Response) ([]byte, error) {
	if res.StatusCode == 404 {
		return nil, fmt.Errorf("API error, status: %d", res.StatusCode)
	}

	bytes, err := res.ReadAll()
	if err != nil {
		return bytes, err
	}

	if res.StatusCode >= 400 && res.StatusCode <= 511 {
		body := &models.Error{}
		if err := json.Unmarshal(bytes, body); err != nil {
			return nil, err
		}

		err := fmt.Errorf("API error, status: %d, code: %s, message: %s, ", res.StatusCode, body.Code, body.Message)
		return nil, err
	}

	return bytes, nil
}

func (c *Client) absoluteURL(path string) string {
	return c.cfg.BaseURL + path
}

func (c *Client) copayerHash(name string) []byte {
	requestPubKey := hex.EncodeToString(c.keys.ReqPubKey.SerializeCompressed())
	return []byte(strings.Join([]string{name, c.keys.AccExtPubKey.String(), requestPubKey}, "|"))
}

func (c *Client) headers(signature []byte) map[string]string {
	return map[string]string{
		"x-client-version": clientVersion,
		"x-identity":       hex.EncodeToString(c.xPubToCopayerID()),
		"x-signature":      hex.EncodeToString(signature),
	}
}

func (c *Client) signRequest(method, url, args string) ([]byte, error) {
	message := strings.Join([]string{strings.ToLower(method), url, args}, "|")
	signature, err := utils.SignMessage([]byte(message), c.keys.ReqPrvKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

func (c *Client) xPubToCopayerID() []byte {
	base58 := c.keys.AccExtPubKey.String()
	data := []byte(base58)
	if c.cfg.Coin != config.CoinBTC {
		data = []byte(c.cfg.Coin + base58)
	}

	return utils.Sha256(data)
}
