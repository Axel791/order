package scenarios

import (
	"context"

	"github.com/Axel791/order/internal/usecases/order/dto"
)

type CreateOrderUseCase interface {
	Execute(ctx context.Context, order dto.CreateOrder) (dto.Order, error)
}
