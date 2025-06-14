package repository

import (
	"errors"
	"wallet-app-server/app/entity"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Wallet repository interface
type IWalletRepository interface {
	VerifyUserWalletPossession(db *gorm.DB, userID string, walletID string) (bool, error)
	ListUserWallets(db *gorm.DB, userID string) ([]entity.Wallet, error)
	GetWalletByID(db *gorm.DB, walletID string) (entity.Wallet, error)
	Deposit(db *gorm.DB, walletID string, amount decimal.Decimal) (decimal.Decimal, error)
	Withdraw(db *gorm.DB, walletID string, amount decimal.Decimal) (decimal.Decimal, error)
	Transfer(db *gorm.DB, userID string, fromWalletID string, toWalletID string, amount decimal.Decimal) error
}

// Wallet repository instance
var WalletRepository IWalletRepository = &walletRepositoryImpl{}

// Wallet repository implementation
type walletRepositoryImpl struct{}

// Verify if the wallet is belong to the user
func (wr *walletRepositoryImpl) VerifyUserWalletPossession(db *gorm.DB, userID string, walletID string) (bool, error) {
	var count int64
	if err := db.Table("user_wallet_bridge").Where("user_id = ? and wallet_id = ?", userID, walletID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// List user's wallets
func (wr *walletRepositoryImpl) ListUserWallets(db *gorm.DB, userID string) ([]entity.Wallet, error) {
	var wallets []entity.Wallet
	if err := db.Table("wallet").Joins("INNER JOIN user_wallet_bridge ON wallet.wallet_id = user_wallet_bridge.wallet_id").
		Where("user_wallet_bridge.user_id = ?", userID).Select("wallet.*").Order("user_wallet_bridge.seq").Find(&wallets).Error; err != nil {
		return []entity.Wallet{}, err
	}
	return wallets, nil
}

// Get wallet by ID
func (wr *walletRepositoryImpl) GetWalletByID(db *gorm.DB, walletID string) (entity.Wallet, error) {
	var wallet entity.Wallet
	if err := db.Table("wallet").Where("wallet_id = ?", walletID).First(&wallet).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.Wallet{}, nil
		}
		return entity.Wallet{}, err
	}
	return wallet, nil
}

// Deposit to wallet
// Should call this method inside a transaction
// Note that the wallet row will be locked during the transaction to achieve consistency
func (wr *walletRepositoryImpl) Deposit(tx *gorm.DB, walletID string, amount decimal.Decimal) (decimal.Decimal, error) {
	// Ensure transaction amount > 0
	if amount.IsNegative() || amount.IsZero() {
		return decimal.Zero, errors.New(ErrNegativeOrZeroAmount)
	}
	// Fetch wallet balance
	// [NOTE] use clause Strengh = "UPDATE" to implement SELECT ... FOR UPDATE in PostgreSQL
	var walletBalance decimal.Decimal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", walletID).Select("balance").Scan(&walletBalance).Error; err != nil {
		return decimal.Zero, err
	}
	// Modify to wallet balance (+ amount)
	newWalletBalance := walletBalance.Add(amount)
	if err := tx.Table("wallet").Where("wallet_id = ?", walletID).Update("balance", newWalletBalance).Error; err != nil {
		return decimal.Zero, err
	}
	return newWalletBalance, nil
}

// Withdraw from wallet
// Should call this method inside a transaction
// Note that the wallet row will be locked during the transaction to achieve consistency
func (wr *walletRepositoryImpl) Withdraw(tx *gorm.DB, walletID string, amount decimal.Decimal) (decimal.Decimal, error) {
	// Ensure transaction amount > 0
	if amount.IsNegative() || amount.IsZero() {
		return decimal.Zero, errors.New(ErrNegativeOrZeroAmount)
	}
	// Fetch wallet balance
	// [NOTE] use clause Strengh = "UPDATE" to implement SELECT ... FOR UPDATE in PostgreSQL
	var walletBalance decimal.Decimal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", walletID).Select("balance").Scan(&walletBalance).Error; err != nil {
		return decimal.Zero, err
	}
	// Check balance sufficiency
	if walletBalance.Cmp(amount) < 0 {
		return decimal.Zero, errors.New(ErrInsufficientBalance)
	}
	// Modify from wallet balance (- amount)
	newWalletBalance := walletBalance.Sub(amount)
	if err := tx.Table("wallet").Where("wallet_id = ?", walletID).Update("balance", newWalletBalance).Error; err != nil {
		return decimal.Zero, err
	}
	return newWalletBalance, nil
}

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
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", fromWalletID).Select("balance").Scan(&fromWalletBalance).Error; err != nil {
		return err
	}
	// Check balance sufficiency
	if fromWalletBalance.Cmp(amount) < 0 {
		return errors.New(ErrInsufficientBalance)
	}
	// Fetch to wallet balance
	// [NOTE] use clause Strengh = "UPDATE" to implement SELECT ... FOR UPDATE in PostgreSQL
	var toWalletBalance decimal.Decimal
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Table("wallet").Where("wallet_id = ?", toWalletID).Select("balance").Scan(&toWalletBalance).Error; err != nil {
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
