package converter

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FilterToDTO(filter model.PartsFilter) inventory_v1.PartsFilter {
	var categorys []inventory_v1.Category
	for _, cat := range filter.Categorys {
		categorys = append(categorys, CategoryToDTO(cat))
	}
	return inventory_v1.PartsFilter{
		Uuids:     filter.Uuids,
		Names:     filter.Names,
		Categorys: categorys,
		Countrys:  filter.Countrys,
		Tags:      filter.Tags,
	}
}

func CategoryToDTO(cat model.Category) inventory_v1.Category {
	switch cat {
	case model.CategoryEngine:
		return inventory_v1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventory_v1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventory_v1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventory_v1.Category_CATEGORY_WING
	default:
		return inventory_v1.Category_CATEGORY_UNSPECIFIED
	}
}

func CategoryToDomain(cat inventory_v1.Category) model.Category {
	switch cat {
	case inventory_v1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventory_v1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventory_v1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventory_v1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnknown
	}
}

func DimensionsToDomain(d *inventory_v1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: d.Length,
		Width:  d.Width,
		Height: d.Height,
		Weight: d.Weight,
	}
}

func ManufactererToDomain(m *inventory_v1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    m.Name,
		Country: m.Country,
		Website: m.Website,
	}
}

func ValueToDomain(val *inventory_v1.Value) model.Value {
	if val == nil {
		return nil
	}

	switch v := val.Value.(type) {
	case *inventory_v1.Value_StringValue:
		return model.StringValue{
			Value: v.StringValue,
		}
	case *inventory_v1.Value_Int64Value:
		return model.Int64Value{
			Value: v.Int64Value,
		}
	case *inventory_v1.Value_DoubleValue:
		return model.DoubleValue{
			Value: v.DoubleValue,
		}
	case *inventory_v1.Value_BoolValue:
		return model.BoolValue{
			Value: v.BoolValue,
		}
	default:
		return nil
	}
}

func PaymentMethodToDTO(pm model.PaymentMethod) payment_v1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodCREDITCARD:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func MapError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return err
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return model.ErrInvalidArgument
	case codes.Internal:
		return model.ErrInternal
	case codes.NotFound:
		return model.ErrNotFound
	case codes.Unavailable:
		return model.ErrUnavailable
	default:
		return model.ErrInternal
	}
}
