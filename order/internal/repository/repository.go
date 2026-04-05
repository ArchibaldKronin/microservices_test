package repository

import (
	"context"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

type OrderRepository interface {
	CreateOrder(ctx context.Context, o *serviceModel.Order)
	UpdateOrder(ctx context.Context, o *serviceModel.Order) *serviceModel.Order
	GetOrder(ctx context.Context, id string) *serviceModel.Order
}
