package order

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/tracing"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func (s *service) PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (string, error) {
	var result string

	ctx, span := tracing.StartSpan(ctx, "order.pay_order",
		trace.WithAttributes(
			attribute.String("order_id", orderId),
			attribute.String("payment_method", string(pm)),
		))
	defer span.End()

	err := s.txManager.WithTransaction(ctx, func(executer repository.OrderRepository) error {
		order, err := executer.GetOrder(ctx, orderId)
		if err != nil {
			logger.Error(ctx, "error getting order", zap.String("id", orderId), zap.Error(err))

			if errors.Is(err, repoModel.ErrNotFound) {
				return model.ErrNotFound
			}
			return model.ErrInternal
		}

		switch order.Status {
		case model.OrderStatusCANCELLED:
			logger.Warn(ctx, "error paing order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyCancelled))
			return model.ErrOrderAlreadyCancelled
		case model.OrderStatusPAID:
			logger.Warn(ctx, "error paing order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyPaid))
			return model.ErrOrderAlreadyPaid
		case model.OrderStatusCOMPLETED:
			logger.Warn(ctx, "error paing order", zap.String("id", orderId), zap.Error(model.ErrOrderAlreadyCompleted))
			return model.ErrOrderAlreadyCompleted
		case model.OrderStatusPENDINGPAYMENT:
		default:
			logger.Warn(ctx, "unexpected order status", zap.String("id", orderId), zap.Error(model.ErrUnexpectedOrderStatus))
			return model.ErrUnexpectedOrderStatus
		}

		userId := order.UserID
		transId, err := s.paymentClient.PayOrder(ctx, userId, orderId, pm)
		if err != nil {
			logger.Error(ctx, "error paying order", zap.String("id", orderId), zap.String("payment_method", string(pm)), zap.Error(err))
			return err
		}

		order.Status = model.OrderStatusPAID
		order.PaymentMethod = &pm
		order.TransactionID = &transId

		err = executer.UpdateOrder(ctx, order)
		if err != nil {
			if errors.Is(err, repoModel.ErrNotFound) {
				logger.Error(ctx, "error NON CONSISTENT DATA", zap.String("id", orderId), zap.Error(err))
				return model.ErrNotFound
			} else {
				logger.Error(ctx, "error updating order", zap.String("id", orderId), zap.Error(err))
				return model.ErrInternal
			}
		}

		eventID := uuid.NewString()
		err = s.orderProducerService.ProduceOrderPaid(ctx, model.OrderPaidEvent{
			EventUuid:       eventID,
			OrderUuid:       order.OrderID,
			UserUuid:        order.UserID,
			PaymentMethod:   pm,
			TransactionUuid: transId,
		})
		if err != nil {
			logger.Error(ctx, "error producing ", zap.String("id", orderId), zap.String("eventID", eventID), zap.Error(err))
			return model.ErrInternal
		}

		result = transId
		return nil
	})
	if err != nil {
		span.RecordError(err)
		return "", err
	}

	logger.Info(
		ctx,
		"Order paid",
		zap.String("order_id", orderId),
	)

	return result, nil
}
