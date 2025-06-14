package service

import (
	"net/http"
	"wallet-app-server/app/model"

	"github.com/shopspring/decimal"
)

// Transaction service interface
type ITransactionService interface {
	Transfer(currentUserID string, fromWalletID string, toWalletID string, amount decimal.Decimal) (string, int, error)
	ListHistory(currentUserID string, walletID string) ([]model.TransactionHistory, int, error)
}

// Transaction service instance
var TransactionService ITransactionService = &transactionServiceImpl{}

// Transaction service implementation
type transactionServiceImpl struct{}

func (ts *transactionServiceImpl) Transfer(currentUserID string, fromWalletID string, toWalletID string, amount decimal.Decimal) (string, int, error) {
	return "", http.StatusOK, nil
}

func (ts *transactionServiceImpl) ListHistory(currentUserID string, walletID string) ([]model.TransactionHistory, int, error) {
	return []model.TransactionHistory{}, http.StatusOK, nil
}
