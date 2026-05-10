package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
)

type TelegramService interface {
	SendOrderPaidNotification(ctx context.Context, orderPaid model.OrderPaidEvent) error
	SendOrderAssembledNotification(ctx context.Context, orderAssembled model.OrderAssembledEvent) error
}

type OrderAssembledConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderPaidConsumerService interface {
	RunConsumer(ctx context.Context) error
}
