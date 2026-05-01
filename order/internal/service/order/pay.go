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

func (s *service) PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (string, error) {
	var result string

	err := s.txManager.WithTransaction(ctx, func(executer repository.OrderRepository) error {
		order, err := executer.GetOrder(ctx, orderId)
		if err != nil {
			logger.Error(ctx, "error getting order", zap.String("id", orderId), zap.Error(err))

			if errors.Is(err, repoModel.ErrNotFound) {
				return model.ErrNotFound
			}
			return model.ErrInternal
		}

		userId := order.UserID
		transId, err := s.paymentClient.PayOrder(ctx, userId, orderId, pm)
		if err != nil {
			logger.Error(ctx, "error paying order", zap.String("id", orderId), zap.String("payment_method", string(pm)), zap.Error(err))
			return err
		}

		order.Status = model.OrderStatusPAID
		order.PaymentMethod = &pm
		order.TransactionID = &transId

		err = executer.UpdateOrder(ctx, order)
		if err != nil {
			if errors.Is(err, repoModel.ErrNotFound) {
				logger.Error(ctx, "error NON CONSISTENT DATA", zap.String("id", orderId), zap.Error(err))
				return model.ErrNotFound
			} else {
				logger.Error(ctx, "eerror updating order", zap.String("id", orderId), zap.Error(err))
				return model.ErrInternal
			}
		}

		result = transId
		return nil
	})
	if err != nil {
		return "", err
	}
	return result, nil
}
