package v1

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/client/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/middleware/grpc"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, userId, orderId string, paymentMethod model.PaymentMethod) (transactionId string, err error) {
	ctx = grpc.ForwardSessionUUIDToGRPC(ctx)

	resp, err := c.generatedClient.PayOrder(ctx,
		&payment_v1.PayOrderRequest{
			UserUuid:      userId,
			OrderUuid:     orderId,
			PaymentMethod: converter.PaymentMethodToDTO(paymentMethod),
		},
	)
	if err != nil {
		return "", converter.MapError(err)
	}
	return resp.TransactionUuid, nil
}
