package service

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"context"
	//"net/http"
	//"strconv"
	//"time"
	//
	//"github.com/gin-gonic/gin"
	//
	//"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	//"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
)

type WalletRepo interface {
	GetWalletTransactions(ctx context.Context) ([]string, error)

	RetrieveWalletByUserIdAndWalletId(ctx context.Context, userId string, walletId string) (*model.Wallet, error)
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
