package part

import (
	"fmt"
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
)

func (s *ServiceSuite) TestListParts() {
	var (
		parts = []*model.Part{
			{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Turbo Engine X1",
				Description:   "High performance aircraft engine",
				Price:         125000.50,
				StockQuantity: 5,
				Category:      model.CategoryFuel,
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
				UpdatedAt: nil,
			},
			{
				Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
				Name:          "Fuel Tank F200",
				Description:   "Composite aircraft fuel tank",
				Price:         32000,
				StockQuantity: 12,
				Category:      model.CategoryFuel,
				Dimensions: model.Dimensions{
					Length: 3.0,
					Width:  1.5,
					Height: 1.5,
					Weight: 200,
				},
				Manufacturer: model.Manufacturer{
					Name:    "FuelTech",
					Country: "Germany",
					Website: "https://fueltech.example.com",
				},
				Tags: []string{"fuel", "tank"},
				Metadata: map[string]model.Value{
					"capacity_liters": model.DoubleValue{Value: 1500.75},
					"has_sensor":      model.BoolValue{Value: true},
					"material":        model.StringValue{Value: "Composite"},
				},
				CreatedAt: time.Now(),
				UpdatedAt: nil,
			},
		}

		filter = &model.PartsFilter{
			Uuids: []string{"550e8400-e29b-41d4-a716-446655440000", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
		}
	)

	s.inventoryRepository.EXPECT().ListParts(s.ctx, filter).Return(parts, nil).Once()

	res, err := s.service.ListParts(s.ctx, filter)
	s.NoError(err)
	s.Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsErrGeneric() {
	var (
		parts []*model.Part

		filter = &model.PartsFilter{
			Uuids: []string{"edsfv23ref"},
		}
	)

	s.inventoryRepository.EXPECT().ListParts(s.ctx, filter).Return(parts, fmt.Errorf("generic error")).Once()

	res, err := s.service.ListParts(s.ctx, filter)
	s.Error(err)
	s.Equal(parts, res)
}

func (s *ServiceSuite) TestListPartsErrInvalidArgument() {
	var (
		parts []*model.Part

		filter *model.PartsFilter
	)

	res, err := s.service.ListParts(s.ctx, filter)
	s.Nil(res)
	s.Equal(parts, res)
	s.Error(err)
	s.ErrorIs(err, model.ErrInvalidArgument)
}
