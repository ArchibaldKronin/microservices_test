package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/service"
	txmanager "github.com/ArchibaldKronin/microservices_test/order/internal/txManager"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepo repository.OrderRepository
	txManager txmanager.TxManager

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(
	r repository.OrderRepository,
	txManager txmanager.TxManager,

	inventoryClient grpc.InventoryClient,
	paymentClient grpc.PaymentClient,
) *service {
	return &service{
		orderRepo: r,
		txManager: txManager,

		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
