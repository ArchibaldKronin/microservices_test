package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
)

func (s *service) CancelOrder(ctx context.Context, orderId string) error {
	err := s.txManager.WithTransaction(ctx, func(executer repository.OrderRepository) error {
		order, err := executer.GetOrder(ctx, orderId)
		if err != nil {
			slog.Error("error getting order", "error", err)

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
					slog.Error("error NON CONSISTENT DATA", "error", err)
					return model.ErrNotFound
				} else {
					slog.Error("error updating order", "error", err)
					return model.ErrInternal
				}
			}
			return nil
		case model.OrderStatusPAID:
			return model.ErrOrderAlreadyPaid
		default:
			return nil
		}
	})

	if err != nil {
		return err
	}
	return nil
}
