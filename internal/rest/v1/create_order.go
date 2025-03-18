package v1

import (
	"encoding/json"
	"net/http"

	"github.com/Axel791/appkit"
	"github.com/Axel791/order/internal/rest/v1/api"
	"github.com/Axel791/order/internal/usecases/order/dto"
	"github.com/Axel791/order/internal/usecases/order/scenarios"
	log "github.com/sirupsen/logrus"
)

type CreateOrderHandler struct {
	logger             *log.Logger
	createOrderUseCase scenarios.CreateOrderUseCase
}

func NewCreateOrderHandler(
	logger *log.Logger,
	createOrderUseCase scenarios.CreateOrderUseCase,
) *CreateOrderHandler {
	return &CreateOrderHandler{
		logger:             logger,
		createOrderUseCase: createOrderUseCase,
	}
}

func (h *CreateOrderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input api.InputCreateOrder
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.logger.Infof("err decode body: %v", err)
		appkit.WriteErrorJSON(w, appkit.BadRequestError("invalid request body"))
		return
	}

	createOrderDTO := dto.CreateOrder{
		UserID:     input.UserID,
		Code:       input.Code,
		TotalPrice: input.TotalPrice,
	}

	order, err := h.createOrderUseCase.Execute(r.Context(), createOrderDTO)
	if err != nil {
		h.logger.Infof("err execute: %v", err)
		appkit.WriteErrorJSON(w, err)
	}
	appkit.WriteJSON(
		w,
		http.StatusOK,
		api.OrderResponse{
			ID:         order.ID,
			UserID:     order.UserID,
			TotalPrice: order.TotalPrice,
			Code:       order.Code,
		},
	)
}
