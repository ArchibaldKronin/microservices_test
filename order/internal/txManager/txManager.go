package txmanager

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ArchibaldKronin/microservices_test/order/internal/model"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/order"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// type callBackFunc func(executer repository.OrderRepository) error

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
		slog.Error("error beginning tx", "error", err)
		return model.ErrInternal
	}

	defer func() {
		rerr := tx.Rollback(ctx)
		if rerr != nil && !errors.Is(rerr, pgx.ErrTxClosed) {
			slog.Warn("rallback failed", "error", rerr)
		}
	}()

	exec := order.NewRepository(tx)

	err = fn(exec)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error("error committing tx:", "error", err)
		return model.ErrInternal
	}

	return nil
}

// func (m *txRepoManager) WithoutTransaction(ctx context.Context, fn callBackFunc) error {
// 	exec := order.NewRepository(m.pgRepo)

// 	err := fn(exec)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
