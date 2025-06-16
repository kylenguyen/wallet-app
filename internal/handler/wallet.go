package handler

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WalletService interface {
	GetWalletTransactions(ctx context.Context) ([]string, error)

	GetWalletInfo(ctx context.Context, userId, walletId string) (*model.Wallet, error)
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

func (h *WalletHandler) GetWalletInfo(c *gin.Context) {
	userId := c.Param("userId")
	walletId := c.Param("walletId")

	if len(userId) == 0 || len(walletId) == 0 {
		restjson.ResponseError(c, http.StatusBadRequest, errors.New("userId or walletId is invalid"))
		return
	}

	result, err := h.wService.GetWalletInfo(c.Request.Context(), userId, walletId)

	if err != nil {
		restjson.ResponseError(c, http.StatusInternalServerError, err)
		return
	}
	restjson.ResponseData(c, result)
}
