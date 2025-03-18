package scenarios

import (
	"context"
	"fmt"

	"github.com/Axel791/order/internal/domains"
	"github.com/Axel791/order/internal/grpc/v1/pb"
	"github.com/Axel791/order/internal/usecases/order/dto"
	"github.com/Axel791/order/internal/usecases/order/repositories"
)

// CreateOrderUseCaseHandler - описание структуры use case
type CreateOrderUseCaseHandler struct {
	orderRepository   repositories.OrderRepository
	grpcLoyaltyClient pb.ConclusionUserBalanceUseCaseClient
}

// NewCreateOrderUseCase - конструктор use case создания order
func NewCreateOrderUseCase(
	orderRepository repositories.OrderRepository,
	grpcLoyaltyClient pb.ConclusionUserBalanceUseCaseClient,
) *CreateOrderUseCaseHandler {
	return &CreateOrderUseCaseHandler{
		orderRepository:   orderRepository,
		grpcLoyaltyClient: grpcLoyaltyClient,
	}
}

// Execute - метод, который выполняет всю бизнес логику на создание заказа
func (s *CreateOrderUseCaseHandler) Execute(ctx context.Context, order dto.CreateOrder) (dto.Order, error) {
	orderDomain := domains.Order{
		UserID:     order.UserID,
		Code:       order.Code,
		TotalPrice: order.TotalPrice,
	}

	if err := orderDomain.ValidateUserID(); err != nil {
		return dto.Order{}, err
	}

	if err := orderDomain.ValidateTotalPrice(); err != nil {
		return dto.Order{}, err
	}

	extractedOrder, err := s.orderRepository.CreateOrder(ctx, orderDomain)
	if err != nil {
		return dto.Order{}, fmt.Errorf("error create order: %w", err)
	}

	resp, err := s.grpcLoyaltyClient.Conclude(
		ctx,
		&pb.ConclusionRequest{
			UserId:  extractedOrder.UserID,
			OrderId: extractedOrder.ID,
			Count:   extractedOrder.TotalPrice,
		},
	)
	if err != nil {
		return dto.Order{}, fmt.Errorf("error order service RPC: %w", err)
	}
	if !resp.Success {
		return dto.Order{}, fmt.Errorf("lerror order service RPC: %s", resp.Message)
	}

	var orderDTO dto.Order

	orderDTO = dto.Order{
		ID:         extractedOrder.ID,
		UserID:     extractedOrder.UserID,
		Code:       extractedOrder.Code,
		TotalPrice: extractedOrder.TotalPrice,
	}

	return orderDTO, nil
}
