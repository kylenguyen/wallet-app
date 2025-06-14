package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"bitbucket.org/ntuclink/ff-order-history-go/internal/model"
)

type OrderSummary struct {
	db *sqlx.DB
}

func NewOrderSummary(db *sqlx.DB) *OrderSummary {
	return &OrderSummary{db: db}
}

func (r *OrderSummary) Calculate(
	ctx context.Context,
	customerID int,
	criteria model.CustomerOrderSummaryGetCriteria,
) (*model.CustomerOrderSummary, error) {
	var params []interface{}

	query := `
	SELECT
		IFNULL(SUM(amount), 0) AS amount,
		COUNT(amount) AS count
	FROM
		orders
	INNER JOIN
		order_types ON orders.type_id = order_types.id
	WHERE
		orders.organization_id = 2
		AND orders.customer_id = ?
`

	params = append(params, customerID)

	if criteria.StartTime != nil && criteria.EndTime != nil {
		query += " AND orders.completed_at BETWEEN ? AND ?"

		params = append(params, *criteria.StartTime, *criteria.EndTime)
	}

	if len(criteria.ExcludeStatus) > 0 {
		inQuery, args, err := sqlx.In(" AND orders.status NOT IN (?)", criteria.ExcludeStatus)
		if err != nil {
			return nil, fmt.Errorf("error building exclude status query: %w", err)
		}

		query += inQuery

		params = append(params, args...)
	}

	if len(criteria.OrderType) > 0 {
		inQuery, args, err := sqlx.In(" AND order_types.name IN (?)", criteria.OrderType)
		if err != nil {
			return nil, fmt.Errorf("error building order type query: %w", err)
		}

		query += inQuery

		params = append(params, args...)
	}

	var orderSummary model.CustomerOrderSummary

	query = r.db.Rebind(query)

	err := r.db.GetContext(ctx, &orderSummary, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &orderSummary, nil
		}

		return nil, fmt.Errorf("error getting order summary: %w", err)
	}

	return &orderSummary, nil
}

func (r *OrderSummary) GetLastOrderDate(
	ctx context.Context,
	customerID int,
	criteria model.CustomerOrderSummaryGetCriteria,
) (*time.Time, error) {
	var params []interface{}

	query := `
	SELECT created_at 
	FROM orders 
	`

	if len(criteria.OrderType) > 0 {
		query = fmt.Sprintf("%s INNER JOIN order_types on orders.type_id = order_types.id", query)
	}

	query = fmt.Sprintf("%s WHERE orders.organization_id = 2 AND orders.customer_id = ?", query)

	params = append(params, customerID)

	if len(criteria.OrderType) > 0 {
		inQuery, args, err := sqlx.In(" AND order_types.name IN (?)", criteria.OrderType)
		if err != nil {
			return nil, fmt.Errorf("error building order type query: %w", err)
		}

		query += inQuery

		params = append(params, args...)
	}

	if len(criteria.ExcludeStatus) > 0 {
		inQuery, args, err := sqlx.In(" AND orders.status NOT IN (?)", criteria.ExcludeStatus)
		if err != nil {
			return nil, fmt.Errorf("error building exclude status query: %w", err)
		}

		query += inQuery

		params = append(params, args...)
	}

	query = fmt.Sprintf("%s ORDER BY orders.created_at DESC LIMIT 1", query)

	query = r.db.Rebind(query)

	var lastOrder time.Time

	err := r.db.GetContext(ctx, &lastOrder, query, params...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &lastOrder, nil
		}

		return nil, fmt.Errorf("error getting last order: %w", err)
	}

	return &lastOrder, nil
}
