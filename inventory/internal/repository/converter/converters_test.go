package converter

import (
	"testing"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"

	// "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
)

func TestConverters(t *testing.T) {
	// t.Run("Part to domain", func(t *testing.T) {
	// 	uuid := gofakeit.UUID()
	// 	CreatedAt := time.Now()
	// 	UpdatedAt := time.Now()

	// 	part := repoModel.Part{
	// 		Uuid:          uuid,
	// 		Name:          "Turbo Engine X1",
	// 		Description:   "High performance aircraft engine",
	// 		Category:      repoModel.CategoryEngine,
	// 		Price:         125000.50,
	// 		StockQuantity: 5,
	// 		Dimensions: repoModel.Dimensions{
	// 			Length: 2.5,
	// 			Width:  1.2,
	// 			Height: 1.5,
	// 			Weight: 850,
	// 		},
	// 		Manufacturer: repoModel.Manufacturer{
	// 			Name:    "JetCorp",
	// 			Country: "USA",
	// 			Website: "https://jetcorp.example.com",
	// 		},
	// 		Tags: []string{"engine", "turbo"},
	// 		Metadata: map[string]repoModel.Value{
	// 			"horsepower":   4500,
	// 			"fuel_type":    "Jet A-1",
	// 			"is_certified": true,
	// 		},
	// 		CreatedAt: CreatedAt,
	// 		UpdatedAt: UpdatedAt,
	// 	}

	// 	expected := model.Part{
	// 		Uuid:          uuid,
	// 		Name:          "Turbo Engine X1",
	// 		Description:   "High performance aircraft engine",
	// 		Category:      model.CategoryEngine,
	// 		Price:         125000.50,
	// 		StockQuantity: 5,
	// 		Dimensions: model.Dimensions{
	// 			Length: 2.5,
	// 			Width:  1.2,
	// 			Height: 1.5,
	// 			Weight: 850,
	// 		},
	// 		Manufacturer: model.Manufacturer{
	// 			Name:    "JetCorp",
	// 			Country: "USA",
	// 			Website: "https://jetcorp.example.com",
	// 		},
	// 		Tags: []string{"engine", "turbo"},
	// 		Metadata: map[string]model.Value{
	// 			"horsepower":   model.Int64Value{Value: 4500},
	// 			"fuel_type":    model.StringValue{Value: "Jet A-1"},
	// 			"is_certified": model.BoolValue{Value: true},
	// 		},
	// 		CreatedAt: CreatedAt,
	// 		UpdatedAt: UpdatedAt,
	// 	}

	// 	res := PartToDomain(&part)
	// 	require.NotNil(t, res)
	// 	require.Equal(t, &expected, res)
	// })
	t.Run("Value to primitive", func(t *testing.T) {
		tests := []struct {
			val      model.Value
			expected any
		}{
			{
				val:      model.StringValue{Value: "test_string"},
				expected: "test_string",
			},
			{
				val:      model.Int64Value{Value: 42},
				expected: int64(42),
			},
			{
				val:      model.DoubleValue{Value: 36.6},
				expected: 36.6,
			},
			{
				val:      model.BoolValue{Value: true},
				expected: true,
			},
		}

		for _, test := range tests {
			res, err := ValueToPrimitive(test.val)
			require.NoError(t, err)
			require.Equal(t, test.expected, res)
		}
	})
	t.Run("Primitive to value", func(t *testing.T) {
		tests := []struct {
			val      any
			expected model.Value
		}{
			{
				val:      "test_string",
				expected: model.StringValue{Value: "test_string"},
			},
			{
				val:      int(42),
				expected: model.Int64Value{Value: 42},
			},
			{
				val:      36.6,
				expected: model.DoubleValue{Value: 36.6},
			},
			{
				val:      true,
				expected: model.BoolValue{Value: true},
			},
		}

		for _, test := range tests {
			res, err := PrimitiveToValue(test.val)
			require.NoError(t, err)
			require.Equal(t, test.expected, res)
		}

		val := model.Part{}
		res, err := PrimitiveToValue(val)
		require.Nil(t, res)
		require.Error(t, err)
		require.ErrorIs(t, err, repoModel.ErrParseValue)
		var targetErr *repoModel.MetadataParseValueError
		require.ErrorAs(t, err, &targetErr)
		require.Equal(t, val, targetErr.Value)
	})
	t.Run("Category to domain", func(t *testing.T) {
		tests := []struct {
			val      repoModel.Category
			expected model.Category
		}{
			{
				val:      repoModel.CategoryEngine,
				expected: model.CategoryEngine,
			},
			{
				val:      repoModel.CategoryFuel,
				expected: model.CategoryFuel,
			},
			{
				val:      repoModel.CategoryPorthole,
				expected: model.CategoryPorthole,
			},
			{
				val:      repoModel.CategoryWing,
				expected: model.CategoryWing,
			},
			{
				val:      repoModel.CategoryUnknown,
				expected: model.CategoryUnknown,
			},
		}

		for _, test := range tests {
			res := CategoryToDomain(test.val)
			require.Equal(t, test.expected, res)
		}
	})
	t.Run("Category to repo", func(t *testing.T) {
		tests := []struct {
			val      model.Category
			expected repoModel.Category
		}{
			{
				val:      model.CategoryEngine,
				expected: repoModel.CategoryEngine,
			},
			{
				val:      model.CategoryFuel,
				expected: repoModel.CategoryFuel,
			},
			{
				val:      model.CategoryPorthole,
				expected: repoModel.CategoryPorthole,
			},
			{
				val:      model.CategoryWing,
				expected: repoModel.CategoryWing,
			},
			{
				val:      model.CategoryUnknown,
				expected: repoModel.CategoryUnknown,
			},
		}

		for _, test := range tests {
			res := CategoryToRepo(test.val)
			require.Equal(t, test.expected, res)
		}
	})
	t.Run("Dimensions to domain", func(t *testing.T) {
		dimensions := repoModel.Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.5,
			Weight: 850,
		}

		expected := model.Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.5,
			Weight: 850,
		}

		res := DimensionsToDomain(dimensions)
		require.Equal(t, expected, res)
	})
	t.Run("Manufacterer to domain", func(t *testing.T) {
		manufacturer := repoModel.Manufacturer{
			Name:    "JetCorp",
			Country: "USA",
			Website: "https://jetcorp.example.com",
		}

		expected := model.Manufacturer{
			Name:    "JetCorp",
			Country: "USA",
			Website: "https://jetcorp.example.com",
		}

		res := ManufactererToDomain(manufacturer)
		require.Equal(t, expected, res)
	})
}
