package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestPaySuccess() {
	var (
		orderId       = gofakeit.UUID()
		userId        = gofakeit.UUID()
		partId        = gofakeit.UUID()
		paymentMethod = model.PaymentMethodSBP
		transactionId = gofakeit.UUID()

		order = &model.Order{
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		orderUpdate = &model.Order{
			OrderID:       orderId,
			UserID:        userId,
			PartIDs:       []string{partId},
			TotalPrice:    42.2,
			Status:        model.OrderStatusPAID,
			TransactionID: &transactionId,
			PaymentMethod: &paymentMethod,
		}
	)

	s.txManager.
		On("WithTransaction", s.ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.OrderRepository) error)

			err := fn(s.orderRepository)
			s.NoError(err)
		}).
		Return(nil).
		Once()
	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order, nil).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return(transactionId, nil).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, orderUpdate).Return(nil).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.NoError(err)
	s.NotNil(res)
	s.Equal(transactionId, res)
}

func (s *ServiceSuite) TestPayErrNotFoundGet() {
	var (
		orderId = gofakeit.UUID()

		paymentMethod = model.PaymentMethodSBP
	)

	s.txManager.
		On("WithTransaction", s.ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.OrderRepository) error)

			err := fn(s.orderRepository)
			s.Error(err)
			s.ErrorIs(err, model.ErrNotFound)
		}).
		Return(model.ErrNotFound).
		Once()
	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(nil, repoModel.ErrNotFound).Once()

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
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		orderUpdate = &model.Order{
			OrderID:       orderId,
			UserID:        userId,
			PartIDs:       []string{partId},
			TotalPrice:    42.2,
			Status:        model.OrderStatusPAID,
			TransactionID: &transactionId,
			PaymentMethod: &paymentMethod,
		}
	)

	s.txManager.
		On("WithTransaction", s.ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.OrderRepository) error)

			err := fn(s.orderRepository)
			s.Error(err)
			s.ErrorIs(err, model.ErrNotFound)
		}).
		Return(model.ErrNotFound).
		Once()
	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order, nil).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return(transactionId, nil).Once()
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, orderUpdate).Return(repoModel.ErrNotFound).Once()

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
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.txManager.
		On("WithTransaction", s.ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.OrderRepository) error)

			err := fn(s.orderRepository)
			s.Error(err)
			s.ErrorIs(err, model.ErrInternal)
		}).
		Return(model.ErrInternal).
		Once()
	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order, nil).Once()
	s.paymentClient.EXPECT().PayOrder(s.ctx, userId, orderId, paymentMethod).Return("", model.ErrInternal).Once()

	res, err := s.service.PayOrder(s.ctx, orderId, paymentMethod)

	s.Error(err)
	s.ErrorIs(err, model.ErrInternal)
	s.Empty(res)
}
