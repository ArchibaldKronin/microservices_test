package main

import "time"

var PartsTest = []*Part{
	{
		Uuid:          "550e8400-e29b-41d4-a716-446655440000",
		Name:          "Turbo Engine X1",
		Description:   "High performance aircraft engine",
		Price:         125000.50,
		StockQuantity: 5,
		Category:      CategoryEngine,
		Dimensions: Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.5,
			Weight: 850,
		},
		Manufacturer: Manufacturer{
			Name:    "JetCorp",
			Country: "USA",
			Website: "https://jetcorp.example.com",
		},
		Tags: []string{"engine", "turbo"},
		Metadata: map[string]Value{
			"horsepower":   Int64Value{Value: 4500},
			"fuel_type":    StringValue{Value: "Jet A-1"},
			"is_certified": BoolValue{Value: true},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Uuid:          "6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		Name:          "Fuel Tank F200",
		Description:   "Composite aircraft fuel tank",
		Price:         32000,
		StockQuantity: 12,
		Category:      CategoryFuel,
		Dimensions: Dimensions{
			Length: 3.0,
			Width:  1.5,
			Height: 1.5,
			Weight: 200,
		},
		Manufacturer: Manufacturer{
			Name:    "FuelTech",
			Country: "Germany",
			Website: "https://fueltech.example.com",
		},
		Tags: []string{"fuel", "tank"},
		Metadata: map[string]Value{
			"capacity_liters": DoubleValue{Value: 1500.75},
			"has_sensor":      BoolValue{Value: true},
			"material":        StringValue{Value: "Composite"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
	{
		Uuid:          "7d444840-9dc0-11d1-b245-5ffdce74fad2",
		Name:          "Wing Section W9",
		Description:   "Carbon fiber wing segment",
		Price:         54000,
		StockQuantity: 3,
		Category:      CategoryWing,
		Dimensions: Dimensions{
			Length: 10.0,
			Width:  3.0,
			Height: 0.5,
			Weight: 1200,
		},
		Manufacturer: Manufacturer{
			Name:    "SkyWorks",
			Country: "UK",
			Website: "https://skyworks.example.com",
		},
		Tags: []string{"wing", "carbon"},
		Metadata: map[string]Value{
			"max_load_tons":    DoubleValue{Value: 18.5},
			"has_fuel_channel": BoolValue{Value: true},
			"revision":         Int64Value{Value: 3},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	},
}
