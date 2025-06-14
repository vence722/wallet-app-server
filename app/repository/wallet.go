package repository

import (
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Wallet repository interface
type IWalletRepository interface {
	Transfer(db *gorm.DB, userID string, fromWalletID string, toWalletID string, amount decimal.Decimal) error
}

// Wallet repository instance
var WalletRepository IWalletRepository = &walletRepositoryImpl{}

// Wallet repository implementation
type walletRepositoryImpl struct{}

// Transfer money from a wallet to another
// Should call this method inside a transaction
// Note that the wallet rows will be locked during the transaction to achieve consistency
func (wr *walletRepositoryImpl) Transfer(tx *gorm.DB, userID string, fromWalletID string, toWalletID string, amount decimal.Decimal) error {
	// Ensure transaction amount > 0
	if amount.IsNegative() || amount.IsZero() {
		return errors.New(ErrNegativeOrZeroAmount)
	}
	// Fetch from wallet balance
	// [NOTE] use clause Strengh = "UPDATE" to implement SELECT ... FOR UPDATE in PostgreSQL
	var fromWalletBalance decimal.Decimal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", fromWalletID).Select("balance").Find(&fromWalletBalance).Error; err != nil {
		return err
	}
	// Check balance sufficiency
	if fromWalletBalance.Cmp(amount) < 0 {
		return errors.New(ErrInsufficientBalance)
	}
	// Fetch to wallet balance
	// [NOTE] use clause Strengh = "UPDATE" to implement SELECT ... FOR UPDATE in PostgreSQL
	var toWalletBalance decimal.Decimal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", toWalletID).Select("balance").Find(&toWalletBalance).Error; err != nil {
		return err
	}
	// Modify from wallet balance (- amount)
	newFromWalletBalance := fromWalletBalance.Sub(amount)
	if err := tx.Table("wallet").Where("wallet_id = ?", fromWalletID).Update("balance", newFromWalletBalance).Error; err != nil {
		return err
	}
	// Modify to wallet balance (+ amount)
	newToWalletBalance := toWalletBalance.Add(amount)
	if err := tx.Table("wallet").Where("wallet_id = ?", toWalletID).Update("balance", newToWalletBalance).Error; err != nil {
		return err
	}
	return nil
}
