package part

import (
	"context"
	"errors"
	"log"
	"time"

	def "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	// mu   sync.RWMutex
	data *mongo.Collection
}

func NewRepository(db *mongo.Database) (*repository, error) {
	collection := db.Collection("parts")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "uuid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		return nil, err
	}

	err = initRepo(ctx, collection)
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

	log.Print("Init DB with initial parts")
	return nil

	// if res.Err() != nil {

	// 	if !errors.Is(res.Err(), mongo.ErrNoDocuments) {
	// 		return res.Err()
	// 	}

	// }

	// var initial []any
	// for _, v := range InitialParts {
	// 	initial = append(initial, v)
	// }
	// _, err := collection.InsertMany(ctx, initial)
	// if err != nil {
	// 	return err
	// }

	// return nil
}

// type repository struct {
// 	mu   sync.RWMutex
// 	data map[string]repoModel.Part
// }

// func NewRepository(init []*repoModel.Part) *repository {
// 	repo := make(map[string]repoModel.Part)
// 	for _, part := range init {
// 		temp := *part

// 		tags := make([]string, 0, len(part.Tags))
// 		tags = append(tags, part.Tags...)
// 		temp.Tags = tags

// 		metadata := make(map[string]repoModel.Value)
// 		for k, v := range part.Metadata {
// 			metadata[k] = v
// 		}
// 		temp.Metadata = metadata

// 		repo[part.Uuid] = temp
// 	}

// 	return &repository{
// 		data: repo,
// 	}
// }
