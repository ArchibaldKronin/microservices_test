package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
)

type ShipAssembledProducer interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}

type OrderPaidConsumer interface {
	RunConsumer(ctx context.Context) error
}
