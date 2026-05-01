package order

import (
	"time"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
)

func (s *ServiceSuite) TestGetPartsInfoSuccess() {
	var (
		partId  = gofakeit.UUID()
		partsId = []string{partId}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		updatedAt = time.Now()
		createdAt = time.Now().Add(-10 * time.Minute)

		part = model.Part{
			Uuid:          partId,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
			Category:      model.CategoryEngine,
			Price:         125000.50,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: model.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: []string{"engine", "turbo"},
			Metadata: map[string]model.Value{
				"horsepower":   model.Int64Value{Value: 4500},
				"fuel_type":    model.StringValue{Value: "Jet A-1"},
				"is_certified": model.BoolValue{Value: true},
			},
			CreatedAt: createdAt,
			UpdatedAt: lo.ToPtr(updatedAt),
		}
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]*model.Part{&part}, nil).Once()

	ids, price, err := s.service.getPartsInfo(s.ctx, partsId)

	s.NoError(err)
	s.Equal(partId, ids[0])
	s.Equal(part.Price, price)
}

func (s *ServiceSuite) TestGetPartsErrInvalidArgumentListParts() {
	logger.SetNopLogger()

	var (
		partId  = gofakeit.UUID()
		partsId = []string{partId}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		expectedError = model.ErrInvalidArgument
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return(nil, expectedError).Once()

	ids, price, err := s.service.getPartsInfo(s.ctx, partsId)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.Empty(ids)
	s.Empty(price)
}

func (s *ServiceSuite) TestGetPartsErrInternalListParts() {
	logger.SetNopLogger()

	var (
		partId  = gofakeit.UUID()
		partsId = []string{partId}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		expectedError = model.ErrInternal
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return(nil, expectedError).Once()

	ids, price, err := s.service.getPartsInfo(s.ctx, partsId)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.Empty(ids)
	s.Empty(price)
}

func (s *ServiceSuite) TestGetPartsErrNotFound() {
	logger.SetNopLogger()

	var (
		partId1 = gofakeit.UUID()
		partId2 = gofakeit.UUID()
		partsId = []string{partId1, partId2}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		updatedAt = time.Now()
		createdAt = time.Now().Add(-10 * time.Minute)

		part1 = model.Part{
			Uuid:          partId1,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
			Category:      model.CategoryEngine,
			Price:         125000.50,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: model.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: []string{"engine", "turbo"},
			Metadata: map[string]model.Value{
				"horsepower":   model.Int64Value{Value: 4500},
				"fuel_type":    model.StringValue{Value: "Jet A-1"},
				"is_certified": model.BoolValue{Value: true},
			},
			CreatedAt: createdAt,
			UpdatedAt: lo.ToPtr(updatedAt),
		}

		expectedError = model.ErrNotFound
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]*model.Part{&part1}, nil).Once()

	ids, price, err := s.service.getPartsInfo(s.ctx, partsId)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.Empty(ids)
	s.Empty(price)
}

func (s *ServiceSuite) TestCreateSuccess() {
	logger.SetNopLogger()

	var (
		partId  = gofakeit.UUID()
		partsId = []string{partId}
		userId  = gofakeit.UUID()

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		updatedAt = time.Now()
		createdAt = time.Now().Add(-10 * time.Minute)

		part = model.Part{
			Uuid:          partId,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
			Category:      model.CategoryEngine,
			Price:         125000.50,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: model.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: []string{"engine", "turbo"},
			Metadata: map[string]model.Value{
				"horsepower":   model.Int64Value{Value: 4500},
				"fuel_type":    model.StringValue{Value: "Jet A-1"},
				"is_certified": model.BoolValue{Value: true},
			},
			CreatedAt: createdAt,
			UpdatedAt: lo.ToPtr(updatedAt),
		}

		order = model.NewOrder(userId, partsId, part.Price)
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]*model.Part{&part}, nil).Once()
	s.orderRepository.EXPECT().CreateOrder(s.ctx, mock.AnythingOfType("*model.Order")).Return(nil).Once()

	res, err := s.service.CreateOrder(s.ctx, userId, partsId)

	s.NoError(err)
	s.Equal(order.PartIDs, res.PartIDs)
	s.Equal(order.PaymentMethod, res.PaymentMethod)
	s.Equal(order.Status, res.Status)
	s.Equal(order.PartIDs, res.PartIDs)
	s.Equal(order.UserID, res.UserID)
}

func (s *ServiceSuite) TestCreateErrInvalidArgument() {
	logger.SetNopLogger()

	var (
		userId = gofakeit.UUID()

		partId  = gofakeit.UUID()
		partsId = []string{partId}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		expectedError = model.ErrInvalidArgument
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return(nil, expectedError).Once()

	res, err := s.service.CreateOrder(s.ctx, userId, partsId)

	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidArgument)
	s.Nil(res)
}

func (s *ServiceSuite) TestCreateErrInternal() {
	logger.SetNopLogger()

	var (
		userId = gofakeit.UUID()

		partId  = gofakeit.UUID()
		partsId = []string{partId}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		expectedError = model.ErrInternal
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return(nil, expectedError).Once()

	res, err := s.service.CreateOrder(s.ctx, userId, partsId)

	s.Error(err)
	s.ErrorIs(err, model.ErrInternal)
	s.Nil(res)
}

func (s *ServiceSuite) TestCreateErrNotFound() {
	logger.SetNopLogger()

	var (
		userId = gofakeit.UUID()

		partId1 = gofakeit.UUID()
		partId2 = gofakeit.UUID()
		partsId = []string{partId1, partId2}

		filter = model.PartsFilter{
			Uuids: partsId,
		}

		updatedAt = (*time.Time)(nil)
		createdAt = time.Now().Add(-10 * time.Minute)

		part1 = model.Part{
			Uuid:          partId1,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
			Category:      model.CategoryEngine,
			Price:         125000.50,
			StockQuantity: 5,
			Dimensions: model.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: model.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: []string{"engine", "turbo"},
			Metadata: map[string]model.Value{
				"horsepower":   model.Int64Value{Value: 4500},
				"fuel_type":    model.StringValue{Value: "Jet A-1"},
				"is_certified": model.BoolValue{Value: true},
			},
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		expectedError = model.ErrNotFound
	)

	s.inventoryClient.EXPECT().ListParts(s.ctx, filter).Return([]*model.Part{&part1}, nil).Once()

	res, err := s.service.CreateOrder(s.ctx, userId, partsId)

	s.Error(err)
	s.ErrorIs(err, expectedError)
	s.Nil(res)
}
