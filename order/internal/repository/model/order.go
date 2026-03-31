package model

type Order struct {
	OrderId       string
	UserId        string
	PartIds       []string
	Total_price   float64
	TransactionID *string
	PaymentMethod *PaymentMethod
	Status        OrderStatus
}
