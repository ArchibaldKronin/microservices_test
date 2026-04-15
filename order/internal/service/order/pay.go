package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
)

func (s *service) PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (string, error) {
	var result string

	err := s.txManager.WithTransaction(ctx, func(executer repository.OrderRepository) error {
		order, err := executer.GetOrder(ctx, orderId)
		if err != nil {
			slog.Error("error getting order", "error", err)

			if errors.Is(err, repoModel.ErrNotFound) {
				return model.ErrNotFound
			}
			return model.ErrInternal
		}

		userId := order.UserID
		transId, err := s.paymentClient.PayOrder(ctx, userId, orderId, pm)
		if err != nil {
			slog.Error("error payment", "error", err)
			return err
		}

		order.Status = model.OrderStatusPAID
		order.PaymentMethod = &pm
		order.TransactionID = &transId

		err = executer.UpdateOrder(ctx, order)
		if err != nil {
			if errors.Is(err, repoModel.ErrNotFound) {
				slog.Error("error NON CONSISTENT DATA", "error", err)
				return model.ErrNotFound
			} else {
				slog.Error("error updating order", "error", err)
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
