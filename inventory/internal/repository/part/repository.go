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

func NewRepository(init []*repoModel.Part) *repository {
	repo := make(map[string]repoModel.Part)
	for _, part := range init {
		temp := *part

		tags := make([]string, 0, len(part.Tags))
		tags = append(tags, part.Tags...)
		temp.Tags = tags

		metadata := make(map[string]repoModel.Value)
		for k, v := range part.Metadata {
			metadata[k] = v
		}
		temp.Metadata = metadata

		repo[part.Uuid] = temp
	}

	return &repository{
		data: repo,
	}
}
