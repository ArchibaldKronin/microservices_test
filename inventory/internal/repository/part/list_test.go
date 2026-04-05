package part

import (
	"context"
	"testing"
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	"github.com/stretchr/testify/require"
)

func TestListSuccess(t *testing.T) {
	var (
		ctx = context.Background()

		expected = []*model.Part{
			{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Turbo Engine X1",
				Description:   "High performance aircraft engine",
				Price:         125000.50,
				StockQuantity: 5,
				Category:      model.CategoryEngine,
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
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
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
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
		}

		filters = []*model.PartsFilter{
			{
				Uuids:     []string{"550e8400-e29b-41d4-a716-446655440000", "6ba7b810-9dad-11d1-80b4-00c04fd430c8"},
				Names:     []string{"Turbo Engine X1", "Fuel Tank F200"},
				Categorys: []model.Category{model.CategoryEngine, model.CategoryFuel},
				Countrys:  []string{"USA", "Germany"},
				Tags:      []string{"engine", "tank"},
			},
			{
				Names:     []string{"Turbo Engine X1", "Fuel Tank F200"},
				Categorys: []model.Category{model.CategoryEngine, model.CategoryFuel},
				Countrys:  []string{"USA", "Germany"},
				Tags:      []string{"engine", "tank"},
			},
			{
				Categorys: []model.Category{model.CategoryEngine, model.CategoryFuel},
				Countrys:  []string{"USA", "Germany"},
				Tags:      []string{"engine", "tank"},
			},
			{
				Countrys: []string{"USA", "Germany"},
				Tags:     []string{"engine", "tank"},
			},
			{
				Tags: []string{"engine", "tank"},
			},
		}
	)

	repo := NewRepository(InitialParts)

	for _, f := range filters {
		res, err := repo.ListParts(ctx, f)
		for _, part := range res {
			part.CreatedAt = time.Time{}
			part.UpdatedAt = time.Time{}
		}

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, expected, res)
	}
}

func TestListNil(t *testing.T) {
	ctx := context.Background()
	filter := &model.PartsFilter{
		Uuids: []string{"550e8400-e29b-41d4-a716-44665544001"},
	}

	repo := NewRepository(InitialParts)

	res, err := repo.ListParts(ctx, filter)

	require.NoError(t, err)
	require.Nil(t, res)
}
