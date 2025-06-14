package model

import "errors"

var (
	ErrInvalidCustomerID           = errors.New("invalid customerId")
	ErrMonthsMustBeGreaterThanZero = errors.New("months must be greater than 0")
)
