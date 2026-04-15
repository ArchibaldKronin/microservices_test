package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestCancelSuccess() {
	var (
		orderId = gofakeit.UUID()
		userId  = gofakeit.UUID()
		partId  = gofakeit.UUID()
		order   = &model.Order{
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		cancelledOrder = &model.Order{
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusCANCELLED,
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
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, cancelledOrder).Return(nil).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.NoError(err)
}

func (s *ServiceSuite) TestCancelErrNotFoundGet() {
	orderId := gofakeit.UUID()

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
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}

		cancelledOrder = &model.Order{
			OrderID:    orderId,
			UserID:     userId,
			PartIDs:    []string{partId},
			TotalPrice: 42.2,
			Status:     model.OrderStatusCANCELLED,
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
	s.orderRepository.EXPECT().UpdateOrder(s.ctx, cancelledOrder).Return(repoModel.ErrNotFound).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, model.ErrNotFound)
}

func (s *ServiceSuite) TestCancelErrOrderAlreadyPaid() {
	var (
		orderId = gofakeit.UUID()
		order   = &model.Order{
			OrderID:    orderId,
			UserID:     gofakeit.UUID(),
			PartIDs:    []string{gofakeit.UUID()},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPAID,
		}
	)

	s.txManager.
		On("WithTransaction", s.ctx, mock.Anything).
		Run(func(args mock.Arguments) {
			fn := args.Get(1).(func(repository.OrderRepository) error)

			err := fn(s.orderRepository)
			s.Error(err)
			s.ErrorIs(err, model.ErrOrderAlreadyPaid)
		}).
		Return(model.ErrOrderAlreadyPaid).
		Once()

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(order, nil).Once()

	err := s.service.CancelOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, model.ErrOrderAlreadyPaid)
}
