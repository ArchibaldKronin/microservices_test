package repository

import (
	"context"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o *serviceModel.Order) error
	UpdateOrder(ctx context.Context, o *serviceModel.Order) error
	GetOrder(ctx context.Context, id string) (*serviceModel.Order, error)
	// UpdateOrderTx(ctx context.Context, tx pgx.Tx, o *serviceModel.Order) error
	// GetOrderTx(ctx context.Context, tx pgx.Tx, id string) (*serviceModel.Order, error)

	// BeginTx(ctx context.Context) (pgx.Tx, error)
}
