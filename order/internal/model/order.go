package model

import "github.com/google/uuid"

type Order struct {
	OrderId       string
	UserId        string
	PartIds       []string
	Total_price   float64
	TransactionID *string
	PaymentMethod *PaymentMethod
	Status        OrderStatus
}

func NewOrder(userId string, partIds []string, totalPrice float64) *Order {
	return &Order{
		OrderId:     uuid.NewString(),
		UserId:      userId,
		PartIds:     partIds,
		Total_price: totalPrice,
		Status:      OrderStatusPENDINGPAYMENT,
	}
}
