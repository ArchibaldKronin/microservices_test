package v1

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func (a *api) PayOrderByUUID(ctx context.Context, req *order_v1.PayOrderRequest, params order_v1.PayOrderByUUIDParams) (order_v1.PayOrderByUUIDRes, error) {
	orderId := params.OrderUUID.String()
	paymentMethod := converter.PaymentMethodToDomain(req.PaymentMethod)

	transactionId, err := a.orderService.PayOrder(ctx, orderId, paymentMethod)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: "order not found",
			}, nil
		case errors.Is(err, model.ErrInternal):
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		case errors.Is(err, model.ErrUnavailable):
			return &order_v1.ServiceUnavailableError{
				Code:    503,
				Message: "pay order error: unavailable",
			}, nil
		default:
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	transactionUUID, err := uuid.Parse(transactionId)
	if err != nil {
		return &order_v1.InternalServerError{
			Code:    500,
			Message: "error parse UUID",
		}, nil
	}

	return &order_v1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
