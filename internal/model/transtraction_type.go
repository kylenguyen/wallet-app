package model

// TransactionType defines the allowed types for a transaction.
type TransactionType string

const (
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeTransfer   TransactionType = "transfer"
)

// IsValid checks if the transaction type is valid.
func (tt TransactionType) IsValid() bool {
	switch tt {
	case TransactionTypeDeposit, TransactionTypeWithdrawal, TransactionTypeTransfer:
		return true
	}
	return false
}
