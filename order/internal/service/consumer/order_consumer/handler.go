package order_consumer

import (
	"context"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) ShipAssembledHandler(ctx context.Context, msg consumer.Message) error {
	event, err := s.orderAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode ShipAssembled", zap.Error(err))
		return err
	}

	err = s.orderService.CompleteOrder(ctx, event.OrderUuid)
	if err != nil {
		logger.Error(ctx, "error in consumer handler", zap.Error(err))
		return fmt.Errorf("error in consumer handler: %w", err)
	}

	return nil
}
