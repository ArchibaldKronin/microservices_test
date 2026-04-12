package model

type Category int

const (
	CategoryUnknown Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

type Dimensions struct {
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

// type Value interface {
// 	isValue()
// }

// type StringValue struct {
// 	Value string
// }

// func (StringValue) isValue() {}

// type Int64Value struct {
// 	Value int64
// }

// func (Int64Value) isValue() {}

// type DoubleValue struct {
// 	Value float64
// }

// func (DoubleValue) isValue() {}

// type BoolValue struct {
// 	Value bool
// }

// func (BoolValue) isValue() {}
