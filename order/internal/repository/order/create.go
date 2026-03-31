package order

import (
	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
)

func (r *repository) CreateOrder(o *serviceModel.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[o.OrderId] = converter.OrderToRepo(*o)
}
