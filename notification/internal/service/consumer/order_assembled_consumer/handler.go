package order_assembled_consumer

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (c *service) orderAssembledHandler(ctx context.Context, msg consumer.Message) error {
	event, err := c.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderAssembled", zap.Error(err))
		return err
	}

	err = c.telegramClient.SendOrderAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send OrderAssembled telegram notification", zap.String("OrderID", event.OrderUuid), zap.Error(err))
	}

	logger.Info(ctx, "Telegram notification OrderAssembled sent", zap.String("OrderID", event.OrderUuid))

	return nil
}
