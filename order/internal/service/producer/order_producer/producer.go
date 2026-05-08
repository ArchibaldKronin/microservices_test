package order_producer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/order/internal/converter/kafka"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.OrderProducerService = (*service)(nil)

type service struct {
	orderPaidProducer kafka.Producer
}

func NewService(orderPaidProducer kafka.Producer) *service {
	return &service{
		orderPaidProducer: orderPaidProducer,
	}
}

func (p *service) ProduceOrderPaid(ctx context.Context, event model.OrderPaidEvent) error {
	msg := kafkaConverter.ModelToOrderPaidEvent(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderPaid", zap.Error(err))
		return err
	}

	err = p.orderPaidProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish OrderPaid", zap.Error(err))
		return err
	}

	return nil
}
