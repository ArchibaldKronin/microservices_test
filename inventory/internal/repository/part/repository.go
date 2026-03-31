package part

import (
	"sync"

	def "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	repo := make(map[string]repoModel.Part)
	for _, part := range InitialParts {
		repo[part.Uuid] = *part
	}

	return &repository{
		data: repo,
	}
}
