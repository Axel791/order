package repositories

import (
	"context"

	"github.com/Axel791/order/internal/domains"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

// SqlOrderRepository - структура репозитория
type SqlOrderRepository struct {
	db *sqlx.DB
}

// NewSqlOrderRepository - структура репозитория Order
func NewSqlOrderRepository(db *sqlx.DB) *SqlOrderRepository {
	return &SqlOrderRepository{
		db: db,
	}
}

// CreateOrder - создание заказа
func (r *SqlOrderRepository) CreateOrder(ctx context.Context, order domains.Order) (domains.Order, error) {
	query, args, err := sq.
		Insert("orders").
		Columns("user_id", "code", "total_price").
		Values(order.UserID, order.Code, order.TotalPrice).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return domains.Order{}, err
	}

	var id int64
	err = r.db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return domains.Order{}, err
	}

	order.ID = id
	return order, nil
}
