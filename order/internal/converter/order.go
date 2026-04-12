package converter

import (
	"fmt"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	order_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
	"github.com/google/uuid"
)

func OrderToDTO(o *model.Order) (order_v1.OrderDto, error) {
	orderIdDTO, err := uuid.Parse(o.OrderID)
	if err != nil {
		return order_v1.OrderDto{}, fmt.Errorf("error parse orderId: %w", err)
	}
	userIdDTO, err := uuid.Parse(o.UserID)
	if err != nil {
		return order_v1.OrderDto{}, fmt.Errorf("error parse orderId: %w", err)
	}
	partsIdsDTO, err := convertStringsToUUIDs(o.PartIDs)
	if err != nil {
		return order_v1.OrderDto{}, fmt.Errorf("error parse strings to UUIDs: %w", err)
	}

	var optTransactionId order_v1.OptUUID
	if o.TransactionID != nil {
		txId, err := uuid.Parse(*o.TransactionID)
		if err != nil {
			return order_v1.OrderDto{}, fmt.Errorf("error parse transactionID: %w", err)
		}
		optTransactionId = order_v1.OptUUID{
			Value: txId,
			Set:   true,
		}
	}

	var optPaymentMethod order_v1.OptPaymentMethod
	if o.PaymentMethod != nil {
		optPaymentMethod = order_v1.OptPaymentMethod{
			Value: PaymentMethodToDTO(*o.PaymentMethod),
			Set:   true,
		}
	}

	return order_v1.OrderDto{
		OrderUUID:       orderIdDTO,
		UserUUID:        userIdDTO,
		PartUuids:       partsIdsDTO,
		TotalPrice:      float32(o.TotalPrice),
		TransactionUUID: optTransactionId,
		PaymentMethod:   optPaymentMethod,
		Status:          OrderStatusToDTO(o.Status),
	}, nil
}
