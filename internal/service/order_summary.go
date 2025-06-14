package service

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/datetime"
)

type CustomerOrderSummaryGetter interface {
	Calculate(
		ctx context.Context, customerID int,
		criteria model.CustomerOrderSummaryGetCriteria) (*model.CustomerOrderSummary, error)
	GetLastOrderDate(
		ctx context.Context, customerID int,
		criteria model.CustomerOrderSummaryGetCriteria) (*time.Time, error)
}

type OrderSummary struct {
	repo     CustomerOrderSummaryGetter
	datetime datetime.Getter
}

func NewOrderSummary(
	repo CustomerOrderSummaryGetter, datetime datetime.Getter,
) *OrderSummary {
	return &OrderSummary{
		repo:     repo,
		datetime: datetime,
	}
}

func (s *OrderSummary) GetCustomerTotalSummary(
	ctx context.Context, customerID int, orderType []string, excludeStatus []string,
) (*model.CustomerOrderSummary, error) {
	criteria := model.CustomerOrderSummaryGetCriteria{
		OrderType:     orderType,
		ExcludeStatus: excludeStatus,
	}

	totalCos, err := s.repo.Calculate(ctx, customerID, criteria)
	if err != nil {
		return nil, fmt.Errorf("error getting customer order summary: %w", err)
	}

	return totalCos, nil
}

func (s *OrderSummary) GetLastOrderDate(
	ctx context.Context, customerID int, orderType []string,
) (*time.Time, error) {
	criteria := model.CustomerOrderSummaryGetCriteria{
		OrderType:     orderType,
		ExcludeStatus: []string{"CANCELLED"},
	}

	lastOrderDate, err := s.repo.GetLastOrderDate(ctx, customerID, criteria)
	if err != nil {
		return nil, fmt.Errorf("error getting last order on: %w", err)
	}

	return lastOrderDate, nil
}

func (s *OrderSummary) GetCustomerMonthlySummaries(
	ctx context.Context, customerID int, orderType []string, excludeStatus []string, months int,
) ([]model.CustomerMonthlyOrderSummary, error) {
	criteria := model.CustomerOrderSummaryGetCriteria{
		OrderType:     orderType,
		ExcludeStatus: excludeStatus,
	}

	var monthlySummaries []model.CustomerMonthlyOrderSummary

	for month := 1; month <= months; month++ {
		start, end := s.datetime.GetMonthStartAndEnd(month)

		criteria.StartTime = &start
		criteria.EndTime = &end

		cos, err := s.repo.Calculate(ctx, customerID, criteria)
		if err != nil {
			return nil, fmt.Errorf("error getting customer order summary for month %d: %w", month, err)
		}

		monthlySummaries = append(monthlySummaries, model.CustomerMonthlyOrderSummary{
			CustomerOrderSummary: *cos,
			Month:                start,
		})
	}

	return monthlySummaries, nil
}
