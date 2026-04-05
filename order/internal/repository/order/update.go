package order

import (
	"context"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
)

func (r *repository) UpdateOrder(_ context.Context, o *serviceModel.Order) *serviceModel.Order {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.data[o.OrderId]
	if !ok {
		return nil
	}

	r.data[o.OrderId] = converter.OrderToRepo(*o)
	return o
}
