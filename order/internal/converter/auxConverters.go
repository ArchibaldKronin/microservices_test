package converter

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func convertStringsToUUIDs(in []string) ([]uuid.UUID, error) {
	uuids := make([]uuid.UUID, 0, len(in))
	for _, id := range in {
		if u, err := uuid.Parse(id); err == nil {
			uuids = append(uuids, u)
		} else {
			return nil, err
		}
	}
	return uuids, nil
}

func UUIDsToString(in []uuid.UUID) []string {
	ids := make([]string, 0, len(in))
	for _, u := range in {
		ids = append(ids, u.String())
	}
	return ids
}

func PaymentMethodToDTO(pm model.PaymentMethod) order_v1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCARD:
		return order_v1.PaymentMethodCARD
	case model.PaymentMethodCREDITCARD:
		return order_v1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return order_v1.PaymentMethodINVESTORMONEY
	case model.PaymentMethodSBP:
		return order_v1.PaymentMethodSBP
	default:
		return order_v1.PaymentMethodUNKNOWN
	}
}

func PaymentMethodToDomain(pm order_v1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case order_v1.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case order_v1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case order_v1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	case order_v1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func OrderStatusToDTO(os model.OrderStatus) order_v1.OrderStatus {
	switch os {
	case model.OrderStatusCANCELLED:
		return order_v1.OrderStatusCANCELLED
	case model.OrderStatusPAID:
		return order_v1.OrderStatusPAID
	case model.OrderStatusCOMPLETED:
		return order_v1.OrderStatusCOMPLETED
	default:
		return order_v1.OrderStatusPENDINGPAYMENT
	}
}
