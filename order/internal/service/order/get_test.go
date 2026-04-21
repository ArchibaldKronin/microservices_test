package order

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		orderId  = gofakeit.UUID()
		expected = &model.Order{
			OrderID:    orderId,
			UserID:     gofakeit.UUID(),
			PartIDs:    []string{gofakeit.UUID()},
			TotalPrice: 42.2,
			Status:     model.OrderStatusPENDINGPAYMENT,
		}
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(expected, nil).Once()

	res, err := s.service.GetOrder(s.ctx, orderId)

	s.NoError(err)
	s.Equal(expected, res)
}

func (s *ServiceSuite) TestGetErrNotFound() {
	var (
		orderId = gofakeit.UUID()

		expected = model.ErrNotFound
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(nil, repoModel.ErrNotFound).Once()

	res, err := s.service.GetOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, expected)
	s.Nil(res)
}

func (s *ServiceSuite) TestGetErrInternal() {
	var (
		orderId = gofakeit.UUID()

		expected = model.ErrInternal
	)

	s.orderRepository.EXPECT().GetOrder(s.ctx, orderId).Return(nil, repoModel.ErrSelectQuery).Once()

	res, err := s.service.GetOrder(s.ctx, orderId)

	s.Error(err)
	s.ErrorIs(err, expected)
	s.Nil(res)
}
