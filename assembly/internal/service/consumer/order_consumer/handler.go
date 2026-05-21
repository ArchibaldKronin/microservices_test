package order_consumer

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/service/metrics"
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

	// засекаем время сборки
	startAssembling := time.Now()

	delayMillisec := time.Millisecond * time.Duration(bigInt.Int64())

	delaySec := int64(delayMillisec.Seconds())

	time.Sleep(delayMillisec)

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

	// останавливаем таймер сборки
	durationAssemnbling := time.Since(startAssembling)
	//записываем метрику
	metrics.AppendAssemblyDurationSecondsMetric(ctx, durationAssemnbling)
	return nil
}
