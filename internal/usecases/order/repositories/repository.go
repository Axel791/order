package repositories

import (
	"context"

	"github.com/Axel791/order/internal/domains"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, order domains.Order) (domains.Order, error)
}
