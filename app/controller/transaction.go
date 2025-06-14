package controller

import (
	"net/http"
	"wallet-app-server/app/service"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

// User transfer money from user's wallet to another wallet
// POST /transaction/transfer
func Transfer(c *gin.Context) {
	// Get current user ID
	currentUserID := c.GetString("current_user_id")

	// Parse request body
	req := struct {
		FromWalletID string          `json:"from_wallet_id"`
		ToWalletID   string          `json:"to_wallet_id"`
		Amount       decimal.Decimal `json:"amount"`
	}{}
	if err := c.BindJSON(&req); err != nil {
		respondeWithError(c, http.StatusBadRequest, err)
		return
	}

	// Make transfer
	txnID, statusCode, err := service.TransactionService.Transfer(currentUserID, req.FromWalletID, req.ToWalletID, req.Amount)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"txn_id": txnID})
}

// List user wallet's transaction history
// POST /transaction/history
func History(c *gin.Context) {
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

	// List transaction history
	txnHistory, statusCode, err := service.TransactionService.ListHistory(currentUserID, req.WalletID)
	if err != nil {
		respondeWithError(c, statusCode, err)
		return
	}

	// Return resposne
	resposneWithData(c, gin.H{"txn_history": txnHistory})
}
