package repo_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
	"bitbucket.org/ntuclink/ff-order-history-go/internal/repo"
)

func TestOrderSummary_Calculate(t *testing.T) {
	testCases := []struct {
		name           string
		customerID     int
		criteria       model.CustomerOrderSummaryGetCriteria
		query          string
		mockRows       *sqlmock.Rows
		mockErr        error
		expectedResult *model.CustomerOrderSummary
		expectedErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "Success Without Filters",
			customerID: 1,
			criteria:   model.CustomerOrderSummaryGetCriteria{},
			query:      `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \?`,
			mockRows: sqlmock.NewRows([]string{"amount", "count"}).
				AddRow(100.00, 5),
			mockErr: nil,
			expectedResult: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromFloat(100),
				OrderCount:  5,
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "Success With Time Range",
			customerID: 1,
			criteria: model.CustomerOrderSummaryGetCriteria{
				StartTime: func() *time.Time {
					t := time.Now().Add(-time.Hour)
					return &t
				}(),
				EndTime: func() *time.Time {
					t := time.Now()
					return &t
				}(),
			},
			query: `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \? AND orders.completed_at BETWEEN \? AND \?`,
			mockRows: sqlmock.NewRows([]string{"amount", "count"}).
				AddRow(90.00, 4),
			mockErr: nil,
			expectedResult: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromFloat(90),
				OrderCount:  4,
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "Success With Exclude Status",
			customerID: 1,
			criteria: model.CustomerOrderSummaryGetCriteria{
				ExcludeStatus: []string{"CANCELLED"},
			},
			query: `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \? AND orders.status NOT IN \(\?\)`,
			mockRows: sqlmock.NewRows([]string{"amount", "count"}).
				AddRow(80.00, 3),
			mockErr: nil,
			expectedResult: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromFloat(80),
				OrderCount:  3,
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "Success With Order Type",
			customerID: 1,
			criteria: model.CustomerOrderSummaryGetCriteria{
				OrderType: []string{"DELIVERY"},
			},
			query: `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \? AND order_types.name IN \(\?\)`,
			mockRows: sqlmock.NewRows([]string{"amount", "count"}).
				AddRow(80.00, 3),
			mockErr: nil,
			expectedResult: &model.CustomerOrderSummary{
				TotalAmount: decimal.NewFromFloat(80),
				OrderCount:  3,
			},
			expectedErr: assert.NoError,
		},
		{
			name:           "No Rows Found",
			customerID:     1,
			criteria:       model.CustomerOrderSummaryGetCriteria{},
			query:          `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \?`,
			mockRows:       sqlmock.NewRows([]string{"amount", "count"}),
			mockErr:        nil,
			expectedResult: &model.CustomerOrderSummary{},
			expectedErr:    assert.NoError,
		},
		{
			name:           "Db error",
			customerID:     1,
			criteria:       model.CustomerOrderSummaryGetCriteria{},
			query:          `SELECT IFNULL\(SUM\(amount\), 0\) AS amount, COUNT\(amount\) AS count FROM orders INNER JOIN order_types ON orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \?`,
			mockRows:       nil,
			mockErr:        sqlmock.ErrCancelled,
			expectedResult: nil,
			expectedErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(tt, err, i...) && assert.EqualError(tt, err, fmt.Sprintf("error getting order summary: %v", sqlmock.ErrCancelled))
			},
		},
	}

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	orderSummaryRepo := repo.NewOrderSummary(sqlxDB)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			switch {
			case tc.mockErr != nil:
				mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
			case tc.mockRows != nil:
				mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
			default:
				mock.ExpectQuery(regexp.QuoteMeta(tc.query))
			}

			result, err := orderSummaryRepo.Calculate(ctx, tc.customerID, tc.criteria)

			tc.expectedErr(t, err)
			assert.Equal(t, tc.expectedResult, result)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestOrderSummary_GetLastOrderDate(t *testing.T) {
	testCases := []struct {
		name           string
		customerID     int
		criteria       model.CustomerOrderSummaryGetCriteria
		query          string
		mockRows       *sqlmock.Rows
		mockErr        error
		expectedResult assert.ValueAssertionFunc
		expectedErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "Success Without Filters",
			customerID: 1,
			criteria:   model.CustomerOrderSummaryGetCriteria{},
			query:      `SELECT created_at FROM orders WHERE orders.organization_id = 2 AND orders.customer_id = \? ORDER BY orders.created_at DESC LIMIT 1`,
			mockRows:   sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()),
			mockErr:    nil,
			expectedResult: func(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
				return assert.NotNil(tt, i1) && assert.IsType(tt, &time.Time{}, i1) && assert.False(tt, i1.(*time.Time).IsZero())
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "Success With Order Type",
			customerID: 1,
			criteria:   model.CustomerOrderSummaryGetCriteria{OrderType: []string{"DELIVERY"}},
			query:      `SELECT created_at FROM orders INNER JOIN order_types on orders.type_id = order_types.id WHERE orders.organization_id = 2 AND orders.customer_id = \? AND order_types.name IN \(\?\) ORDER BY orders.created_at DESC LIMIT 1`,
			mockRows:   sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()),
			mockErr:    nil,
			expectedResult: func(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
				return assert.NotNil(tt, i1) && assert.IsType(tt, &time.Time{}, i1) && assert.False(tt, i1.(*time.Time).IsZero())
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "Success With Exclude Status",
			customerID: 1,
			criteria:   model.CustomerOrderSummaryGetCriteria{ExcludeStatus: []string{"CANCELLED"}},
			query:      `SELECT created_at FROM orders WHERE orders.organization_id = 2 AND orders.customer_id = \? AND orders.status NOT IN \(\?\) ORDER BY orders.created_at DESC LIMIT 1`,
			mockRows:   sqlmock.NewRows([]string{"created_at"}).AddRow(time.Now()),
			mockErr:    nil,
			expectedResult: func(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
				return assert.NotNil(tt, i1) && assert.IsType(tt, &time.Time{}, i1) && assert.False(tt, i1.(*time.Time).IsZero())
			},
			expectedErr: assert.NoError,
		},
		{
			name:       "No Rows Found",
			customerID: 1,
			criteria:   model.CustomerOrderSummaryGetCriteria{},
			query:      `SELECT created_at FROM orders WHERE orders.organization_id = 2 AND orders.customer_id = \? ORDER BY orders.created_at DESC LIMIT 1`,
			mockRows:   sqlmock.NewRows([]string{"created_at"}),
			mockErr:    nil,
			expectedResult: func(tt assert.TestingT, i1 interface{}, _ ...interface{}) bool {
				return assert.NotNil(tt, i1) && assert.IsType(tt, &time.Time{}, i1) && assert.True(tt, i1.(*time.Time).IsZero())
			},
			expectedErr: assert.NoError,
		},
		{
			name:           "Db Error",
			customerID:     1,
			criteria:       model.CustomerOrderSummaryGetCriteria{},
			query:          `SELECT created_at FROM orders WHERE orders.organization_id = 2 AND orders.customer_id = \? ORDER BY orders.created_at DESC LIMIT 1`,
			mockRows:       nil,
			mockErr:        sqlmock.ErrCancelled,
			expectedResult: assert.Empty,
			expectedErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.Error(tt, err, i...) && assert.EqualError(tt, err, fmt.Sprintf("error getting last order: %v", sqlmock.ErrCancelled))
			},
		},
	}

	ctx := context.Background()

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("error creating mock database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	orderSummaryRepo := repo.NewOrderSummary(sqlxDB)

	for _, tc := range testCases {
		switch {
		case tc.mockErr != nil:
			mock.ExpectQuery(tc.query).WillReturnError(tc.mockErr)
		case tc.mockRows != nil:
			mock.ExpectQuery(tc.query).WillReturnRows(tc.mockRows)
		default:
			mock.ExpectQuery(regexp.QuoteMeta(tc.query))
		}

		result, err := orderSummaryRepo.GetLastOrderDate(ctx, tc.customerID, tc.criteria)

		tc.expectedErr(t, err)
		tc.expectedResult(t, result)

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}
}
