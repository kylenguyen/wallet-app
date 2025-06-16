package service

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"context"
	"fmt"
	"github.com/shopspring/decimal"
)

type WalletRepo interface {
	GetWalletTransactions(ctx context.Context) ([]string, error)

	RetrieveWalletByUserIdAndWalletId(ctx context.Context, userId string, walletId string) (*model.Wallet, error)
	Deposit(ctx context.Context, userIDStr string, walletIDStr string, amount decimal.Decimal) (*model.Transaction, error)
}

type WalletServiceImpl struct {
	wRepo WalletRepo
}

func NewWalletImpl(wr WalletRepo) *WalletServiceImpl {
	return &WalletServiceImpl{wr}
}

func (ws *WalletServiceImpl) GetWalletTransactions(ctx context.Context) ([]string, error) {
	return ws.wRepo.GetWalletTransactions(ctx)
}

func (ws *WalletServiceImpl) GetWalletInfo(ctx context.Context, userId, walletId string) (*model.Wallet, error) {
	return ws.wRepo.RetrieveWalletByUserIdAndWalletId(ctx, userId, walletId)
}

func (ws *WalletServiceImpl) Deposit(ctx context.Context, userId, walletId string, amount decimal.Decimal) (*model.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("deposit amount must be positive")
	}
	transaction, err := ws.wRepo.Deposit(ctx, userId, walletId, amount)
	if err != nil {
		return nil, fmt.Errorf("service.Deposit: %w", err)
	}
	return transaction, nil
}
