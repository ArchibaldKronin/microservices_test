package order_producer

import (
	"context"

	kafkaConverter "github.com/ArchibaldKronin/microservices_test/assembly/internal/converter/kafka"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/assembly/internal/service"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

var _ def.ShipAssembledProducer = (*servive)(nil)

type servive struct {
	shipAssembledProducer kafka.Producer
}

func NewService(shipAssembledProducer kafka.Producer) *servive {
	return &servive{
		shipAssembledProducer: shipAssembledProducer,
	}
}

func (p *servive) ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error {
	msg := kafkaConverter.ModelToShipAssembledEvent(event)

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal ShipAssembled", zap.Error(err))
		return err
	}

	err = p.shipAssembledProducer.Send(ctx, []byte(event.EventUuid), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish ShipAssembled", zap.Error(err))
		return err
	}

	return nil
}
