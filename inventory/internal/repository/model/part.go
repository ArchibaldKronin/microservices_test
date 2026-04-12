package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	Uuid          string             `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int                `bson:"stock_quantity"`
	Category      Category           `bson:"category"`
	Dimensions    Dimensions         `bson:"dimensions"`
	Manufacturer  Manufacturer       `bson:"manufacturer"`
	Tags          []string           `bson:"tags"`
	Metadata      map[string]any     `bson:"metadata"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty"`
}

// type Part struct {
// 	Uuid          string
// 	Name          string
// 	Description   string
// 	Price         float64
// 	StockQuantity intLfdfq
// 	Category      Category
// 	Dimensions    Dimensions
// 	Manufacturer  Manufacturer
// 	Tags          []string
// 	Metadata      map[string]Value
// 	CreatedAt     time.Time
// 	UpdatedAt     time.Time
// }
