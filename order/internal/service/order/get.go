package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

func (s *service) GetOrder(ctx context.Context, orderId string) (order *model.Order, err error) {
	order = s.orderRepo.GetOrder(orderId)
	if order == nil {
		return nil, model.ErrNotFound
	}
	return order, nil
}
