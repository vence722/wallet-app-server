package entity

import (
	"database/sql"
	"time"

	"github.com/shopspring/decimal"
)

type User struct {
	UserID     string       `gorm:"primaryKey;column:user_id"`
	UserName   string       `gorm:"column:user_name"`
	UserHash   string       `gorm:"column:user_hash"`
	CreateTime time.Time    `gorm:"column:create_time"`
	UpdateTime sql.NullTime `gorm:"column:update_time"`
}

func (u *User) TableName() string {
	return "user"
}

type Wallet struct {
	WalletID   string          `gorm:"primaryKey;column:wallet_id"`
	WalletName string          `gorm:"column:wallet_name"`
	Balance    decimal.Decimal `gorm:"column:balance"`
	CreateTime time.Time       `gorm:"column:create_time"`
	UpdateTime sql.NullTime    `gorm:"column:update_time"`
}

func (w *Wallet) TableName() string {
	return "wallet"
}

type TxnHistory struct {
	TxnID        string          `gorm:"primaryKey;column:txn_id"`
	FromWalletID string          `gorm:"column:from_wallet_id"`
	ToWalletID   string          `gorm:"column:to_wallet_id"`
	TxnType      string          `gorm:"column:txn_type"`
	TxnAmount    decimal.Decimal `gorm:"column:txn_amount"`
	TxnTime      time.Time       `gorm:"column:txn_time"`
}

func (th *TxnHistory) TableName() string {
	return "txn_history"
}

type UserActivity struct {
	UserActID     string    `gorm:"primaryKey;column:user_act_id"`
	UserActType   string    `gorm:"column:user_act_type"`
	UserActDetail string    `gorm:"column:user_act_detail"`
	UserWalletID  string    `gorm:"column:user_wallet_id"`
	UserActTime   time.Time `gorm:"column:user_act_time"`
}

func (ua *UserActivity) TableName() string {
	return "user_activity"
}
