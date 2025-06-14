package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/service"
	ordersummarymocks "bitbucket.org/ntuclink/ff-order-history-go/internal/service/mocks"
	datetimemocks "bitbucket.org/ntuclink/ff-order-history-go/pkg/datetime/mocks"
)

func TestOrderSummary_GetCustomerTotalSummary(t *testing.T) {
	testCases := []struct {
		name             string
		customerID       int
		orderType        []string
		excludeStatus    []string
		mockCalculate    *model.CustomerOrderSummary
		mockCalculateErr error
		expectedResult   *model.CustomerOrderSummary
		expectedErr      assert.ErrorAssertionFunc
	}{
		{
			name:           "success",
			customerID:     1,
			orderType:      []string{"DELIVERY"},
			excludeStatus:  []string{"CANCELLED"},
			mockCalculate:  &model.CustomerOrderSummary{TotalAmount: decimal.NewFromInt(100), OrderCount: 5},
			expectedResult: &model.CustomerOrderSummary{TotalAmount: decimal.NewFromInt(100), OrderCount: 5},
			expectedErr:    assert.NoError,
		},
		{
			name:             "repo error",
			customerID:       1,
			orderType:        []string{"DELIVERY"},
			excludeStatus:    []string{"CANCELLED"},
			mockCalculateErr: errors.New("repo error"),
			expectedResult:   nil,
			expectedErr:      assert.Error,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(ordersummarymocks.CustomerOrderSummaryGetterMock)
			mockDatetime := new(datetimemocks.GetterMock)

			mockRepo.On("Calculate", mock.Anything, tc.customerID, mock.Anything).Return(tc.mockCalculate, tc.mockCalculateErr)

			os := service.NewOrderSummary(mockRepo, mockDatetime)

			result, err := os.GetCustomerTotalSummary(context.Background(), tc.customerID, tc.orderType, tc.excludeStatus)

			tc.expectedErr(t, err)
			assert.Equal(t, tc.expectedResult, result)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrderSummary_GetLastOrderDate(t *testing.T) {
	testCases := []struct {
		name                string
		customerID          int
		orderType           []string
		mockGetLastOrderOn  *time.Time
		mockGetLastOrderErr error
		expectedResult      assert.ValueAssertionFunc
		expectedErr         assert.ErrorAssertionFunc
	}{
		{
			name:               "success",
			customerID:         1,
			orderType:          []string{"DELIVERY"},
			mockGetLastOrderOn: func() *time.Time { t := time.Date(2025, 3, 20, 0, 0, 0, 0, time.Local); return &t }(),
			expectedResult: func(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
				return assert.NotNil(tt, i1) && assert.Equal(tt, time.Date(2025, 3, 20, 0, 0, 0, 0, time.Local).Nanosecond(), i1.(*time.Time).Nanosecond())
			},
			expectedErr: assert.NoError,
		},
		{
			name:                "repo error",
			customerID:          1,
			orderType:           []string{"DELIVERY"},
			mockGetLastOrderErr: errors.New("repo error"),
			expectedResult:      assert.Nil,
			expectedErr:         assert.Error,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(ordersummarymocks.CustomerOrderSummaryGetterMock)
			mockDatetime := new(datetimemocks.GetterMock)

			mockRepo.On("GetLastOrderDate", mock.Anything, tc.customerID, mock.Anything).Return(tc.mockGetLastOrderOn, tc.mockGetLastOrderErr)

			os := service.NewOrderSummary(mockRepo, mockDatetime)

			result, err := os.GetLastOrderDate(context.Background(), tc.customerID, tc.orderType)

			tc.expectedErr(t, err)
			tc.expectedResult(t, result)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrderSummary_GetCustomerMonthlySummaries(t *testing.T) {
	mockStart := time.Date(2023, 10, 1, 0, 0, 0, 0, time.Local)
	mockEnd := time.Date(2023, 10, 31, 23, 59, 59, 0, time.Local)

	testCases := []struct {
		name             string
		customerID       int
		orderType        []string
		excludeStatus    []string
		months           int
		mockStart        time.Time
		mockEnd          time.Time
		mockCalculate    *model.CustomerOrderSummary
		mockCalculateErr error
		expectedErr      assert.ErrorAssertionFunc
	}{
		{
			name:          "success",
			customerID:    1,
			orderType:     []string{"DELIVERY"},
			excludeStatus: []string{"CANCELLED"},
			months:        3,
			mockStart:     mockStart,
			mockEnd:       mockEnd,
			mockCalculate: &model.CustomerOrderSummary{TotalAmount: decimal.NewFromInt(100), OrderCount: 5},
			expectedErr:   assert.NoError,
		},
		{
			name:             "repo error",
			customerID:       1,
			orderType:        []string{"DELIVERY"},
			excludeStatus:    []string{"CANCELLED"},
			months:           3,
			mockStart:        mockStart,
			mockEnd:          mockEnd,
			mockCalculateErr: errors.New("repo error"),
			expectedErr:      assert.Error,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(ordersummarymocks.CustomerOrderSummaryGetterMock)
			mockDatetime := new(datetimemocks.GetterMock)

			mockDatetime.On("GetMonthStartAndEnd", mock.Anything).Return(tc.mockStart, tc.mockEnd)
			mockRepo.On("Calculate", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockCalculate, tc.mockCalculateErr)

			os := service.NewOrderSummary(mockRepo, mockDatetime)

			_, err := os.GetCustomerMonthlySummaries(context.Background(), tc.customerID, tc.orderType, tc.excludeStatus, tc.months)

			tc.expectedErr(t, err)

			mockRepo.AssertExpectations(t)
			mockDatetime.AssertExpectations(t)
		})
	}
}
