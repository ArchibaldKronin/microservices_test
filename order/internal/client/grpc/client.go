package grpc

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

type PaymentClient interface {
	PayOrder(ctx context.Context, userId, orderId string, paymentMethod model.PaymentMethod) (transactionId string, err error)
}

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error)
}
