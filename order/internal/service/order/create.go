package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/metrics"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/tracing"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	metrics.IncrOrdersTotalMetric(ctx)

	return order, nil
}

func (s *service) getPartsInfo(ctx context.Context, partIds []string) ([]string, float64, error) {
	// Создаем спан для вызова Inventory сервиса
	ctx, span := tracing.StartSpan(ctx, "order.call_inventory",
		trace.WithAttributes(
			attribute.StringSlice("parts_IDs", partIds)),
	)
	defer span.End()

	parts, err := s.inventoryClient.ListParts(
		ctx,
		model.PartsFilter{
			Uuids: partIds,
		},
	)
	if err != nil {
		span.RecordError(err)
		return nil, 0, err
	}

	if len(partIds) != 0 {
		if len(parts) != len(partIds) {
			span.RecordError(err)
			return nil, 0, model.ErrNotFound
		}
	}

	ids := make([]string, 0, len(parts))
	totalPrice := 0.0
	for _, part := range parts {
		ids = append(ids, part.Uuid)
		totalPrice += part.Price
	}

	span.SetAttributes(
		attribute.Float64("total_price", totalPrice),
	)

	return ids, totalPrice, nil
}
