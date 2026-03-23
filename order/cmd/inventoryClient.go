package main

import (
	"context"

	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
)

type InventoryClient interface {
	GetPart(ctx context.Context, uuid string) (*inventory_v1.Part, error)
	ListPart(ctx context.Context, uuids []string) ([]*inventory_v1.Part, error)
}

type inventoryClient struct {
	client inventory_v1.PartServiceClient
}

func NewInventoryClient(conn *grpc.ClientConn) *inventoryClient {
	client := inventory_v1.NewPartServiceClient(conn)

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

func (c *inventoryClient) ListPart(ctx context.Context, uuids []string) ([]*inventory_v1.Part, error) {
	resp, err := c.client.ListPart(ctx, &inventory_v1.ListPartRequest{
		Uuids: uuids,
	})
	if err != nil {
		return nil, err
	}
	return resp.Parts, nil
}
