package part

import (
	"context"
	"errors"
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (r *repository) GetPart(ctx context.Context, id string) (*model.Part, error) {
	var part repoModel.Part

	err := r.data.FindOne(ctx, bson.M{"uuid": id}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, repoModel.ErrNotFound
		}
		return nil, fmt.Errorf("find part by id %s: %w", id, err)
	}

	partDomain, err := converter.PartToDomain(&part)
	if err != nil {
		return partDomain, fmt.Errorf("error get part by id %s: %w", id, err)
	}

	return partDomain, nil
}

// func (r *repository) GetPart(_ context.Context, id string) (*model.Part, error) {
// 	r.mu.RLock()
// 	defer r.mu.RUnlock()

// 	part, ok := r.data[id]
// 	if !ok {
// 		return nil, repoModel.ErrNotFound
// 	}

// 	return converter.PartToDomain(&part), nil
// }
