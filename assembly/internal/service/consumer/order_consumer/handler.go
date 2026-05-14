package order_consumer

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/kafka/consumer"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (c *service) orderPaidHandler(ctx context.Context, msg consumer.Message) error {
	event, err := c.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode OrderPaid", zap.Error(err))
		return err
	}

	bigInt, err := rand.Int(rand.Reader, big.NewInt(10000))
	if err != nil {
		logger.Error(ctx, "Failed to generate random bigInt", zap.Error(err))
		return err
	}
	delayMillisec := time.Millisecond * time.Duration(bigInt.Int64())
	ctx, cancel := context.WithTimeout(ctx, delayMillisec)
	defer cancel()

	delaySec := int64(delayMillisec.Seconds())
	// delaySec := int64(math.Round(float64(delayMillisec) / 1000_000_000))

	<-ctx.Done()

	eventId := uuid.NewString()
	err = c.assembleProducerService.ProduceShipAssembled(ctx, model.ShipAssembledEvent{
		EventUuid:    eventId,
		OrderUuid:    event.OrderUuid,
		UserUuid:     event.UserUuid,
		BuildTimeSec: delaySec,
	})
	if err != nil {
		logger.Error(ctx, "error in consumer handler", zap.Error(err))
		return fmt.Errorf("error in consumer handler: %w", err)
	}

	return nil
}
