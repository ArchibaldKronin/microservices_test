package converter

import (
	"github.com/ArchibaldKronin/microservices_test/payment/internal/model"
	payment_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/payment/v1"
)

func PaymentMethodToDomain(pm payment_v1.PaymentMethod) model.PaymentMethod {
	switch pm {
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case payment_v1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func PaymentMethodToProto(pm model.PaymentMethod) payment_v1.PaymentMethod {
	switch pm {
	case model.PaymentMethodCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CARD
	case model.PaymentMethodSBP:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_SBP
	case model.PaymentMethodCreditCard:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return payment_v1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}
