package order

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetOrder(ctx context.Context, orderId string) (order *model.Order, err error) {
	order, err = s.orderRepo.GetOrder(ctx, orderId)
	if err != nil {
		logger.Error(ctx, "error getting order", zap.String("id", orderId), zap.Error(err))

		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, model.ErrInternal
	}

	return order, nil
}
