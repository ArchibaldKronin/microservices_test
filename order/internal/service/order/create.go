package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

func (s *service) CreateOrder(ctx context.Context, userId string, partIds []string) (order *model.Order, err error) {
	ids, totalPrice, err := s.getPartsInfo(ctx, partIds)
	if err != nil {
		return nil, err
	}

	order = model.NewOrder(userId, ids, totalPrice)
	s.orderRepo.CreateOrder(ctx, order)

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
