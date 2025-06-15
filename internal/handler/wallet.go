package handler

import (
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
	"context"
	"github.com/gin-gonic/gin"
)

type WalletService interface {
	GetWalletTransactions(ctx context.Context) ([]string, error)
}

func NewWalletImpl(wService WalletService) *WalletHandler {
	return &WalletHandler{wService}
}

type WalletHandler struct {
	wService WalletService
}

func (h *WalletHandler) GetWalletTransactions(c *gin.Context) {

	result, _ := h.wService.GetWalletTransactions(nil)
	restjson.ResponseData(c, result)
}
