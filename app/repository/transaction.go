package repository

import (
	"time"
	"wallet-app-server/app/entity"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

// Transaction repository interface
type ITransactionRepository interface {
	ListTransactionHistory(db *gorm.DB, walletID string) ([]entity.TxnHistory, error)
	CreateTransactionHistory(db *gorm.DB, fromWalletID string, toWalletID string, txnType string, txnAmount decimal.Decimal, txnTime time.Time) (string, error)
}

// Transaction repository instance
var TransactionRepository ITransactionRepository = &transactionRepositoryImpl{}

// Transaction repository implementation
type transactionRepositoryImpl struct{}

// List transaction history by wallet ID
func (tr *transactionRepositoryImpl) ListTransactionHistory(db *gorm.DB, walletID string) ([]entity.TxnHistory, error) {
	var result []entity.TxnHistory
	err := db.Where("from_wallet_id = ? or to_wallet_id = ?", walletID, walletID).Find(&result).Error
	return result, err
}

// Create new transaction history
func (tr *transactionRepositoryImpl) CreateTransactionHistory(db *gorm.DB, fromWalletID string, toWalletID string, txnType string, txnAmount decimal.Decimal, txnTime time.Time) (string, error) {
	txnID := uuid.New().String()
	txnHistory := entity.TxnHistory{
		TxnID:        txnID,
		FromWalletID: fromWalletID,
		ToWalletID:   toWalletID,
		TxnType:      txnType,
		TxnAmount:    txnAmount,
		TxnTime:      txnTime,
	}
	if err := db.Create(txnHistory).Error; err != nil {
		return "", err
	}
	return txnID, nil
}
