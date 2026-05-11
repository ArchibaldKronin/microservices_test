package order_paid_consumer

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (c *service) orderPaidHandler(ctx context.Context, msg consumer.Message) error {
	event, err := c.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	err = c.telegramClient.SendOrderPaidNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send OrderPaid telegram notification", zap.String("OrderID", event.OrderUuid), zap.Error(err))
	}

	logger.Info(ctx, "Telegram notification OrderPaid sent", zap.String("OrderID", event.OrderUuid))

	return nil
}
