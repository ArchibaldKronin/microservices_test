package order_consumer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/order/internal/converter/kafka"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

var _ def.OrderConsumerService = (*service)(nil)

type service struct {
	orderAssebledConsumer kafka.Consumer
	orderAssembledDecoder kafkaConverter.ShipAssembledDecoder

	orderService def.OrderService
}

func NewService(orderAssebledConsumer kafka.Consumer, orderAssembledDecoder kafkaConverter.ShipAssembledDecoder, orderService def.OrderService) *service {
	return &service{
		orderAssebledConsumer: orderAssebledConsumer,
		orderAssembledDecoder: orderAssembledDecoder,

		orderService: orderService,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order ShipAssembledConsumer service")

	err := s.orderAssebledConsumer.Consume(ctx, s.shipAssembledHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
