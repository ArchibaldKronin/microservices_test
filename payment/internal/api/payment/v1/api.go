package v1

import (
	"github.com/ArchibaldKronin/microservices_test/payment/internal/service"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
)

type api struct {
	payment_v1.UnimplementedPaymentServiceServer

	paymentService service.PaymentService
}

func NewAPI(paymentService service.PaymentService) *api {
	return &api{
		paymentService: paymentService,
	}
}
