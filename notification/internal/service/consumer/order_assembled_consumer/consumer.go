package order_assembled_consumer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/notification/internal/converter/kafka"
	def "github.com/ArchibaldKronin/microservices_test/notification/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
)

var _ def.OrderAssembledConsumerService = (*service)(nil)

type service struct {
	orderAssembledConsumer kafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder

	telegramClient def.TelegramService
}

func NewService(
	orderAssembledConsumer kafka.Consumer,
	orderAssembledDecoder kafkaConverter.OrderAssembledDecoder,
	telegramClient def.TelegramService,
) *service {
	return &service{
		orderAssembledConsumer: orderAssembledConsumer,
		orderAssembledDecoder:  orderAssembledDecoder,
		telegramClient:         telegramClient,
	}
}

func (c *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "starting notification OrderAssembledConsumer service")

	err := c.orderAssembledConsumer.Consume(ctx, c.orderAssembledHandler)
	if err != nil {
		return err
	}

	return nil
}
