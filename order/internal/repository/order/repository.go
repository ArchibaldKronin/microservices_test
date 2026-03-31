package order

import (
	"sync"

	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]model.Order
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]model.Order),
	}
}
