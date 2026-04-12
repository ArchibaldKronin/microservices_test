package v1

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ArchibaldKronin/microservices_test/order/internal/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderByUUID(ctx context.Context, params order_v1.GetOrderByUUIDParams) (order_v1.GetOrderByUUIDRes, error) {
	id := params.OrderUUID.String()
	order, err := a.orderService.GetOrder(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrNotFound):
			return &order_v1.NotFoundError{
				Code:    404,
				Message: fmt.Sprintf("Order for id %s not found", id),
			}, nil
		case errors.Is(err, model.ErrUnavailable):
			return &order_v1.ServiceUnavailableError{
				Code:    503,
				Message: "get order error: unavailable",
			}, nil
		default:
			log.Printf("get order failed: %v\n", err)
			return &order_v1.InternalServerError{
				Code:    500,
				Message: "internal server error",
			}, nil
		}
	}
	orderDTO, err := converter.OrderToDTO(order)
	if err != nil {
		log.Printf("order convertion failed: %v\n", err)

		return &order_v1.InternalServerError{
			Code:    500,
			Message: "internal server error",
		}, nil
	}
	return &orderDTO, nil
}
