package v1

import (
	"context"
	"testing"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type ApiSuite struct {
	suite.Suite

	ctx context.Context //nolint: containedctx

	inventoryService *mocks.InventoryService

	api *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()

	a.inventoryService = mocks.NewInventoryService(a.T())

	a.api = NewApi(a.inventoryService)
}

func (a *ApiSuite) TearDownTest() {}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
