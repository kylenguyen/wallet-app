package repo

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

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
