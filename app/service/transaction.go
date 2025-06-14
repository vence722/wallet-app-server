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
	// Verify from wallet is belong to the current user
	valid, err := repository.WalletRepository.VerifyUserWalletPossession(db.DB, currentUserID, fromWalletID)
	if err != nil {
		return "", http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	if !valid {
		return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
	}
	var result string
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		// Record current time
		currTime := time.Now()
		// Transfer money
		if err := repository.WalletRepository.Transfer(tx, currentUserID, fromWalletID, toWalletID, amount); err != nil {
			return err
		}
		// Create transaction history
		txnID, err := repository.TransactionRepository.CreateTransactionHistory(tx, fromWalletID, toWalletID, constant.TxnTypeTransfer, amount, currTime)
		if err != nil {
			return err
		}
		result = txnID
		// Create user activity
		activityDetail := fmt.Sprintf("User transfer amount %s from wallet %s to wallet %s", amount.StringFixed(2), fromWalletID, toWalletID)
		if err := repository.UserRepository.CreateUserActivity(tx, currentUserID, constant.UserActTypeTransfer, activityDetail, fromWalletID, currTime); err != nil {
			return err
		}
		return nil
	}); err != nil {
		// If the underlying error is business logic related error
		// return bad request status code
		// otherwise return internal server error status code
		if err.Error() == repository.ErrNegativeOrZeroAmount {
			return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageNegativeOrZeroAmount, nil)
		}
		if err.Error() == repository.ErrInsufficientBalance {
			return "", http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageInsufficientBalance, nil)
		}
		return "", http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Return txnID as result, and success status code
	return result, http.StatusOK, nil
}

func (ts *transactionServiceImpl) ListHistory(currentUserID string, walletID string) ([]model.TransactionHistory, int, error) {
	// Verify from wallet is belong to the current user
	valid, err := repository.WalletRepository.VerifyUserWalletPossession(db.DB, currentUserID, walletID)
	if err != nil {
		return nil, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	if !valid {
		return nil, http.StatusBadRequest, newServiceError(ErrTypeInvalidRequestBody, ErrMessageWalletIDInvalid, nil)
	}
	// List transaction history
	txnHistoryList, err := repository.TransactionRepository.ListTransactionHistory(db.DB, walletID)
	if err != nil {
		return nil, http.StatusInternalServerError, newServiceError(ErrTypeInternalServerError, ErrMessageDBError, err)
	}
	// Construct result list of model.TransactionHistory
	var result = make([]model.TransactionHistory, 0, len(txnHistoryList))
	for _, txnHistory := range txnHistoryList {
		result = append(result, model.TransactionHistory{
			TxnID:        txnHistory.TxnID,
			FromWalletID: txnHistory.FromWalletID,
			ToWalletID:   txnHistory.ToWalletID,
			TxnAmount:    txnHistory.TxnAmount,
			TxnTypeDesc:  txnHistory.TxnType,
			TxnTime:      txnHistory.TxnTime,
		})
	}
	return result, http.StatusOK, nil
}
