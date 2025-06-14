package service

import (
	"fmt"
	"net/http"
	"time"
	"wallet-app-server/app/constant"
	"wallet-app-server/app/db"
	"wallet-app-server/app/model"
	"wallet-app-server/app/repository"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Wallet service interface
type IWalletService interface {
	ListUserWallets(currentUserID string) ([]model.WalletInfo, int, error)
	CheckWalletBallance(currentUserID string, walletID string) (decimal.Decimal, int, error)
	Deposit(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error)
	Withdraw(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error)
}

// Wallet service instance
var WalletService IWalletService = &walletServiceImpl{}

// Wallet service implementation
type walletServiceImpl struct{}

func (ws *walletServiceImpl) ListUserWallets(currentUserID string) ([]model.WalletInfo, int, error) {
	wallets, err := repository.WalletRepository.ListUserWallets(db.DB, currentUserID)
	if err != nil {
		// If record not found, return empty OK response
		if err == gorm.ErrRecordNotFound {
			return []model.WalletInfo{}, http.StatusOK, nil
		}
		// Other repository error
		return []model.WalletInfo{}, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Construct result list of model.WalletInfo
	result := make([]model.WalletInfo, 0, len(wallets))
	for _, wallet := range wallets {
		result = append(result, model.WalletInfo{
			WalletID:   wallet.WalletID,
			WalletName: wallet.WalletName,
		})
	}
	return result, http.StatusOK, nil
}

func (ws *walletServiceImpl) CheckWalletBallance(currentUserID string, walletID string) (decimal.Decimal, int, error) {
	// Verify from wallet is belong to the current user
	valid, err := repository.WalletRepository.VerifyUserWalletPossession(db.DB, currentUserID, walletID)
	if err != nil {
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	if !valid {
		return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
	}
	wallet, err := repository.WalletRepository.GetWalletByID(db.DB, walletID)
	if err != nil {
		// Wallet not found
		if err == gorm.ErrRecordNotFound {
			return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
		}
		// Other repository error
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Return wallet.Balance
	return wallet.Balance, http.StatusOK, nil
}

func (ws *walletServiceImpl) Deposit(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error) {
	// Verify from wallet is belong to the current user
	valid, err := repository.WalletRepository.VerifyUserWalletPossession(db.DB, currentUserID, walletID)
	if err != nil {
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	if !valid {
		return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
	}
	var result decimal.Decimal
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Record current time
		currTime := time.Now()
		// Deposit
		latestBalance, err := repository.WalletRepository.Deposit(tx, walletID, amount)
		if err != nil {
			return err
		}
		result = latestBalance
		// Create transaction history
		if _, err := repository.TransactionRepository.CreateTransactionHistory(tx, walletID, walletID, constant.TxnTypeDeposit, amount, currTime); err != nil {
			return err
		}
		// Create user activity
		activityDetail := fmt.Sprintf("User deposit amount %s to wallet %s", amount.StringFixed(2), walletID)
		if err := repository.UserRepository.CreateUserActivity(tx, constant.UserActTypeDeposit, activityDetail, walletID, currTime); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// If the underlying error is business logic related error
		// return bad request status code
		// otherwise return internal server error status code
		if err.Error() == repository.ErrNegativeOrZeroAmount {
			return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageNegativeOrZeroAmount, nil)
		}
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	return result, http.StatusOK, nil
}

func (ws *walletServiceImpl) Withdraw(currentUserID string, walletID string, amount decimal.Decimal) (decimal.Decimal, int, error) {
	// Verify from wallet is belong to the current user
	valid, err := repository.WalletRepository.VerifyUserWalletPossession(db.DB, currentUserID, walletID)
	if err != nil {
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	if !valid {
		return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
	}
	var result decimal.Decimal
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Record current time
		currTime := time.Now()
		// Withdraw
		latestBalance, err := repository.WalletRepository.Withdraw(tx, walletID, amount)
		if err != nil {
			return err
		}
		result = latestBalance
		// Create transaction history
		if _, err := repository.TransactionRepository.CreateTransactionHistory(tx, walletID, walletID, constant.TxnTypeWithdraw, amount, currTime); err != nil {
			return err
		}
		// Create user activity
		activityDetail := fmt.Sprintf("User withdraw amount %s to wallet %s", amount.StringFixed(2), walletID)
		if err := repository.UserRepository.CreateUserActivity(tx, constant.UserActTypeWithdraw, activityDetail, walletID, currTime); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// If the underlying error is business logic related error
		// return bad request status code
		// otherwise return internal server error status code
		if err.Error() == repository.ErrNegativeOrZeroAmount {
			return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageNegativeOrZeroAmount, nil)
		}
		if err.Error() == repository.ErrInsufficientBalance {
			return decimal.Zero, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageInsufficientBalance, nil)
		}
		return decimal.Zero, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	return result, http.StatusOK, nil
}
