package converter

import (
	"github.com/ArchibaldKronin/microservices_test/notification/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
)

func PaymentMethodToModel(pm events_v1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case events_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCARD
	case events_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCREDITCARD
	case events_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodINVESTORMONEY
	case events_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	default:
		return model.PaymentMethodUNKNOWN
	}
}
