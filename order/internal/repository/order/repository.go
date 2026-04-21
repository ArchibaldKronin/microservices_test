package order

import (
	"context"

	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var _ def.OrderRepository = (*repository)(nil)

type pgQueryInterface interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type repository struct {
	// pgRepo *pgxpool.Pool
	pgExecuter pgQueryInterface
}

// func NewRepository(pool *pgxpool.Pool) *repository {
func NewRepository(exec pgQueryInterface) *repository {
	return &repository{
		pgExecuter: exec,
	}
}

// func (r *repository) BeginTx(ctx context.Context) (pgx.Tx, error) {
// 	return r.pgRepo.Begin(ctx)
// }
