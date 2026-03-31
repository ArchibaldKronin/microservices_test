package repository

import serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"

type OrderRepository interface {
	CreateOrder(o *serviceModel.Order)
	UpdateOrder(o *serviceModel.Order)
	GetOrder(id string) *serviceModel.Order
}
