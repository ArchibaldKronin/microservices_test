package txmanager

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/order"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type TxManager interface {
	WithTransaction(ctx context.Context, fn func(executer repository.OrderRepository) error) error
	// WithoutTransaction(ctx context.Context, fn callBackFunc) error
}

var _ TxManager = (*txRepoManager)(nil)

type txRepoManager struct {
	pgRepo *pgxpool.Pool
}

func NewTxRepoManager(pool *pgxpool.Pool) *txRepoManager {
	return &txRepoManager{
		pgRepo: pool,
	}
}

func (m *txRepoManager) WithTransaction(ctx context.Context, fn func(executer repository.OrderRepository) error) error {
	tx, err := m.pgRepo.Begin(ctx)
	if err != nil {
		logger.Error(ctx, "error beginning tx", zap.Error(err))

		return model.ErrInternal
	}

	defer func() {
		rerr := tx.Rollback(ctx)
		if rerr != nil && !errors.Is(rerr, pgx.ErrTxClosed) {
			logger.Error(ctx, "rallback failed", zap.Error(err))
		}
	}()

	exec := order.NewRepository(tx)

	err = fn(exec)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logger.Error(ctx, "error committing tx", zap.Error(err))
		return model.ErrInternal
	}

	return nil
}
