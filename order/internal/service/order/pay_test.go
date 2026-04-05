package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderId       = gofakeit.UUID()
		userId        = gofakeit.UUID()
		partId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodSBP
		transactionId = gofakeit.UUID()

		order = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusPENDINGPAYMENT,
		}

		orderUpdate = &model.Order{
			OrderId:       orderId,
			UserId:        userId,
			PartIds:       []string{partId},
			Total_price:   42.2,
			Status:        model.OrderStatusPAID,
			TransactionID: &transactionId,
			PaymentMethod: &paymentMethod,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return(transactionId, nil).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, orderUpdate).Return(orderUpdate).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(transactionId, res)
}

func (s *ServiceSuite) TestPayErrNotFoundGet() {
	var (
		orderId = gofakeit.UUID()
		// userId        = gofakeit.UUID()
		// partId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodSBP
		// transactionId = gofakeit.UUID()

		// order = &model.Order{
		// 	OrderId:     orderId,
		// 	UserId:      userId,
		// 	PartIds:     []string{partId},
		// 	Total_price: 42.2,
		// 	Status:      model.OrderStatusPENDINGPAYMENT,
		// }

		// orderUpdate = &model.Order{
		// 	OrderId:       orderId,
		// 	UserId:        userId,
		// 	PartIds:       []string{partId},
		// 	Total_price:   42.2,
		// 	Status:        model.OrderStatusPAID,
		// 	TransactionID: &transactionId,
		// 	PaymentMethod: &paymentMethod,
		// }
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(nil).Once()
	// s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return(transactionId, nil).Once()
	// s.orderRepository.EXPECT().UpdateOrder(s.ctx, orderUpdate).Return(orderUpdate).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayErrNotFoundUpdate() {
	var (
		orderId       = gofakeit.UUID()
		userId        = gofakeit.UUID()
		partId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodSBP
		transactionId = gofakeit.UUID()

		order = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusPENDINGPAYMENT,
		}

		orderUpdate = &model.Order{
			OrderId:       orderId,
			UserId:        userId,
			PartIds:       []string{partId},
			Total_price:   42.2,
			Status:        model.OrderStatusPAID,
			TransactionID: &transactionId,
			PaymentMethod: &paymentMethod,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return(transactionId, nil).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, orderUpdate).Return(nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
	s.Empty(res)
}

func (s *ServiceSuite) TestPayErrInternal() {
	var (
		orderId       = gofakeit.UUID()
		userId        = gofakeit.UUID()
		partId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodSBP

		order = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return("", model.ErrInternal).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrInternal)
	s.Empty(res)
}
