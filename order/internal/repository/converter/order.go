package converter

import (
	"slices"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/samber/lo"
)

func OrderToDomain(o model.Order) serviceModel.Order {
	var pmService *serviceModel.PaymentMethod
	if o.PaymentMethod != nil {
		pmService = lo.ToPtr(PaymentMethodToDomain(*o.PaymentMethod))
	}

	return serviceModel.Order{
		OrderID:       o.OrderID,
		UserID:        o.UserID,
		PartIDs:       slices.Clone(o.PartIDs),
		TotalPrice:    o.TotalPrice,
		TransactionID: o.TransactionID,
		PaymentMethod: pmService,
		Status:        StatusToDomain(o.Status),
	}
}

func OrderToRepo(o serviceModel.Order) model.Order {
	var pmRepo *model.PaymentMethod
	if o.PaymentMethod != nil {
		pmRepo = lo.ToPtr(PaymentMethodToRepo(*o.PaymentMethod))
	}

	return model.Order{
		OrderID:       o.OrderID,
		UserID:        o.UserID,
		PartIDs:       slices.Clone(o.PartIDs),
		TotalPrice:    o.TotalPrice,
		TransactionID: o.TransactionID,
		PaymentMethod: pmRepo,
		Status:        StatusToRepo(o.Status),
	}
}
