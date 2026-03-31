package payment

import (
	def "github.com/ArchibaldKronin/microservices_test/payment/internal/service"
)

var _ def.PaymentService = (*service)(nil)

type service struct{}

func NewService() *service {
	return &service{}
}
