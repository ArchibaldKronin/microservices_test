package order

import (
	"context"

	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	// mu   sync.RWMutex
	pgRepo *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pgRepo: pool,
	}
	// return &repository{
	// 	data: make(map[string]model.Order),
	// }
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pgRepo.Begin(ctx)
}

// type repository struct {
// 	mu   sync.RWMutex
// 	data map[string]model.Order
// }

// func NewRepository() *repository {
// 	return &repository{
// 		data: make(map[string]model.Order),
// 	}
// }
