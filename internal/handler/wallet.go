package handler

import (
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/repo"
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"net/http"
)

type WalletService interface {
	GetWalletTransactions(ctx context.Context) ([]string, error)

	GetWalletInfo(ctx context.Context, userId, walletId string) (*model.Wallet, error)
	Deposit(ctx context.Context, userId, walletId string, amount decimal.Decimal) (*model.Transaction, error)
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

// Deposit handles depositing funds into a wallet.
// POST /api/v1/users/{user_id}/wallets/{wallet_id}/deposit
func (h *WalletHandler) Deposit(c *gin.Context) {
	userId := c.Param("userId")
	walletId := c.Param("walletId")

	if userId == "" || walletId == "" {
		restjson.ResponseError(c, http.StatusBadRequest, errors.New("userId or walletId is invalid in path"))
		return
	}

	var req model.AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restjson.ResponseError(c, http.StatusBadRequest, err) // Gin's binding error is usually descriptive
		return
	}

	// Amount validation (e.g. positive) is handled by model binding `gt=0`
	// and can also be double-checked in the service layer.

	transaction, err := h.wService.Deposit(c.Request.Context(), userId, walletId, req.Amount)
	if err != nil {
		if errors.Is(err, repo.ErrWalletNotFound) {
			restjson.ResponseError(c, http.StatusNotFound, err)
		} else {
			// log.Printf("Error in Deposit: %v", err)
			restjson.ResponseError(c, http.StatusInternalServerError, errors.New("failed to process deposit"))
		}
		return
	}
	restjson.ResponseData(c, transaction)
}
