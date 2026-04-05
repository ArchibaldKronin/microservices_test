package v1

import (
	"github.com/ArchibaldKronin/microservices_test/order/internal/service"
	def "github.com/ArchibaldKronin/microservices_test/shared/pkg/openapi/order/v1"
)

var _ def.Handler = (*api)(nil)

type api struct {
	orderService service.OrderService
}

func NewApi(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
