package metrics

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

const (
	serviceName = "order-service"
)

var meter = otel.Meter(serviceName)

var (
	// ordersTotal Счетчик общего числа заказов
	ordersTotal metric.Int64Counter

	// Суммарная выручка
	ordersRevenueTotal metric.Float64Counter
)

func InitMetrics() error {
	var err error

	ordersTotal, err = meter.Int64Counter(
		"orders_create_total",
		metric.WithDescription("Total number of created Orders"),
	)
	if err != nil {
		return err
	}

	ordersRevenueTotal, err = meter.Float64Counter(
		"orders_total_revenue",
		metric.WithDescription("Total revenue"),
	)
	if err != nil {
		return err
	}

	return nil
}

func IncrOrdersTotalMetric(ctx context.Context) {
	ordersTotal.Add(ctx, 1)
}

func IncrOrdersRevenueTotalMetric(ctx context.Context, sum float64) {
	ordersRevenueTotal.Add(ctx, sum)
}
