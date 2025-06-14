package rest_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	httphandler "bitbucket.org/ntuclink/ff-order-history-go/internal/handler/rest"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/handler/rest/mocks"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/pkg/restjson"
)

func TestOrderSummary_GetCustomerOrderSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockTime := time.Date(2025, 3, 20, 0, 0, 0, 0, time.Local)

	testCases := []struct {
		name                     string
		customerID               string
		queryParams              url.Values
		mockGetTotalSummary      *model.CustomerOrderSummary
		mockGetTotalSummaryErr   error
		mockGetLastOrderDate     *time.Time
		mockGetLastOrderDateErr  error
		mockGetMonthlySummary    []model.CustomerMonthlyOrderSummary
		mockGetMonthlySummaryErr error
		expectedCode             int
		expectedResponse         model.OrderSummaryResponse
		expectedError            string
	}{
		{
			name:       "Success",
			customerID: "1",
			queryParams: url.Values{
				"orderType":     []string{"DELIVERY"},
				"excludeStatus": []string{"CANCELLED"},
				"months":        []string{"3"},
			},
			mockGetTotalSummary: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromInt(100),
				OrderCount:  5,
			},
			mockGetLastOrderDate: func() *time.Time { t := mockTime; return &t }(),
			mockGetMonthlySummary: []model.CustomerMonthlyOrderSummary{
				{
					CustomerOrderSummary: model.CustomerOrderSummary{
						TotalAmount: decimal.NewFromInt(40),
						OrderCount:  2,
					},
					Month: time.Date(2025, 3, 1, 0, 0, 0, 0, time.Local),
				},
				{
					CustomerOrderSummary: model.CustomerOrderSummary{
						TotalAmount: decimal.NewFromInt(60),
						OrderCount:  3,
					},
					Month: time.Date(2025, 2, 1, 0, 0, 0, 0, time.Local),
				},
			},
			expectedCode: http.StatusOK,
			expectedResponse: model.OrderSummaryResponse{
				Summary: model.CustomerOrderSummaryResponse{
					TotalAmount:   decimal.NewFromInt(100),
					TotalOrders:   5,
					LastOrderedOn: mockTime.Format(time.DateTime),
					MonthlySummary: []model.CustomerOrderMonthlySummaryResponse{
						{
							OrderAmount: decimal.NewFromInt(40),
							OrderCount:  2,
							Month:       "2025-03",
						},
						{
							OrderAmount: decimal.NewFromInt(60),
							OrderCount:  3,
							Month:       "2025-02",
						},
					},
				},
			},
		},
		{
			name:          "Invalid Customer ID",
			customerID:    "abc",
			expectedCode:  http.StatusBadRequest,
			expectedError: "invalid customerId",
		},
		{
			name:       "Invalid Months",
			customerID: "1",
			queryParams: url.Values{
				"months": []string{"0"},
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: model.ErrMonthsMustBeGreaterThanZero.Error(),
		},
		{
			name:       "GetCustomerTotalSummary error",
			customerID: "1",
			queryParams: url.Values{
				"months": []string{"3"},
			},
			mockGetTotalSummaryErr: errors.New("GetCustomerTotalSummary error"),
			expectedCode:           http.StatusInternalServerError,
			expectedError:          "GetCustomerTotalSummary error",
		},
		{
			name:       "GetLastOrderDate Error",
			customerID: "1",
			queryParams: url.Values{
				"months": []string{"3"},
			},
			mockGetTotalSummary: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromInt(100),
				OrderCount:  5,
			},
			mockGetLastOrderDateErr: errors.New("GetLastOrderDate error"),
			expectedCode:            http.StatusInternalServerError,
			expectedError:           "GetLastOrderDate error",
		},
		{
			name:       "GetCustomerMonthlySummaries Error",
			customerID: "1",
			queryParams: url.Values{
				"months": []string{"3"},
			},
			mockGetTotalSummary: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromInt(100),
				OrderCount:  5,
			},
			mockGetLastOrderDate:     func() *time.Time { t := mockTime; return &t }(),
			mockGetMonthlySummaryErr: errors.New("GetCustomerMonthlySummaries error"),
			expectedCode:             http.StatusInternalServerError,
			expectedError:            "GetCustomerMonthlySummaries error",
		},
		{
			name:       "Bind Query Error",
			customerID: "1",
			queryParams: url.Values{
				"months": []string{"abc"},
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "strconv.ParseInt: parsing \"abc\": invalid syntax",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Mock setup
			mockService := new(mocks.OrderSummaryServiceMock)

			if tc.mockGetTotalSummary != nil || tc.mockGetTotalSummaryErr != nil {
				mockService.On("GetCustomerTotalSummary", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.mockGetTotalSummary, tc.mockGetTotalSummaryErr)
			}

			if tc.mockGetLastOrderDate != nil || tc.mockGetLastOrderDateErr != nil {
				mockService.On("GetLastOrderDate", mock.Anything, mock.Anything, mock.Anything).Return(tc.mockGetLastOrderDate, tc.mockGetLastOrderDateErr)
			}

			if tc.mockGetMonthlySummary != nil || tc.mockGetMonthlySummaryErr != nil {
				mockService.On("GetCustomerMonthlySummaries", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(tc.mockGetMonthlySummary, tc.mockGetMonthlySummaryErr)
			}

			// Handler setup
			h := httphandler.NewOrderSummary(mockService)

			w := httptest.NewRecorder()

			c, _ := gin.CreateTestContext(w)

			req := &http.Request{
				URL: &url.URL{
					Path: fmt.Sprintf("/order-history/v1/customers/%s/order-summary", tc.customerID),
				},
				Method: http.MethodGet,
			}

			req.URL.RawQuery = tc.queryParams.Encode()

			c.Request = req
			c.Params = []gin.Param{
				{
					Key:   "customerId",
					Value: tc.customerID,
				},
			}

			// Call handler
			h.GetCustomerOrderSummary(c)

			// Assertions
			assert.Equal(t, tc.expectedCode, w.Code, "Unexpected HTTP status code")

			if tc.expectedError != "" {
				assertErrorResponse(t, w, tc.expectedError)
			} else {
				assertSuccessResponse(t, w, tc.expectedResponse)
			}

			// Mock assertions
			mockService.AssertExpectations(t)
		})
	}
}

func assertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedError string) {
	t.Helper()

	require.NotEmpty(t, w.Body.Bytes(), "Response body should not be empty for error responses")

	var resp restjson.ErrorResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "Failed to unmarshal error response")
	assert.Equal(t, expectedError, resp.Message, "Unexpected error message")
}

func assertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder, expectedResponse model.OrderSummaryResponse) {
	t.Helper()

	require.NotEmpty(t, w.Body.Bytes(), "Response body should not be empty for success responses")

	var resp restjson.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "Failed to unmarshal success response")

	var actualResponse model.OrderSummaryResponse

	dataBytes, err := json.Marshal(resp.Data)
	require.NoError(t, err, "Failed to marshal response data")
	err = json.Unmarshal(dataBytes, &actualResponse)
	require.NoError(t, err, "Failed to unmarshal response data")
	assert.Equal(t, expectedResponse, actualResponse, "Unexpected response data")
}
