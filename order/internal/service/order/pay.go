package order

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
)

func (s *service) PayOrder(ctx context.Context, orderId string, pm model.PaymentMethod) (transId string, err error) {
	order := s.orderRepo.GetOrder(orderId)
	if order == nil {
		return "", model.ErrNotFound
	}

	userId := order.UserId
	transId, err = s.paymentClient.PayOrder(ctx, userId, orderId, pm)
	if err != nil {
		return "", err
	}

	order.Status = model.OrderStatusPAID
	order.PaymentMethod = &pm
	order.TransactionID = &transId

	s.orderRepo.UpdateOrder(order)

	return
}
