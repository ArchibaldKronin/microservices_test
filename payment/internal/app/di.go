package app

import (
	"context"

	v1 "github.com/ArchibaldKronin/microservices_test/payment/internal/api/payment/v1"
	"github.com/ArchibaldKronin/microservices_test/payment/internal/service"
	paymentService "github.com/ArchibaldKronin/microservices_test/payment/internal/service/payment"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API payment_v1.PaymentServiceServer

	paymentService service.PaymentService
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) PaymentV1API(ctx context.Context) payment_v1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = v1.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
