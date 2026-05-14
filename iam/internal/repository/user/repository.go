package user

import (
	def "github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *repository {
	return &repository{pool: pool}
}
