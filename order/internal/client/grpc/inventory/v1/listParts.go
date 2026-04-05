package v1

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/client/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"github.com/samber/lo"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]*model.Part, error) {
	resp, err := c.generatedClient.ListParts(ctx, &inventory_v1.ListPartsRequest{
		Filter: lo.ToPtr(converter.FilterToDTO(filter)),
	})
	if err != nil {
		return nil, converter.MapError(err)
	}

	parts := make([]*model.Part, 0, len(resp.Parts))
	for _, p := range resp.Parts {
		parts = append(parts, converter.PartToDomain(p))
	}
	return parts, nil
}
