package kafka

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	events_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/events/v1"
)

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}

func ModelToOrderPaidEvent(model model.OrderPaidEvent) *events_v1.OrderPaid {
	return &events_v1.OrderPaid{
		EventUuid:       model.EventUuid,
		OrderUuid:       model.OrderUuid,
		UserUuid:        model.UserUuid,
		PaymentMethod:   PaymentMethodToEvent(model.PaymentMethod),
		TransactionUuid: model.TransactionUuid,
	}
}

func PaymentMethodToEvent(pm model.PaymentMethod) events_v1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCARD:
		return events_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodCREDITCARD:
		return events_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodINVESTORMONEY:
		return events_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	case model.PaymentMethodSBP:
		return events_v1.PaymentMethod_PAYMENT_METHOD_SBP
	default:
		return events_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func ShipAssembledEventToModel(pb *events_v1.ShipAssembled) model.ShipAssembledEvent {
	return model.ShipAssembledEvent{
		EventUuid:    pb.EventUuid,
		OrderUuid:    pb.OrderUuid,
		UserUuid:     pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}
}
