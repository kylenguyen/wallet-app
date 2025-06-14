package model

import "github.com/shopspring/decimal"

type OrderSummaryResponse struct {
	Summary CustomerOrderSummaryResponse `json:"summary"`
}

type CustomerOrderSummaryResponse struct {
	TotalAmount    decimal.Decimal                       `json:"totalAmount"`
	TotalOrders    int                                   `json:"totalOrders"`
	LastOrderedOn  string                                `json:"lastOrderedOn,omitempty"`
	MonthlySummary []CustomerOrderMonthlySummaryResponse `json:"monthlySummary"`
}

type CustomerOrderMonthlySummaryResponse struct {
	OrderAmount decimal.Decimal `json:"orderAmount"`
	OrderCount  int             `json:"orderCount"`
	Month       string          `json:"month"`
}
