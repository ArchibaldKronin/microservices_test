package app

import (
	"context"

	v1 "github.com/ArchibaldKronin/microservices_test/inventory/internal/api/inventory/v1"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/config"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/repository"
	inventoryRepository "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/part"
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/service"
	inventoryService "github.com/ArchibaldKronin/microservices_test/inventory/internal/service/part"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type diContainer struct {
	inventoryV1API inventory_v1.InventoryServiceServer

	inventoryService service.InventoryService

	inventoryRepository repository.InventoryRepository

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) (inventory_v1.InventoryServiceServer, error) {
	if d.inventoryV1API == nil {
		service, err := d.InventoryService(ctx)
		if err != nil {
			return nil, err
		}

		d.inventoryV1API = v1.NewApi(service)
	}

	return d.inventoryV1API, nil
}

func (d *diContainer) InventoryService(ctx context.Context) (service.InventoryService, error) {
	if d.inventoryService == nil {
		repo, err := d.InventoryRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.inventoryService = inventoryService.NewService(repo)
	}

	return d.inventoryService, nil
}

func (d *diContainer) InventoryRepository(ctx context.Context) (repository.InventoryRepository, error) {
	if d.inventoryRepository == nil {
		db, err := d.MongoDBHandle(ctx)
		if err != nil {
			return nil, err
		}

		repo, err := inventoryRepository.NewRepository(ctx, db)
		if err != nil {
			return nil, err
		}

		d.inventoryRepository = repo
	}

	return d.inventoryRepository, nil
}

func (d *diContainer) MongoDBHandle(ctx context.Context) (*mongo.Database, error) {
	if d.mongoDBHandle == nil {
		mongoClient, err := d.MongoDBClient(ctx)
		if err != nil {
			return nil, err
		}

		mongoDB := mongoClient.Database(config.AppConfig().Mongo.DatabaseName())

		d.mongoDBHandle = mongoDB
	}

	return d.mongoDBHandle, nil
}

func (d *diContainer) MongoDBClient(ctx context.Context) (*mongo.Client, error) {
	if d.mongoDBClient == nil {
		mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			return nil, err
		}

		err = mongoClient.Ping(ctx, nil)
		if err != nil {
			return nil, err
		}

		d.mongoDBClient = mongoClient
	}

	return d.mongoDBClient, nil
}
