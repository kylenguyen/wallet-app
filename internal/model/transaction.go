package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// Transaction represents the structure of the 'transactions' table.
type Transaction struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	WalletID        uuid.UUID       `json:"wallet_id" db:"wallet_id"`
	Type            TransactionType `json:"type" db:"type"`
	Amount          decimal.Decimal `json:"amount" db:"amount"`
	RelatedWalletID *uuid.UUID      `json:"related_wallet_id,omitempty" db:"related_wallet_id"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
}
