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
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value interface {
	isValue()
}

type StringValue struct {
	Value string
}

func (StringValue) isValue() {}

type Int64Value struct {
	Value int64
}

func (Int64Value) isValue() {}

type DoubleValue struct {
	Value float64
}

func (DoubleValue) isValue() {}

type BoolValue struct {
	Value bool
}

func (BoolValue) isValue() {}
