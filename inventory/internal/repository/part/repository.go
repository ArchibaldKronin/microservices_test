package part

import (
	"context"
	"errors"
	"time"

	def "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	// mu   sync.RWMutex
	data *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) (*repository, error) {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	constructorCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(constructorCtx, indexModels)
	if err != nil {
		return nil, err
	}

	err = initRepo(constructorCtx, collection)
	if err != nil {
		return nil, err
	}

	return &repository{
		data: collection,
	}, nil
}

func initRepo(ctx context.Context, collection *mongo.Collection) error {
	res := collection.FindOne(ctx, bson.M{})
	if res.Err() == nil {
		return nil
	}

	if !errors.Is(res.Err(), mongo.ErrNoDocuments) {
		return res.Err()
	}

	var initial []any
	for _, v := range InitialParts {
		initial = append(initial, v)
	}
	_, err := collection.InsertMany(ctx, initial)
	if err != nil {
		return err
	}

	logger.Info(ctx, "Init DB with initial parts")
	return nil
}
