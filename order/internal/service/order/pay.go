package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/jackc/pgx/v5"
)

func (s *service) PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (string, error) {
	tx, err := s.orderRepo.BeginTx(ctx)
	if err != nil {
		slog.Error("error beginning tx", "error", err)
		return "", model.ErrInternal
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
			return "", model.ErrNotFound
		}
		return "", model.ErrInternal
	}

	userId := order.UserID
	transId, err := s.paymentClient.PayOrder(ctx, userId, orderId, pm)
	if err != nil {
		slog.Error("error payment", "error", err)
		return "", err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &pm
	order.TransactionID = &transId

	err = s.orderRepo.UpdateOrderTx(ctx, tx, order)
	if err != nil {
		slog.Error("error updating order", "error", err)

		if errors.Is(err, repoModel.ErrNotFound) {
			return "", model.ErrNotFound
		}
		return "", model.ErrInternal
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error committing tx:", "error", err)
		return "", model.ErrInternal
	}

	return transId, nil
}
