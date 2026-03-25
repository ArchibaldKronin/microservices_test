package main

import (
	"context"

	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type InventoryClient interface {
	GetPart(ctx context.Context, uuid string) (*inventory_v1.Part, error)
	ListParts(ctx context.Context, uuids []string) ([]*inventory_v1.Part, error)
}

type inventoryClient struct {
	client inventory_v1.InventoryServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) InventoryClient {
	client := inventory_v1.NewInventoryServiceClient(conn)

	return &inventoryClient{
		client: client,
	}
}

func (c *inventoryClient) GetPart(ctx context.Context, uuid string) (*inventory_v1.Part, error) {
	resp, err := c.client.GetPart(ctx,
		&inventory_v1.GetPartRequest{
			Uuid: uuid,
		})
	if err != nil {
		return nil, err
	}
	return resp.Part, nil
}

func (c *inventoryClient) ListParts(ctx context.Context, uuids []string) ([]*inventory_v1.Part, error) {
	resp, err := c.client.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: &inventory_v1.PartsFilter{
			Uuids: uuids,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}
