package model

type (
	OrderStatus   string
	PaymentMethod string
)

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"

	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

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

type PartsFilter struct {
	Uuids     []string
	Names     []string
	Categorys []Category
	Countrys  []string
	Tags      []string
}
