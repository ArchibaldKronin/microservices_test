package order

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) CompleteOrder(ctx context.Context, orderId string) error {
	order, err := s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		logger.Error(ctx, "error getting order", zap.String("id", orderId), zap.Error(err))

		if errors.Is(err, repoModel.ErrNotFound) {
			return model.ErrNotFound
		}
		return model.ErrInternal
	}

	switch order.Status {
	case model.OrderStatusPAID:
		// change
		order.Status = model.OrderStatusCOMPLETED
		err := s.orderRepo.UpdateOrder(ctx, order)
		if err != nil {
			if errors.Is(err, repoModel.ErrNotFound) {
				logger.Error(ctx, "error NON CONSISTENT DATA", zap.String("id", orderId), zap.Error(err))

				return model.ErrNotFound
			} else {
				logger.Error(ctx, "error updating order", zap.String("id", orderId), zap.Error(err))

				return model.ErrInternal
			}
		}
	default:
		logger.Error(ctx, "order status not PAID", zap.Error(model.ErrUnexpectedOrderStatus))
		return model.ErrUnexpectedOrderStatus
	}

	return nil
}
