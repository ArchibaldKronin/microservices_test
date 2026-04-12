package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (s *service) CancelOrder(ctx context.Context, orderId string) error {
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		slog.Error("error beginning tx", "error", err)
		return model.ErrInternal
	}

	defer func() {
		rerr := tx.Rollback(ctx)
		if rerr != nil && !errors.Is(rerr, pgx.ErrTxClosed) {
			slog.Warn("rollback failed", "error", rerr)
		}
	}()

	order, err := s.orderRepo.GetOrderTx(ctx, tx, orderId)
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
		err = s.orderRepo.UpdateOrderTx(ctx, tx, order)
		if err != nil {
			slog.Error("error updating order", "error", err)

			if errors.Is(err, repoModel.ErrNotFound) {
				return model.ErrNotFound
			}
			return model.ErrInternal
		}

		err = tx.Commit(ctx)
		if err != nil {
			slog.Error("error committing tx:", "error", err)
			return model.ErrInternal
		}

		return nil

	case model.OrderStatusPAID:
		return model.ErrOrderAlreadyPaid

		// already cancelled
	default:
		return nil
	}
}
