package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

func (s *service) CancelOrder(ctx context.Context, orderId string) (err error) {
	order := s.orderRepo.GetOrder(orderId)
	if order == nil {
		return model.ErrNotFound
	}

	switch order.Status {
	case model.OrderStatusPENDINGPAYMENT:
		order.Status = model.OrderStatusCANCELLED
		s.orderRepo.UpdateOrder(order)
		return nil
	case model.OrderStatusPAID:
		return model.ErrOrderAlreadyPaid
		// already cancelled
	default:
		return nil
	}
}
