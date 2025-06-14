package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
)

type OrderSummaryService interface {
	GetCustomerTotalSummary(
		ctx context.Context, customerID int, orderType []string, excludeStatus []string,
	) (*model.CustomerOrderSummary, error)

	GetLastOrderDate(
		ctx context.Context, customerID int, orderType []string,
	) (*time.Time, error)

	GetCustomerMonthlySummaries(
		ctx context.Context, customerID int, orderType []string, excludeStatus []string, months int,
	) ([]model.CustomerMonthlyOrderSummary, error)
}

type OrderSummary struct {
	service OrderSummaryService
}

func NewOrderSummary(service OrderSummaryService) *OrderSummary {
	return &OrderSummary{service: service}
}

func (h *OrderSummary) GetCustomerOrderSummary(c *gin.Context) {
	customerID, err := strconv.Atoi(c.Param("customerId"))
	if err != nil {
		restjson.ResponseError(c, http.StatusBadRequest, model.ErrInvalidCustomerID)
		return
	}

	var req struct {
		OrderType     []string `form:"orderType"`
		ExcludeStatus []string `form:"excludeStatus"`
		Months        int      `form:"months,default=3"`
	}

	if err = c.ShouldBindQuery(&req); err != nil {
		restjson.ResponseError(c, http.StatusBadRequest, err)
		return
	}

	if req.Months <= 0 {
		restjson.ResponseError(c, http.StatusBadRequest, model.ErrMonthsMustBeGreaterThanZero)
		return
	}

	ctx := c.Request.Context()

	summary, err := h.service.GetCustomerTotalSummary(ctx, customerID, req.OrderType, req.ExcludeStatus)
	if err != nil {
		restjson.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	lastOrderDate, err := h.service.GetLastOrderDate(ctx, customerID, req.OrderType)
	if err != nil {
		restjson.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	monthlySummaries, err := h.service.GetCustomerMonthlySummaries(
		ctx, customerID, req.OrderType, req.ExcludeStatus, req.Months)
	if err != nil {
		restjson.ResponseError(c, http.StatusInternalServerError, err)
		return
	}

	response := model.OrderSummaryResponse{
		Summary: model.CustomerOrderSummaryResponse{
			TotalAmount: summary.TotalAmount,
			TotalOrders: summary.OrderCount,
		},
	}

	if lastOrderDate != nil && !lastOrderDate.IsZero() {
		response.Summary.LastOrderedOn = lastOrderDate.Format(time.DateTime)
	}

	for _, summary := range monthlySummaries {
		response.Summary.MonthlySummary = append(response.Summary.MonthlySummary, model.CustomerOrderMonthlySummaryResponse{
			OrderAmount: summary.TotalAmount,
			OrderCount:  summary.OrderCount,
			Month:       summary.Month.Format("2006-01"),
		})
	}

	restjson.ResponseData(c, response)
}

func (h *OrderSummary) GetUserTransactions(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	restjson.ResponseData(c, userId)
}
