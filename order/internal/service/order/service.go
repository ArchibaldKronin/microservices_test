package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	r repository.OrderRepository,

	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		orderRepo: r,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
