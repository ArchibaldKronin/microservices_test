package converter

import (
	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
)

func PaymentMethodToDomain(pm model.PaymentMethod) serviceModel.PaymentMethod {
	switch pm {
	case model.PaymentMethodCARD:
		return serviceModel.PaymentMethodCARD
	case model.PaymentMethodCREDITCARD:
		return serviceModel.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return serviceModel.PaymentMethodINVESTORMONEY
	case model.PaymentMethodSBP:
		return serviceModel.PaymentMethodSBP
	default:
		return serviceModel.PaymentMethodUNKNOWN
	}
}

func PaymentMethodToRepo(pm serviceModel.PaymentMethod) model.PaymentMethod {
	switch pm {
	case serviceModel.PaymentMethodCARD:
		return model.PaymentMethodCARD
	case serviceModel.PaymentMethodCREDITCARD:
		return model.PaymentMethodCREDITCARD
	case serviceModel.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodINVESTORMONEY
	case serviceModel.PaymentMethodSBP:
		return model.PaymentMethodSBP
	default:
		return model.PaymentMethodUNKNOWN
	}
}

func StatusToDomain(s model.OrderStatus) serviceModel.OrderStatus {
	switch s {
	case model.OrderStatusPENDINGPAYMENT:
		return serviceModel.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return serviceModel.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return serviceModel.OrderStatusCANCELLED
	case model.OrderStatusCOMPLETED:
		return serviceModel.OrderStatusCOMPLETED
	default:
		return serviceModel.OrderStatusCANCELLED
	}
}

func StatusToRepo(s serviceModel.OrderStatus) model.OrderStatus {
	switch s {
	case serviceModel.OrderStatusPENDINGPAYMENT:
		return model.OrderStatusPENDINGPAYMENT
	case serviceModel.OrderStatusPAID:
		return model.OrderStatusPAID
	case serviceModel.OrderStatusCANCELLED:
		return model.OrderStatusCANCELLED
	case serviceModel.OrderStatusCOMPLETED:
		return model.OrderStatusCOMPLETED
	default:
		return model.OrderStatusCANCELLED
	}
}
