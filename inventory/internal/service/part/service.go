package part

import (
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository"
	def "github.com/ArchibaldKronin/microservices_test/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	repo repository.InventoryRepository
}

func NewService(r repository.InventoryRepository) *service {
	return &service{
		repo: r,
	}
}
