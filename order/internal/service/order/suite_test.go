package order

// import (
// 	"context"
// 	"testing"

// 	clientMocks "github.com/ArchibaldKronin/microservices_test/order/internal/client/grpc/mocks"
// 	orderMock "github.com/ArchibaldKronin/microservices_test/order/internal/repository/mocks"
// 	txmanagerMock "github.com/ArchibaldKronin/microservices_test/order/internal/txManager/mocks"
// 	"github.com/stretchr/testify/suite"
// )

// type ServiceSuite struct {
// 	suite.Suite

// 	ctx context.Context //nolint: containedctx

// 	txManager       *txmanagerMock.TxManager
// 	orderRepository *orderMock.OrderRepository
// 	paymentClient   *clientMocks.PaymentClient
// 	inventoryClient *clientMocks.InventoryClient

// 	service *service
// }

// func (s *ServiceSuite) SetupTest() {
// 	s.ctx = context.Background()

// 	s.txManager = txmanagerMock.NewTxManager(s.T())
// 	s.orderRepository = orderMock.NewOrderRepository(s.T())

// 	s.paymentClient = clientMocks.NewPaymentClient(s.T())
// 	s.inventoryClient = clientMocks.NewInventoryClient(s.T())

// 	// s.service = NewService(s.orderRepository, s.txManager, s.inventoryClient, s.paymentClient)
// }

// func (s *ServiceSuite) TearDownTest() {}

// func TestServiceIntegration(t *testing.T) {
// 	suite.Run(t, new(ServiceSuite))
// }
