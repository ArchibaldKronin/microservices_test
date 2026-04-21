package part

// import (
// 	"context"
// 	"testing"
// 	"time"

// 	"github.com/ArchibaldKronin/microservices_test/inventory/internal/model"
// 	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
// 	"github.com/stretchr/testify/require"
// )

// func TestNewRepo(t *testing.T) {
// 	expected := map[string]repoModel.Part{
// 		"550e8400-e29b-41d4-a716-446655440000": {
// 			Uuid:          "550e8400-e29b-41d4-a716-446655440000",
// 			Name:          "Turbo Engine X1",
// 			Description:   "High performance aircraft engine",
// 			Price:         125000.50,
// 			StockQuantity: 5,
// 			Category:      repoModel.CategoryEngine,
// 			Dimensions: repoModel.Dimensions{
// 				Length: 2.5,
// 				Width:  1.2,
// 				Height: 1.5,
// 				Weight: 850,
// 			},
// 			Manufacturer: repoModel.Manufacturer{
// 				Name:    "JetCorp",
// 				Country: "USA",
// 				Website: "https://jetcorp.example.com",
// 			},
// 			Tags: []string{"engine", "turbo"},
// 			Metadata: map[string]repoModel.Value{
// 				"horsepower":   repoModel.Int64Value{Value: 4500},
// 				"fuel_type":    repoModel.StringValue{Value: "Jet A-1"},
// 				"is_certified": repoModel.BoolValue{Value: true},
// 			},
// 			CreatedAt: time.Time{},
// 			UpdatedAt: time.Time{},
// 		},
// 	}

// 	res := NewRepository(InitialParts)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Uuid, res.data["550e8400-e29b-41d4-a716-446655440000"].Uuid)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Name, res.data["550e8400-e29b-41d4-a716-446655440000"].Name)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Metadata, res.data["550e8400-e29b-41d4-a716-446655440000"].Metadata)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Description, res.data["550e8400-e29b-41d4-a716-446655440000"].Description)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Price, res.data["550e8400-e29b-41d4-a716-446655440000"].Price)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].StockQuantity, res.data["550e8400-e29b-41d4-a716-446655440000"].StockQuantity)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Category, res.data["550e8400-e29b-41d4-a716-446655440000"].Category)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Dimensions, res.data["550e8400-e29b-41d4-a716-446655440000"].Dimensions)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Manufacturer, res.data["550e8400-e29b-41d4-a716-446655440000"].Manufacturer)
// 	require.Equal(t, expected["550e8400-e29b-41d4-a716-446655440000"].Tags, res.data["550e8400-e29b-41d4-a716-446655440000"].Tags)
// }

// func TestGetSuccess(t *testing.T) {
// 	repo := NewRepository(InitialParts)
// 	ctx := context.Background()
// 	id := "550e8400-e29b-41d4-a716-446655440000"

// 	expected := &model.Part{
// 		Uuid:          "550e8400-e29b-41d4-a716-446655440000",
// 		Name:          "Turbo Engine X1",
// 		Description:   "High performance aircraft engine",
// 		Price:         125000.50,
// 		StockQuantity: 5,
// 		Category:      model.CategoryEngine,
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
// 		CreatedAt: time.Time{},
// 		UpdatedAt: time.Time{},
// 	}

// 	res, err := repo.GetPart(ctx, id)
// 	res.CreatedAt = time.Time{}
// 	res.UpdatedAt = time.Time{}

// 	require.NoError(t, err)
// 	require.NotNil(t, res)
// 	require.Equal(t, expected, res)
// }

// func TestGetErrNotFound(t *testing.T) {
// 	repo := NewRepository(InitialParts)
// 	ctx := context.Background()
// 	id := "550e8400-e29b-41d4-a716-446655440001"

// 	res, err := repo.GetPart(ctx, id)

// 	require.Nil(t, res)
// 	require.Error(t, err)
// 	require.ErrorIs(t, err, repoModel.ErrNotFound)
// }
