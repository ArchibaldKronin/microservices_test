package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
)

type InventoryService interface {
	GetPart(context.Context, string) (*model.Part, error)
	ListParts(context.Context, *model.PartsFilter) ([]*model.Part, error)
}
