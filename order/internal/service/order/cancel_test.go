package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderId = gofakeit.UUID()
		userId  = gofakeit.UUID()
		partId  = gofakeit.UUID()
		order   = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusPENDINGPAYMENT,
		}

		cancelledOrder = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusCANCELLED,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, cancelledOrder).Return(cancelledOrder).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelErrNotFoundGet() {
	orderId := gofakeit.UUID()

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
}

func (s *ServiceSuite) TestCancelErrNotFoundUpdate() {
	var (
		orderId = gofakeit.UUID()
		userId  = gofakeit.UUID()
		partId  = gofakeit.UUID()
		order   = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusPENDINGPAYMENT,
		}

		cancelledOrder = &model.Order{
			OrderId:     orderId,
			UserId:      userId,
			PartIds:     []string{partId},
			Total_price: 42.2,
			Status:      model.OrderStatusCANCELLED,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, cancelledOrder).Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
}

func (s *ServiceSuite) TestCancelErrOrderAlreadyPaid() {
	var (
		orderId = gofakeit.UUID()
		order   = &model.Order{
			OrderId:     orderId,
			UserId:      gofakeit.UUID(),
			PartIds:     []string{gofakeit.UUID()},
			Total_price: 42.2,
			Status:      model.OrderStatusPAID,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}
