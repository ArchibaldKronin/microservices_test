package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

type OrderService interface {
	CreateOrder(ctx context.Context, userId string, partIds []string) (order *model.Order, err error)
	CancelOrder(ctx context.Context, orderId string) (err error)
	GetOrder(ctx context.Context, orderId string) (order *model.Order, err error)
	PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (transId string, err error)
}
