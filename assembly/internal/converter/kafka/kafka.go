package kafka

import (
	"github.com/ArchibaldKronin/microservices_test/assembly/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
)

type OrderPaidDecoder interface {
	Decode([]byte) (model.OrderPaidEvent, error)
}

func OrderPaidEventToModel(rpcModel *events_v1.OrderPaid) model.OrderPaidEvent {
	return model.OrderPaidEvent{
		OrderUuid:       rpcModel.OrderUuid,
		UserUuid:        rpcModel.UserUuid,
		EventUuid:       rpcModel.EventUuid,
		PaymentMethod:   PaymentMethodToModel(rpcModel.PaymentMethod),
		TransactionUuid: rpcModel.TransactionUuid,
	}
}

func ModelToShipAssembledEvent(model model.ShipAssembledEvent) *events_v1.ShipAssembled {
	return &events_v1.ShipAssembled{
		EventUuid:    model.EventUuid,
		OrderUuid:    model.OrderUuid,
		UserUuid:     model.UserUuid,
		BuildTimeSec: model.BuildTimeSec,
	}
}

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
