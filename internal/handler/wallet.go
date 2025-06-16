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
	Withdraw(ctx context.Context, userId, walletId string, amount decimal.Decimal) (*model.Transaction, error)
	Transfer(ctx context.Context, sourceUserId, sourceWalletId, destinationWalletId string, amount decimal.Decimal) (*model.Transaction, error)
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

// Withdraw handles withdrawing funds from a wallet.
// POST /api/v1/users/{user_id}/wallets/{wallet_id}/withdraw
func (h *WalletHandler) Withdraw(c *gin.Context) {
	userId := c.Param("userId")
	walletId := c.Param("walletId")

	if userId == "" || walletId == "" {
		restjson.ResponseError(c, http.StatusBadRequest, errors.New("userId or walletId is invalid in path"))
		return
	}

	var req model.AmountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restjson.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	transaction, err := h.wService.Withdraw(c.Request.Context(), userId, walletId, req.Amount)
	if err != nil {
		if errors.Is(err, repo.ErrWalletNotFound) {
			restjson.ResponseError(c, http.StatusNotFound, err)
		} else if errors.Is(err, repo.ErrInsufficientFunds) {
			restjson.ResponseError(c, http.StatusBadRequest, err) // Or http.StatusUnprocessableEntity
		} else {
			// log.Printf("Error in Withdraw: %v", err)
			restjson.ResponseError(c, http.StatusInternalServerError, errors.New("failed to process withdrawal"))
		}
		return
	}
	restjson.ResponseData(c, transaction)
}

// Transfer handles transferring funds between wallets.
// POST /api/v1/users/{user_id}/wallets/{from_wallet_id}/transfer
func (h *WalletHandler) Transfer(c *gin.Context) {
	sourceUserId := c.Param("userId") // Renamed from user_id in path to sourceUserId for clarity
	sourceWalletId := c.Param("walletId")

	if sourceUserId == "" || sourceWalletId == "" {
		restjson.ResponseError(c, http.StatusBadRequest, errors.New("source userId or source walletId is invalid in path"))
		return
	}

	var req model.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restjson.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if sourceWalletId == req.DestinationWalletID {
		restjson.ResponseError(c, http.StatusBadRequest, errors.New("source and destination wallets cannot be the same"))
		return
	}

	transaction, err := h.wService.Transfer(c.Request.Context(), sourceUserId, sourceWalletId, req.DestinationWalletID, req.Amount)
	if err != nil {
		if errors.Is(err, repo.ErrWalletNotFound) { // This could be source or destination
			restjson.ResponseError(c, http.StatusNotFound, err) // Consider more specific error messages if needed
		} else if errors.Is(err, repo.ErrInsufficientFunds) {
			restjson.ResponseError(c, http.StatusBadRequest, err)
		} else {
			// log.Printf("Error in Transfer: %v", err)
			restjson.ResponseError(c, http.StatusInternalServerError, errors.New("failed to process transfer"))
		}
		return
	}
	restjson.ResponseData(c, transaction)
}
