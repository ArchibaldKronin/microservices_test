package converter

import (
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
)

func PartToProto(p *model.Part) *inventory_v1.Part {
	categoryProto := CategoryToProto(p.Category)
	dimensionsProto := DimensionsToProto(p.Dimensions)
	manufacturerProto := ManufactererToProto(p.Manufacturer)

	metadataProto := make(map[string]*inventory_v1.Value)
	for k, v := range p.Metadata {
		metadataProto[k] = ValueToProto(v)
	}

	return &inventory_v1.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: int64(p.StockQuantity),
		Category:      categoryProto,
		Dimensions:    &dimensionsProto,
		Manufacturer:  &manufacturerProto,
		Tags:          p.Tags,
		Metadata:      metadataProto,
		CreatedAt:     TimeToProto(&p.CreatedAt),
		UpdatedAt:     TimeToProto(p.UpdatedAt),
	}
}
