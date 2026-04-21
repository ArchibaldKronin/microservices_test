package v1

import (
	"context"
	"errors"
	"log"

	"github.com/ArchibaldKronin/microservices_test/order/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (a *api) CreateOrder(ctx context.Context, req *order_v1.CreateOrderRequest) (order_v1.CreateOrderRes, error) {
	partIds := converter.UUIDsToString(req.PartUuids)
	order, err := a.orderService.CreateOrder(ctx, req.UserUUID.String(), partIds)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "some parts not found",
			}, nil
		case errors.Is(err, model.ErrUnavailable):
			return &order_v1.ServiceUnavailableError{
				Code:    503,
				Message: "create order error: unavailable",
			}, nil
		default:
			log.Printf("create order failed: %v\n", err)
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	orderId, err := uuid.Parse(order.OrderID)
	if err != nil {
		log.Printf("ошибка парсинга uuid: %v\n", err)
		return &order_v1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}

	resp := order_v1.CreateOrderResponse{
		OrderUUID:  orderId,
		TotalPrice: float32(order.TotalPrice),
	}

	return &resp, nil
}
