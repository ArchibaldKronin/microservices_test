package v1

import (
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/service"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventory_v1.UnimplementedInventoryServiceServer

	inventoryService service.InventoryService
}

func NewApi(invService service.InventoryService) *api {
	return &api{
		inventoryService: invService,
	}
}
