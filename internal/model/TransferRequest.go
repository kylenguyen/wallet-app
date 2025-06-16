package model

import "github.com/shopspring/decimal"

// TransferRequest is the request body for transferring funds.
type TransferRequest struct {
	Amount              decimal.Decimal `json:"amount" binding:"required"`
	DestinationWalletID string          `json:"destination_wallet_id" binding:"required"`
}
