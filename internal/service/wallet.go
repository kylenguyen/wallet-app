package service

import (
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
