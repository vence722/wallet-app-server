package controller

import (
	"net/http"
	"wallet-app-server/app/service"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// List user's wallets
// GET /wallet/list
func ListWallets(c *gin.Context) {
	// Get current user ID
	currentUserID := c.GetString("current_user_id")

	// List wallets
	wallets, statusCode, err := service.WalletService.ListWallets(currentUserID)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"wallets": wallets})
}

// Check wallet's balance
// POST /wallet/checkBalance
func CheckWalletBalance(c *gin.Context) {
	// Get current user ID
	currentUserID := c.GetString("current_user_id")

	// Parse request body
	req := struct {
		WalletID string `json:"wallet_id"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		respondeWithError(c, http.StatusBadRequest, err)
		return
	}

	// Check wallet balance
	balance, statusCode, err := service.WalletService.CheckWalletBallance(currentUserID, req.WalletID)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"balance": balance})
}

// Deposit to user's wallet and return the latest balance
// POST /wallet/deposit
func Deposit(c *gin.Context) {
	// Get current user ID
	currentUserID := c.GetString("current_user_id")

	// Parse request body
	req := struct {
		WalletID string          `json:"wallet_id"`
		Amount   decimal.Decimal `json:"amount"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		respondeWithError(c, http.StatusBadRequest, err)
		return
	}

	// Deposit to user wallet
	latestBalance, statusCode, err := service.WalletService.Deposit(currentUserID, req.WalletID, req.Amount)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"balance": latestBalance})
}

// Withdraw from user's wallet and return the latest balance
// POST /wallet/withdraw
func Withdraw(c *gin.Context) {
	// Get current user ID
	currentUserID := c.GetString("current_user_id")

	// Parse request body
	req := struct {
		WalletID string          `json:"wallet_id"`
		Amount   decimal.Decimal `json:"amount"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		respondeWithError(c, http.StatusBadRequest, err)
		return
	}

	// Withdraw from user wallet
	latestBalance, statusCode, err := service.WalletService.Withdraw(currentUserID, req.WalletID, req.Amount)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"balance": latestBalance})
}
