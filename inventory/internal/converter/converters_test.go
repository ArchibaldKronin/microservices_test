package converter

import (
	"testing"
	"time"

	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
	inventory_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/inventory/v1"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestConverters(t *testing.T) {
	t.Run("PartToProto", func(t *testing.T) {
		uuid := gofakeit.UUID()

		part := model.Part{
			Uuid:          uuid,
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expected := inventory_v1.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: int64(part.StockQuantity),
			Category:      inventory_v1.Category_CATEGORY_ENGINE,
			Dimensions: &inventory_v1.Dimensions{
				Length: 2.5,
				Width:  1.2,
				Height: 1.5,
				Weight: 850,
			},
			Manufacturer: &inventory_v1.Manufacturer{
				Name:    "JetCorp",
				Country: "USA",
				Website: "https://jetcorp.example.com",
			},
			Tags: part.Tags,
			Metadata: map[string]*inventory_v1.Value{
				"horsepower":   {Value: &inventory_v1.Value_Int64Value{Int64Value: 4500}},
				"fuel_type":    {Value: &inventory_v1.Value_StringValue{StringValue: "Jet A-1"}},
				"is_certified": {Value: &inventory_v1.Value_BoolValue{BoolValue: true}},
			},
			CreatedAt: timestamppb.New(part.CreatedAt),
			UpdatedAt: timestamppb.New(part.UpdatedAt),
		}

		res := PartToProto(&part)
		require.NotNil(t, res)
		require.Equal(t, &expected, res)
	})

	t.Run("ValueToProto", func(t *testing.T) {
		tests := []struct {
			val      model.Value
			expected *inventory_v1.Value
		}{
			{
				val:      model.StringValue{Value: "test_string"},
				expected: &inventory_v1.Value{Value: &inventory_v1.Value_StringValue{StringValue: "test_string"}},
			},
			{
				val:      model.Int64Value{Value: 42},
				expected: &inventory_v1.Value{Value: &inventory_v1.Value_Int64Value{Int64Value: 42}},
			},
			{
				val:      model.DoubleValue{Value: 36.6},
				expected: &inventory_v1.Value{Value: &inventory_v1.Value_DoubleValue{DoubleValue: 36.6}},
			},
			{
				val:      model.BoolValue{Value: true},
				expected: &inventory_v1.Value{Value: &inventory_v1.Value_BoolValue{BoolValue: true}},
			},
		}

		for i := range tests {
			test := &tests[i]
			res := ValueToProto(test.val)
			require.Equal(t, test.expected, res)
		}
	})
}
