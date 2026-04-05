package part

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func (r *repository) GetPart(_ context.Context, id string) (*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.data[id]
	if !ok {
		return nil, repoModel.ErrNotFound
	}

	return converter.PartToDomain(&part), nil
}
