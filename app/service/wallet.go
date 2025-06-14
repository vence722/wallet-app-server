package service

import (
	"net/http"
	"wallet-app-server/app/model"

	"github.com/shopspring/decimal"
)

// Wallet service interface
type IWalletService interface {
	ListWallets(currentUserID string) ([]model.WalletInfo, int, error)
	CheckWalletBallance(currentUserID string, walletID string) (decimal.Decimal, int, error)
	Deposit(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error)
	Withdraw(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error)
}

// Wallet service instance
var WalletService IWalletService = &walletServiceImpl{}

// Wallet service implementation
type walletServiceImpl struct{}

func (ws *walletServiceImpl) ListWallets(currentUserID string) ([]model.WalletInfo, int, error) {
	return []model.WalletInfo{}, http.StatusOK, nil
}

func (ws *walletServiceImpl) CheckWalletBallance(currentUserID string, walletID string) (decimal.Decimal, int, error) {
	return decimal.Decimal{}, http.StatusOK, nil
}

func (ws *walletServiceImpl) Deposit(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error) {
	return decimal.Decimal{}, http.StatusOK, nil
}

func (ws *walletServiceImpl) Withdraw(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error) {
	return decimal.Decimal{}, http.StatusOK, nil
}
