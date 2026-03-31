package converter

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
)

func PartToDomain(part *inventory_v1.Part) *model.Part {
	metadataProto := make(map[string]model.Value)
	for k, v := range part.Metadata {
		metadataProto[k] = ValueToDomain(v)
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: int(part.StockQuantity),
		Category:      CategoryToDomain(part.Category),
		Dimensions:    DimensionsToDomain(part.Dimensions),
		Manufacturer:  ManufactererToDomain(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      metadataProto,
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     part.UpdatedAt.AsTime(),
	}
}
