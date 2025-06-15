package repo

import (
	"context"
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
