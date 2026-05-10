package order_paid_consumer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/notification/internal/converter/kafka"
	def "github.com/ArchibaldKronin/microservices_test/notification/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
)

var _ def.OrderPaidConsumerService = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	telegramClient def.TelegramService
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	telegramClient def.TelegramService,
) *service {
	return &service{
		orderPaidConsumer: orderPaidConsumer,
		orderPaidDecoder:  orderPaidDecoder,
		telegramClient:    telegramClient,
	}
}

func (c *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "starting notification OrderPaidConsumer service")

	err := c.orderPaidConsumer.Consume(ctx, c.orderPaidHandler)
	if err != nil {
		return err
	}

	return nil
}
