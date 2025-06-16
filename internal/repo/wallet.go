package repo

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"time"
)

// ErrWalletNotFound indicates that the requested wallet was not found.
var ErrWalletNotFound = errors.New("wallet not found")

// ErrInsufficientFunds indicates that the wallet does not have enough balance for the operation.
var ErrInsufficientFunds = errors.New("insufficient funds")

type WalletRepoImpl struct {
	db *sqlx.DB
}

func NewWalletImpl(db *sqlx.DB) *WalletRepoImpl {
	return &WalletRepoImpl{db}
}

func (wr *WalletRepoImpl) GetWalletTransactions(ctx context.Context) ([]string, error) {
	return []string{"1234", "5678"}, nil
}

// GetWalletByUserIDAndWalletID retrieves a specific wallet for a user.
// userIDStr and walletIDStr are strings that will be parsed to uuid.UUID.
func (wr *WalletRepoImpl) RetrieveWalletByUserIdAndWalletId(ctx context.Context, userIDStr string, walletIDStr string) (*model.Wallet, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		return nil, errors.New("invalid wallet ID format")
	}

	var wallet model.Wallet
	query := `SELECT id, user_id, name, balance, created_at, updated_at
              FROM wallets
              WHERE user_id = $1 AND id = $2`

	err = wr.db.GetContext(ctx, &wallet, query, userID, walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err // Return other database-related errors
	}
	return &wallet, nil
}

// Deposit adds funds to a wallet and creates a transaction record.
func (wr *WalletRepoImpl) Deposit(ctx context.Context, userIDStr string, walletIDStr string, amount decimal.Decimal) (*model.Transaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}
	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID format: %w", err)
	}

	tx, err := wr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if not committed

	// 1. Retrieve and lock the wallet row
	var wallet model.Wallet
	queryWallet := `SELECT id, user_id, name, balance, created_at, updated_at
                    FROM wallets
                    WHERE user_id = $1 AND id = $2 FOR UPDATE`
	err = tx.GetContext(ctx, &wallet, queryWallet, userID, walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("failed to retrieve wallet for deposit: %w", err)
	}

	// 2. Update wallet balance
	newBalance := wallet.Balance.Add(amount)
	updateQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateQuery, newBalance, time.Now(), walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance for deposit: %w", err)
	}

	// 3. Create transaction record
	transaction := &model.Transaction{
		ID:        uuid.New(),
		WalletID:  walletID,
		Type:      model.TransactionTypeDeposit,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	insertTxQuery := `INSERT INTO transactions (id, wallet_id, type, amount, created_at)
                      VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, insertTxQuery, transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create deposit transaction record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit deposit transaction: %w", err)
	}

	return transaction, nil
}
