package order

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) CancelOrder(ctx context.Context, orderId string) error {
	err := s.txManager.WithTransaction(ctx, func(executer repository.OrderRepository) error {
		order, err := executer.GetOrder(ctx, orderId)
		if err != nil {
			logger.Error(ctx, "error getting order", zap.String("id", orderId), zap.Error(err))

			if errors.Is(err, repoModel.ErrNotFound) {
				return model.ErrNotFound
			}
			return model.ErrInternal
		}

		switch order.Status {
		case model.OrderStatusPENDINGPAYMENT:
			order.Status = model.OrderStatusCANCELLED
			err = executer.UpdateOrder(ctx, order)
			if err != nil {
				if errors.Is(err, repoModel.ErrNotFound) {
					logger.Error(ctx, "error NON CONSISTENT DATA", zap.String("id", orderId), zap.Error(err))

					return model.ErrNotFound
				} else {
					logger.Error(ctx, "error updating order", zap.String("id", orderId), zap.Error(err))

					return model.ErrInternal
				}
			}
			return nil
		case model.OrderStatusCANCELLED:
			logger.Warn(ctx, "error canclelling order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyCancelled))
			return model.ErrOrderAlreadyCancelled
		case model.OrderStatusPAID:
			logger.Warn(ctx, "error canclelling order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyPaid))
			return model.ErrOrderAlreadyPaid
		case model.OrderStatusCOMPLETED:
			logger.Warn(ctx, "error canclelling order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyCompleted))
			return model.ErrOrderAlreadyCompleted
		default:
			logger.Warn(ctx, "unexpected order status", zap.String("id", orderId), zap.Error(model.ErrUnexpectedOrderStatus))
			return model.ErrUnexpectedOrderStatus
		}
	})
	if err != nil {
		return err
	}

	logger.Info(
		ctx,
		"Order canceled",
		zap.String("order_id", orderId),
	)

	return nil
}
