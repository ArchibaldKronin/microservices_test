package part

import (
	"time"

	repoModel "github.com/ArchibaldKronin/microservices_test/inventory/internal/repository/model"
	"github.com/samber/lo"
)

var InitialParts = []*repoModel.Part{
	{
		Uuid:          "550e8400-e29b-41d4-a716-446655440000",
		Name:          "Turbo Engine X1",
		Description:   "High performance aircraft engine",
		Price:         125000.50,
		StockQuantity: 5,
		Category:      repoModel.CategoryEngine,
		Dimensions: repoModel.Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.5,
			Weight: 850,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "JetCorp",
			Country: "USA",
			Website: "https://jetcorp.example.com",
		},
		Tags: []string{"engine", "turbo"},
		Metadata: map[string]any{
			"horsepower":   4500,
			"fuel_type":    "Jet A-1",
			"is_certified": true,
		},
		CreatedAt: time.Now().Add(-10 * time.Minute),
		UpdatedAt: lo.ToPtr(time.Now()),
	},
	{
		Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:          "Fuel Tank F200",
		Description:   "Composite aircraft fuel tank",
		Price:         32000,
		StockQuantity: 12,
		Category:      repoModel.CategoryFuel,
		Dimensions: repoModel.Dimensions{
			Length: 3.0,
			Width:  1.5,
			Height: 1.5,
			Weight: 200,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "FuelTech",
			Country: "Germany",
			Website: "https://fueltech.example.com",
		},
		Tags: []string{"fuel", "tank"},
		Metadata: map[string]any{
			"capacity_liters": 1500.75,
			"has_sensor":      true,
			"material":        "Composite",
		},
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	},
	{
		Uuid:          "7d444840-9dc0-11d1-b245-5ffdce74fad2",
		Name:          "Wing Section W9",
		Description:   "Carbon fiber wing segment",
		Price:         54000,
		StockQuantity: 3,
		Category:      repoModel.CategoryWing,
		Dimensions: repoModel.Dimensions{
			Length: 10.0,
			Width:  3.0,
			Height: 0.5,
			Weight: 1200,
		},
		Manufacturer: repoModel.Manufacturer{
			Name:    "SkyWorks",
			Country: "UK",
			Website: "https://skyworks.example.com",
		},
		Tags: []string{"wing", "carbon"},
		Metadata: map[string]any{
			"max_load_tons":    18.5,
			"has_fuel_channel": true,
			"revision":         3,
		},
		CreatedAt: time.Now(),
		UpdatedAt: nil,
	},
}
