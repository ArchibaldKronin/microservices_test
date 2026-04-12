package order

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
)

func (s *service) GetOrder(ctx context.Context, orderId string) (order *model.Order, err error) {
	order, err = s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		slog.Error("error getting order", "error", err)

		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, model.ErrInternal
	}

	return order, nil
}
