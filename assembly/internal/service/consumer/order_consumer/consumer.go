package order_consumer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/assembly/internal/converter/kafka"
	def "github.com/ArchibaldKronin/microservices_test/assembly/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.OrderPaidConsumer = (*service)(nil)

type service struct {
	orderPaidConsumer kafka.Consumer
	orderPaidDecoder  kafkaConverter.OrderPaidDecoder

	assembleProducerService def.ShipAssembledProducer
}

func NewService(
	orderPaidConsumer kafka.Consumer,
	orderPaidDecoder kafkaConverter.OrderPaidDecoder,
	assembleProducerService def.ShipAssembledProducer,
) *service {
	return &service{
		orderPaidDecoder:        orderPaidDecoder,
		orderPaidConsumer:       orderPaidConsumer,
		assembleProducerService: assembleProducerService,
	}
}

func (c *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "starting assembly OrderPaidConsumer service")

	err := c.orderPaidConsumer.Consume(ctx, c.orderPaidHandler)
	if err != nil {
		logger.Error(ctx, "consume from order.paid topic error", zap.Error(err))
		return err
	}

	return nil
}
