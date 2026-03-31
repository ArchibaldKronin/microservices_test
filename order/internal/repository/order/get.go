package order

import (
	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
	"github.com/samber/lo"
)

func (r *repository) GetOrder(id string) *serviceModel.Order {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[id]
	if !ok {
		return nil
	}
	return lo.ToPtr(converter.OrderToDomain(order))
}
