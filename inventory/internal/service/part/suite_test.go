package part

import (
	"context"
	"testing"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context //nolint: containedctx

	inventoryRepository *mocks.InventoryRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.service = NewService(s.inventoryRepository)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
