package model

import (
	"time"

	"github.com/shopspring/decimal"
)

type TransactionHistory struct {
	TxnID        string          `json:"txn_id"`
	FromWalletID string          `json:"from_wallet_id"`
	ToWalletID   string          `json:"to_wallet_id"`
	TxnAmount    decimal.Decimal `json:"txn_amount"`
	TxnTypeDesc  string          `json:"txn_type_desc"`
	TxnTime      time.Time       `json:"txn_time"`
}
