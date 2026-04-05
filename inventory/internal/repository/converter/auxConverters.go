package converter

import (
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func CategoryToRepo(c model.Category) repoModel.Category {
	switch c {
	case 1:
		return repoModel.CategoryEngine
	case 2:
		return repoModel.CategoryFuel
	case 3:
		return repoModel.CategoryPorthole
	case 4:
		return repoModel.CategoryWing
	default:
		return repoModel.CategoryUnknown
	}
}

func CategoryToDomain(c repoModel.Category) model.Category {
	switch c {
	case 1:
		return model.CategoryEngine
	case 2:
		return model.CategoryFuel
	case 3:
		return model.CategoryPorthole
	case 4:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func DimensionsToDomain(d repoModel.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufactererToDomain(m repoModel.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ValueToDomain(val repoModel.Value) model.Value {
	switch v := val.(type) {
	case repoModel.StringValue:
		return model.StringValue{
			Value: v.Value,
		}
	case repoModel.Int64Value:
		return model.Int64Value{
			Value: v.Value,
		}
	case repoModel.DoubleValue:
		return model.DoubleValue{
			Value: v.Value,
		}
	case repoModel.BoolValue:
		return model.BoolValue{
			Value: v.Value,
		}
	default:
		return nil
	}
}
