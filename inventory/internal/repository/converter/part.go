package converter

import (
	"slices"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func PartToDomain(p *repoModel.Part) *model.Part {
	categoryProto := CategoryToDomain(p.Category)
	dimensionsProto := DimensionsToDomain(p.Dimensions)
	manufacturerProto := ManufactererToDomain(p.Manufacturer)

	metadataProto := make(map[string]model.Value)
	for k, v := range p.Metadata {
		metadataProto[k] = ValueToDomain(v)
	}

	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      categoryProto,
		Dimensions:    dimensionsProto,
		Manufacturer:  manufacturerProto,
		Tags:          slices.Clone(p.Tags),
		Metadata:      metadataProto,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}
}
