package service

import (
	"context"
	"fmt"
	"github.com/kylenguyen/wallet-app/internal/model"
	"github.com/shopspring/decimal"
)

type WalletRepo interface {
	GetWalletInfo(ctx context.Context, userId string, walletId string) (*model.Wallet, error)
	GetTransactionsByWalletID(ctx context.Context, userIDStr string, walletIDStr string) ([]model.Transaction, error)
	Deposit(ctx context.Context, userIDStr string, walletIDStr string, amount decimal.Decimal) (*model.Transaction, error)
	Withdraw(ctx context.Context, userIDStr string, walletIDStr string, amount decimal.Decimal) (*model.Transaction, error)
	Transfer(ctx context.Context, sourceUserIDStr string, sourceWalletIDStr string, destinationWalletIDStr string, amount decimal.Decimal) (*model.Transaction, error)
}

type WalletServiceImpl struct {
	wRepo WalletRepo
}

func NewWalletImpl(wr WalletRepo) *WalletServiceImpl {
	return &WalletServiceImpl{wr}
}

func (ws *WalletServiceImpl) GetWalletInfo(ctx context.Context, userId, walletId string) (*model.Wallet, error) {
	return ws.wRepo.GetWalletInfo(ctx, userId, walletId)
}

func (ws *WalletServiceImpl) GetWalletTransactionsByWalletID(ctx context.Context, userId, walletId string) ([]model.Transaction, error) {
	transactions, err := ws.wRepo.GetTransactionsByWalletID(ctx, userId, walletId)
	if err != nil {
		return nil, fmt.Errorf("service.GetWalletTransactionsByWalletID: %w", err)
	}
	return transactions, nil
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

func (ws *WalletServiceImpl) Withdraw(ctx context.Context, userId, walletId string, amount decimal.Decimal) (*model.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("withdrawal amount must be positive")
	}
	transaction, err := ws.wRepo.Withdraw(ctx, userId, walletId, amount)
	if err != nil {
		return nil, fmt.Errorf("service.Withdraw: %w", err)
	}
	return transaction, nil
}

func (ws *WalletServiceImpl) Transfer(ctx context.Context, sourceUserId, sourceWalletId, destinationWalletId string, amount decimal.Decimal) (*model.Transaction, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("transfer amount must be positive")
	}
	if sourceWalletId == destinationWalletId {
		return nil, fmt.Errorf("source and destination wallets cannot be the same")
	}
	transaction, err := ws.wRepo.Transfer(ctx, sourceUserId, sourceWalletId, destinationWalletId, amount)
	if err != nil {
		return nil, fmt.Errorf("service.Transfer: %w", err)
	}
	return transaction, nil
}
