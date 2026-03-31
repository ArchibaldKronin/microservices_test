package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/payment/internal/model"
)

type PaymentService interface {
	PayOrder(context.Context, model.PaymentMethod) (string, error)
}
