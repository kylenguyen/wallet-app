package model

import (
	"time"

	"github.com/shopspring/decimal"
)

// CustomerOrderSummaryGetCriteria represents the criteria used to fetch a customer's order summary.
type CustomerOrderSummaryGetCriteria struct {
	OrderType     []string
	StartTime     *time.Time
	EndTime       *time.Time
	ExcludeStatus []string
}

// CustomerOrderSummary is an entity object representing a summary of a customer's orders.
type CustomerOrderSummary struct {
	TotalAmount decimal.Decimal `db:"amount"`
	OrderCount  int             `db:"count"`
}

// CustomerMonthlyOrderSummary represents a customer's order summary for a specific month.
type CustomerMonthlyOrderSummary struct {
	CustomerOrderSummary
	Month time.Time
}
