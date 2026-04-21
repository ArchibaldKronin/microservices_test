package converter

import (
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CategoryToProto(c model.Category) inventory_v1.Category {
	switch c {
	case 1:
		return inventory_v1.Category_CATEGORY_ENGINE
	case 2:
		return inventory_v1.Category_CATEGORY_FUEL
	case 3:
		return inventory_v1.Category_CATEGORY_PORTHOLE
	case 4:
		return inventory_v1.Category_CATEGORY_WING
	default:
		return inventory_v1.Category_CATEGORY_UNSPECIFIED
	}
}

func CategoryToDomain(c inventory_v1.Category) model.Category {
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

func DimensionsToProto(d model.Dimensions) inventory_v1.Dimensions {
	return inventory_v1.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufactererToProto(m model.Manufacturer) inventory_v1.Manufacturer {
	return inventory_v1.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ValueToProto(val model.Value) *inventory_v1.Value {
	switch v := val.(type) {
	case model.StringValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_StringValue{
				StringValue: v.Value,
			},
		}
	case model.Int64Value:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_Int64Value{
				Int64Value: v.Value,
			},
		}
	case model.DoubleValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_DoubleValue{
				DoubleValue: v.Value,
			},
		}
	case model.BoolValue:
		return &inventory_v1.Value{
			Value: &inventory_v1.Value_BoolValue{
				BoolValue: v.Value,
			},
		}
	default:
		return nil
	}
}

func TimeToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func FilterToDomain(f *inventory_v1.PartsFilter) *model.PartsFilter {
	var categorysDomain []model.Category
	for _, c := range f.Categorys {
		categorysDomain = append(categorysDomain, CategoryToDomain(c))
	}

	return &model.PartsFilter{
		Uuids:     f.Uuids,
		Names:     f.Names,
		Categorys: categorysDomain,
		Countrys:  f.Countrys,
		Tags:      f.Tags,
	}
}
