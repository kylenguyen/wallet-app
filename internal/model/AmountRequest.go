package model

import "github.com/shopspring/decimal"

type AmountRequest struct {
	Amount decimal.Decimal `json:"amount" binding:"required"`
}
