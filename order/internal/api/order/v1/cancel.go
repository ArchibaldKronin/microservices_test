package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrderByUUID(ctx context.Context, params order_v1.CancelOrderByUUIDParams) (order_v1.CancelOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	err := a.orderService.CancelOrder(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order for id %s not found", id),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &order_v1.ConflictError{
				Code:    409,
				Message: fmt.Sprintf("Order %s is already paid", id),
			}, nil
		default:
			log.Printf("cancel order failed: %v\n", err)
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}

	return &order_v1.CancelOrderByUUIDNoContent{}, nil
}
