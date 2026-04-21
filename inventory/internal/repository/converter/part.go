package converter

import (
	"slices"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func PartToDomain(p *repoModel.Part) (*model.Part, error) {
	var err error = nil

	category := CategoryToDomain(p.Category)
	dimensions := DimensionsToDomain(p.Dimensions)
	manufacturer := ManufactererToDomain(p.Manufacturer)

	metadata := make(map[string]model.Value)
	for k, v := range p.Metadata {
		value, errConv := PrimitiveToValue(v)
		if errConv != nil {
			if err == nil {
				err = errConv
			}
			continue
		}

		metadata[k] = value
	}

	return &model.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          slices.Clone(p.Tags),
		Metadata:      metadata,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, err
}

func PartToRepo(p *model.Part) (*repoModel.Part, error) {
	var err error = nil
	// handErr := false

	category := CategoryToRepo(p.Category)
	dimensions := DimensionsToRepo(p.Dimensions)
	manufacturer := ManufactererToRepo(p.Manufacturer)

	metadata := make(map[string]any)
	for k, v := range p.Metadata {
		value, errConv := ValueToPrimitive(v)
		if errConv != nil {
			if err == nil {
				err = errConv
			}
			continue
		}
		// var value any
		// if !handErr {
		// 	value, err = ValueToPrimitive(v)
		// 	if err != nil {
		// 		handErr = true
		// 		continue
		// 	}
		// } else {
		// 	var errNext error
		// 	value, errNext = ValueToPrimitive(v)
		// 	if errNext != nil {
		// 		continue
		// 	}
		// }
		metadata[k] = value
	}

	return &repoModel.Part{
		Uuid:          p.Uuid,
		Name:          p.Name,
		Description:   p.Description,
		Price:         p.Price,
		StockQuantity: p.StockQuantity,
		Category:      category,
		Dimensions:    dimensions,
		Manufacturer:  manufacturer,
		Tags:          slices.Clone(p.Tags),
		Metadata:      metadata,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, err
}
