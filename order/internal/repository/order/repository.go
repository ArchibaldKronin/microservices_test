package order

import (
	"context"

	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	pgRepo *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{
		pgRepo: pool,
	}
}

func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.pgRepo.Begin(ctx)
}
