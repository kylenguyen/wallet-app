package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kylenguyen/wallet-app/internal/model"
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

// GetTransactionsByWalletID retrieves all transactions for a specific wallet.
// It also checks if the wallet belongs to the given user for authorization.
func (wr *WalletRepoImpl) GetTransactionsByWalletID(ctx context.Context, userIDStr string, walletIDStr string) ([]model.Transaction, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID format: %w", err)
	}

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid wallet ID format: %w", err)
	}

	// First, verify the wallet exists and belongs to the user
	var walletExists bool
	checkWalletQuery := `SELECT EXISTS(SELECT 1 FROM wallets WHERE id = $1 AND user_id = $2)`
	err = wr.db.GetContext(ctx, &walletExists, checkWalletQuery, walletID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check wallet existence: %w", err)
	}
	if !walletExists {
		return nil, ErrWalletNotFound // Or a more specific "access denied" / "wallet not found for user"
	}

	var transactions []model.Transaction
	query := `SELECT id, wallet_id, type, amount, related_wallet_id, created_at
              FROM transactions
              WHERE wallet_id = $1
              ORDER BY created_at DESC` // Order by most recent

	err = wr.db.SelectContext(ctx, &transactions, query, walletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []model.Transaction{}, nil // Return empty slice if no transactions
		}
		return nil, fmt.Errorf("database error retrieving transactions: %w", err)
	}
	return transactions, nil
}

// GetWalletByUserIDAndWalletID retrieves a specific wallet for a user.
// userIDStr and walletIDStr are strings that will be parsed to uuid.UUID.
func (wr *WalletRepoImpl) GetWalletInfo(ctx context.Context, userIDStr string, walletIDStr string) (*model.Wallet, error) {
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

// Withdraw removes funds from a wallet and creates a transaction record.
func (wr *WalletRepoImpl) Withdraw(ctx context.Context, userIDStr string, walletIDStr string, amount decimal.Decimal) (*model.Transaction, error) {
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
	defer tx.Rollback()

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
		return nil, fmt.Errorf("failed to retrieve wallet for withdrawal: %w", err)
	}

	// 2. Check for sufficient funds
	if wallet.Balance.LessThan(amount) {
		return nil, ErrInsufficientFunds
	}

	// 3. Update wallet balance
	newBalance := wallet.Balance.Sub(amount)
	updateQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateQuery, newBalance, time.Now(), walletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update wallet balance for withdrawal: %w", err)
	}

	// 4. Create transaction record
	transaction := &model.Transaction{
		ID:        uuid.New(),
		WalletID:  walletID,
		Type:      model.TransactionTypeWithdrawal,
		Amount:    amount,
		CreatedAt: time.Now(),
	}
	insertTxQuery := `INSERT INTO transactions (id, wallet_id, type, amount, created_at)
                      VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.ExecContext(ctx, insertTxQuery, transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create withdrawal transaction record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit withdrawal transaction: %w", err)
	}

	return transaction, nil
}

// Transfer moves funds from a source wallet to a destination wallet and creates a transaction record.
func (wr *WalletRepoImpl) Transfer(ctx context.Context, sourceUserIDStr string, sourceWalletIDStr string, destinationWalletIDStr string, amount decimal.Decimal) (*model.Transaction, error) {
	sourceUserID, err := uuid.Parse(sourceUserIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid source user ID format: %w", err)
	}
	sourceWalletID, err := uuid.Parse(sourceWalletIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid source wallet ID format: %w", err)
	}
	destinationWalletID, err := uuid.Parse(destinationWalletIDStr)
	if err != nil {
		return nil, fmt.Errorf("invalid destination wallet ID format: %w", err)
	}

	if sourceWalletID == destinationWalletID {
		return nil, errors.New("source and destination wallets cannot be the same")
	}

	tx, err := wr.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// 1. Retrieve and lock the source wallet
	var sourceWallet model.Wallet
	// Ensure wallets are locked in a consistent order (e.g., by ID) to prevent deadlocks if concurrent transfers happen between the same two wallets in reverse.
	// For simplicity here, we assume different users or infrequent enough operations that deadlock isn't an immediate major concern for this example.
	// A robust solution would involve sorting wallet IDs before locking.
	querySourceWallet := `SELECT id, user_id, name, balance, created_at, updated_at
                          FROM wallets
                          WHERE user_id = $1 AND id = $2 FOR UPDATE`
	err = tx.GetContext(ctx, &sourceWallet, querySourceWallet, sourceUserID, sourceWalletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("source wallet not found: %w", ErrWalletNotFound)
		}
		return nil, fmt.Errorf("failed to retrieve source wallet for transfer: %w", err)
	}

	// 2. Check for sufficient funds in source wallet
	if sourceWallet.Balance.LessThan(amount) {
		return nil, ErrInsufficientFunds
	}

	// 3. Retrieve and lock the destination wallet
	var destinationWallet model.Wallet
	queryDestWallet := `SELECT id, user_id, name, balance, created_at, updated_at
                        FROM wallets
                        WHERE id = $1 FOR UPDATE` // Destination wallet can belong to any user
	err = tx.GetContext(ctx, &destinationWallet, queryDestWallet, destinationWalletID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("destination wallet not found: %w", ErrWalletNotFound)
		}
		return nil, fmt.Errorf("failed to retrieve destination wallet for transfer: %w", err)
	}

	// 4. Update source wallet balance
	newSourceBalance := sourceWallet.Balance.Sub(amount)
	updateSourceQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateSourceQuery, newSourceBalance, time.Now(), sourceWalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update source wallet balance for transfer: %w", err)
	}

	// 5. Update destination wallet balance
	newDestinationBalance := destinationWallet.Balance.Add(amount)
	updateDestQuery := `UPDATE wallets SET balance = $1, updated_at = $2 WHERE id = $3`
	_, err = tx.ExecContext(ctx, updateDestQuery, newDestinationBalance, time.Now(), destinationWalletID)
	if err != nil {
		return nil, fmt.Errorf("failed to update destination wallet balance for transfer: %w", err)
	}

	// 6. Create transaction record for the transfer (from the perspective of the source wallet)
	transaction := &model.Transaction{
		ID:              uuid.New(),
		WalletID:        sourceWalletID,
		Type:            model.TransactionTypeTransfer,
		Amount:          amount, // Amount is positive, representing outgoing from source
		RelatedWalletID: &destinationWalletID,
		CreatedAt:       time.Now(),
	}
	insertTxQuery := `INSERT INTO transactions (id, wallet_id, type, amount, related_wallet_id, created_at)
                      VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, insertTxQuery, transaction.ID, transaction.WalletID, transaction.Type, transaction.Amount, transaction.RelatedWalletID, transaction.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer transaction record: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transfer transaction: %w", err)
	}

	return transaction, nil
}
