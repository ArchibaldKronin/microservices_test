package converter

import (
	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
)

func ValueToPrimitive(v model.Value) (any, error) {
	switch val := v.(type) {
	case model.StringValue:
		return val.Value, nil
	case model.Int64Value:
		return val.Value, nil
	case model.DoubleValue:
		return val.Value, nil
	case model.BoolValue:
		return val.Value, nil
	default:
		return nil, repoModel.NewMetadataParseValueError(v, nil)
	}
}

func PrimitiveToValue(v any) (model.Value, error) {
	switch val := v.(type) {
	case string:
		return model.StringValue{Value: val}, nil
	case int:
		return model.Int64Value{Value: int64(val)}, nil
	case int32:
		return model.Int64Value{Value: int64(val)}, nil
	case int64:
		return model.Int64Value{Value: val}, nil
	case float32:
		return model.DoubleValue{Value: float64(val)}, nil
	case float64:
		return model.DoubleValue{Value: val}, nil
	case bool:
		return model.BoolValue{Value: val}, nil
	default:
		return nil, repoModel.NewMetadataParseValueError(v, nil)
	}
}

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

func DimensionsToRepo(d model.Dimensions) repoModel.Dimensions {
	return repoModel.Dimensions{
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

func ManufactererToRepo(m model.Manufacturer) repoModel.Manufacturer {
	return repoModel.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

// func ValueToDomain(val repoModel.Value) model.Value {
// 	switch v := val.(type) {
// 	case repoModel.StringValue:
// 		return model.StringValue{
// 			Value: v.Value,
// 		}
// 	case repoModel.Int64Value:
// 		return model.Int64Value{
// 			Value: v.Value,
// 		}
// 	case repoModel.DoubleValue:
// 		return model.DoubleValue{
// 			Value: v.Value,
// 		}
// 	case repoModel.BoolValue:
// 		return model.BoolValue{
// 			Value: v.Value,
// 		}
// 	default:
// 		return nil
// 	}
// }
