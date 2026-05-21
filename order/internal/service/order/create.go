package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) CreateOrder(ctx context.Context, userId string, partIds []string) (order *model.Order, err error) {
	ids, totalPrice, err := s.getPartsInfo(ctx, partIds)
	if err != nil {
		logger.Error(ctx, "error get parts info to create order", zap.Strings("parts_ids", partIds), zap.Error(err))

		return nil, err
	}

	order = model.NewOrder(userId, ids, totalPrice)
	if err = s.orderRepo.CreateOrder(ctx, order); err != nil {

		logger.Error(ctx, "error creating order", zap.Strings("parts_ids", partIds), zap.Error(err))

		return nil, model.ErrInternal
	}

	logger.Info(
		ctx,
		"Order created",
		zap.String("order_id", order.OrderID),
	)

	return order, nil
}

func (s *service) getPartsInfo(ctx context.Context, partIds []string) ([]string, float64, error) {
	parts, err := s.inventoryClient.ListParts(
		ctx,
		model.PartsFilter{
			Uuids: partIds,
		},
	)
	if err != nil {
		return nil, 0, err
	}

	if len(partIds) != 0 {
		if len(parts) != len(partIds) {
			return nil, 0, model.ErrNotFound
		}
	}

	ids := make([]string, 0, len(parts))
	totalPrice := 0.0
	for _, part := range parts {
		ids = append(ids, part.Uuid)
		totalPrice += part.Price
	}

	return ids, totalPrice, nil
}
