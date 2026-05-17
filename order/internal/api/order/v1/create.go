package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (a *api) CreateOrder(
	ctx context.Context,
	req *order_v1.CreateOrderRequest,
	params order_v1.CreateOrderParams,
) (order_v1.CreateOrderRes, error) {
	ctx = logger.WithUserID(ctx, req.UserUUID.String())

	partIds := converter.UUIDsToString(req.PartUuids)
	order, err := a.orderService.CreateOrder(ctx, req.UserUUID.String(), partIds)
	logger.Info(ctx, "Успешно создан заказ", zap.String("ID", order.OrderID))
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
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	orderId, err := uuid.Parse(order.OrderID)
	if err != nil {
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
