package model

import "github.com/google/uuid"

type Order struct {
	OrderID       string
	UserID        string
	PartIDs       []string
	TotalPrice    float64
	TransactionID *string
	PaymentMethod *PaymentMethod
	Status        OrderStatus
}

func NewOrder(userId string, partIds []string, totalPrice float64) *Order {
	return &Order{
		OrderID:    uuid.NewString(),
		UserID:     userId,
		PartIDs:    partIds,
		TotalPrice: totalPrice,
		Status:     OrderStatusPENDINGPAYMENT,
	}
}
