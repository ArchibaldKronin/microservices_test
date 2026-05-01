package v1

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/converter.go"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *payment_v1.PayOrderRequest) (*payment_v1.PayOrderResponse, error) {
	pm := converter.PaymentMethodToDomain(req.PaymentMethod)
	transactionId, err := a.paymentService.PayOrder(ctx, pm)
	if err != nil {
		logger.Error(ctx, "internal error when payd order", zap.String("order_ID", req.OrderUuid))
		return nil, status.Errorf(codes.Internal, "internal error when payd order %s", req.OrderUuid)
	}
	return &payment_v1.PayOrderResponse{
		TransactionUuid: transactionId,
	}, nil
}
