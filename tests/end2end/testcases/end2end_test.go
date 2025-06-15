package testcases

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"wallet-app-server/app/model"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

const (
	ApiRoot = "https://localhost:8227/api/v1"
)

// Config http client to bypass SSL verification
func init() {
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
}

/*
Test case 1 (Happy Path):
1. Login as vence.lin user
2. List wallets (expected 2 wallets returned)
3. Deposit 10000.00 to wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (balance 10000.00)
4. Withdraw 3000.00 from wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (balance 7000.00)
5. Transfer 2000.00 from wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d to wallet 68e95347-29ad-4324-9725-eed1feaa8594 (balance 5000.00)
6. Check current balance for wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (expect 5000.00)
5. List transaction history for wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (expect 3 records returned)
*/
func TestHappyPath(t *testing.T) {
	// Test login
	accessToken, err := testLogin(t, "vence.lin", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken: %s", accessToken)

	// Test list wallets
	wallets, err := testListWallets(t, accessToken)
	assert.NoError(t, err, "Failed to list wallets")
	assert.Equal(t, len(wallets), 2)

	walletID := "a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d"

	// Test deposit
	newBalance, err := testDeposit(t, accessToken, walletID, decimal.NewFromFloat(10000.00))
	assert.NoError(t, err, "Failed to deposit")
	assert.Equal(t, newBalance, decimal.NewFromFloat(10000.00))

	// Test withdraw
	newBalance, err = testWithdraw(t, accessToken, walletID, decimal.NewFromFloat(3000.00))
	assert.NoError(t, err, "Failed to withdraw")
	assert.Equal(t, newBalance, decimal.NewFromFloat(7000.00))

	// Test transfer
	toWalletID := "68e95347-29ad-4324-9725-eed1feaa8594"
	txnID, err := testTransfer(t, accessToken, walletID, toWalletID, decimal.NewFromFloat(2000.00))
	assert.NoError(t, err, "Failed to transfer")
	assert.NotEmpty(t, txnID, "Transaction ID should not be empty")
	t.Logf("txnID: %s", txnID)

	// Test check balance
	latestBalance, err := testCheckBalance(t, accessToken, walletID)
	assert.NoError(t, err, "Failed to check balance")
	assert.Equal(t, latestBalance, decimal.NewFromFloat(5000.00))
}

/*
Test case 2 (Withdraw with Insufficient Amount):
1. Login as vence.lin user
2. Deposit 10000.00 to wallet 34fad474-1df7-40a1-8675-0af586d02435 (balance 10000.00)
3. Withdraw 12000.00 from wallet 34fad474-1df7-40a1-8675-0af586d02435 (expect insufficient amount error)
*/
func TestWithdrawWithInsufficientAmount(t *testing.T) {
	// Test login
	accessToken, err := testLogin(t, "vence.lin", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken: %s", accessToken)

	walletID := "34fad474-1df7-40a1-8675-0af586d02435"

	// Test deposit
	newBalance, err := testDeposit(t, accessToken, walletID, decimal.NewFromFloat(10000.00))
	assert.NoError(t, err, "Failed to deposit")
	assert.Equal(t, newBalance, decimal.NewFromFloat(10000.00))

	// Test withdraw
	_, err = testWithdraw(t, accessToken, walletID, decimal.NewFromFloat(12000.00))
	assert.Error(t, err, "Should have error")
	t.Logf("withdraw err: %s", err.Error())
}

/*
Test case 3 (Transfer with Insufficient Amount, check both wallet balance are not changed)
 1. Login as vence.lin user
 2. Transfer 7000.00 from wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (current balance 5000.00) to wallet 68e95347-29ad-4324-9725-eed1feaa8594 (current wallet 2000.00)
    (expect insufficient amount error)
 3. Check balance of wallet a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d (expect 5000.00)
 4. Login as nick.lee user
 5. Check balance of wallet 68e95347-29ad-4324-9725-eed1feaa8594 (expect 2000.00)
*/
func TestTransferWithInsufficientAmount(t *testing.T) {
	// Test login
	accessToken, err := testLogin(t, "vence.lin", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken(vence.lin): %s", accessToken)

	walletID := "a5344dde-a6a2-4c7a-8b9d-78841ef0ab3d"
	toWalletID := "68e95347-29ad-4324-9725-eed1feaa8594"

	// Test transfer
	_, err = testTransfer(t, accessToken, walletID, toWalletID, decimal.NewFromFloat(7000.00))
	assert.Error(t, err, "Should have error")
	t.Logf("transfer err: %s", err.Error())

	// Test check balance
	latestBalance, err := testCheckBalance(t, accessToken, walletID)
	assert.NoError(t, err, "Failed to check balance")
	assert.Equal(t, latestBalance, decimal.NewFromFloat(5000.00))

	// Test login
	accessToken, err = testLogin(t, "nick.lee", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken(nick.lee): %s", accessToken)

	// Test check balance
	latestBalance, err = testCheckBalance(t, accessToken, toWalletID)
	assert.NoError(t, err, "Failed to check balance")
	assert.Equal(t, latestBalance, decimal.NewFromFloat(2000.00))
}

/*
Test case 4 (Deposit to a wallet not related current user)
 1. Login as vence.lin user
 2. Deposit 1000.00 to wallet 68e95347-29ad-4324-9725-eed1feaa8594 (expect invalid wallet ID error)
*/
func TestDepositToWalletNotSelf(t *testing.T) {
	// Test login
	accessToken, err := testLogin(t, "vence.lin", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken: %s", accessToken)

	walletID := "68e95347-29ad-4324-9725-eed1feaa8594"

	// Test withdraw
	_, err = testDeposit(t, accessToken, walletID, decimal.NewFromFloat(1000.00))
	assert.Error(t, err, "Should have error")
	t.Logf("withdraw err: %s", err.Error())
}

/*
Test case 5 (Withdraw from a wallet not related current user)
 1. Login as vence.lin user
 2. Withdraw 1000.00 from wallet 68e95347-29ad-4324-9725-eed1feaa8594 (expect invalid wallet ID error)
*/
func TestWithdrawFromWalletNotSelf(t *testing.T) {
	// Test login
	accessToken, err := testLogin(t, "vence.lin", "P@ssw0rd")
	assert.NoError(t, err, "Failed to login")
	assert.NotEmpty(t, accessToken, "Access token should not be empty")
	t.Logf("accessToken: %s", accessToken)

	walletID := "68e95347-29ad-4324-9725-eed1feaa8594"

	// Test withdraw
	_, err = testWithdraw(t, accessToken, walletID, decimal.NewFromFloat(1000.00))
	assert.Error(t, err, "Should have error")
	t.Logf("withdraw err: %s", err.Error())
}

func testLogin(t *testing.T, username string, password string) (string, error) {
	// Construct request body
	reqBody := map[string]any{
		"username": username,
		"password": password,
	}
	body, _ := json.Marshal(reqBody)
	t.Logf("[testLogin] --> %s", string(body))
	req, err := http.NewRequest("POST", ApiRoot+"/user/login", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	t.Logf("[testLogin] <-- %s", string(respBody))
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	if response["success"] == false {
		return "", errors.New(response["error"].(string))
	}
	return response["access_token"].(string), nil
}

func testListWallets(t *testing.T, accessToken string) ([]model.WalletInfo, error) {
	t.Logf("[testListWallets] --> none")
	req, err := http.NewRequest("GET", ApiRoot+"/wallet/list", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	t.Logf("[testListWallets] <-- %s", string(respBody))
	var response = struct {
		Success bool               `json:"success"`
		Wallets []model.WalletInfo `json:"wallets"`
		Error   string             `json:"error"`
	}{}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return nil, err
	}
	if response.Success == false {
		return nil, errors.New(response.Error)
	}
	return response.Wallets, nil
}

func testDeposit(t *testing.T, accessToken string, walletID string, amount decimal.Decimal) (decimal.Decimal, error) {
	reqBody := map[string]any{
		"wallet_id": walletID,
		"amount":    amount,
	}
	body, _ := json.Marshal(reqBody)
	t.Logf("[testDeposit] --> %s", string(body))
	req, err := http.NewRequest("POST", ApiRoot+"/wallet/deposit", bytes.NewBuffer(body))
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}
	t.Logf("[testDeposit] <-- %s", string(respBody))
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return decimal.Zero, err
	}
	if response["success"] == false {
		return decimal.Zero, errors.New(response["error"].(string))
	}
	return decimal.NewFromFloat(response["balance"].(float64)), nil
}

func testWithdraw(t *testing.T, accessToken string, walletID string, amount decimal.Decimal) (decimal.Decimal, error) {
	reqBody := map[string]any{
		"wallet_id": walletID,
		"amount":    amount,
	}
	body, _ := json.Marshal(reqBody)
	t.Logf("[testWithdraw] --> %s", string(body))
	req, err := http.NewRequest("POST", ApiRoot+"/wallet/withdraw", bytes.NewBuffer(body))
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}
	t.Logf("[testWithdraw] <-- %s", string(respBody))
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return decimal.Zero, err
	}
	if response["success"] == false {
		return decimal.Zero, errors.New(response["error"].(string))
	}
	return decimal.NewFromFloat(response["balance"].(float64)), nil
}

func testTransfer(t *testing.T, accessToken string, fromWalletID string, toWalletID string, amount decimal.Decimal) (string, error) {
	reqBody := map[string]any{
		"from_wallet_id": fromWalletID,
		"to_wallet_id":   toWalletID,
		"amount":         amount,
	}
	body, _ := json.Marshal(reqBody)
	t.Logf("[testTransfer] --> %s", string(body))
	req, err := http.NewRequest("POST", ApiRoot+"/transaction/transfer", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	t.Logf("[testTransfer] <-- %s", string(respBody))
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return "", err
	}
	if response["success"] == false {
		return "", errors.New(response["error"].(string))
	}
	return response["txn_id"].(string), nil
}

func testCheckBalance(t *testing.T, accessToken string, walletID string) (decimal.Decimal, error) {
	// Construct request body
	reqBody := map[string]any{
		"wallet_id": walletID,
	}
	body, _ := json.Marshal(reqBody)
	t.Logf("[testCheckBalance] --> %s", string(body))
	req, err := http.NewRequest("POST", ApiRoot+"/wallet/checkBalance", bytes.NewBuffer(body))
	if err != nil {
		return decimal.Zero, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return decimal.Zero, err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return decimal.Zero, err
	}
	t.Logf("[testCheckBalance] <-- %s", string(respBody))
	var response map[string]any
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		return decimal.Zero, err
	}
	if response["success"] == false {
		return decimal.Zero, errors.New(response["error"].(string))
	}
	return decimal.NewFromFloat(response["balance"].(float64)), nil
}
