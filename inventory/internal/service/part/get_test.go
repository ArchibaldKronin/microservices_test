package part

import (
	"fmt"
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestGetSuccess() {
	var (
		uuid = gofakeit.UUID()

		part = model.Part{
			Uuid:          uuid,
			Name:          "Turbo Engine X1",
			Description:   "High performance aircraft engine",
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	)

	s.inventoryRepository.EXPECT().GetPart(s.ctx, uuid).Return(&part, nil)

	res, err := s.service.GetPart(s.ctx, uuid)
	s.NoError(err)
	s.Equal(part, *res)
}

func (s *ServiceSuite) TestGetErrNotFound() {
	var (
		uuid = gofakeit.UUID()
		err  = repoModel.ErrNotFound
	)

	s.inventoryRepository.EXPECT().GetPart(s.ctx, uuid).Return(nil, err)

	res, resErr := s.service.GetPart(s.ctx, uuid)
	s.ErrorIs(model.ErrNotFound, resErr)
	s.Equal(res, (*model.Part)(nil))
}

func (s *ServiceSuite) TestGetErrGeneric() {
	uuid := gofakeit.UUID()

	s.inventoryRepository.EXPECT().GetPart(s.ctx, uuid).Return(nil, fmt.Errorf("generic error"))

	res, err := s.service.GetPart(s.ctx, uuid)
	s.Error(err)
	s.Equal(res, (*model.Part)(nil))
}
