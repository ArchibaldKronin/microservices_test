package order

import (
	"context"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
)

func (r *repository) CreateOrder(_ context.Context, o *serviceModel.Order) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[o.OrderId] = converter.OrderToRepo(*o)
}
