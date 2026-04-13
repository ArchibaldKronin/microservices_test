package order

import (
	"context"
	"log/slog"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *repository) UpdateOrder(ctx context.Context, o *serviceModel.Order) error {
	return r.updateOrder(ctx, r.pgRepo, o)
}

func (r *repository) UpdateOrderTx(ctx context.Context, tx pgx.Tx, o *serviceModel.Order) error {
	return r.updateOrder(ctx, tx, o)
}

func (r *repository) updateOrder(ctx context.Context, e interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
}, o *serviceModel.Order,
) error {
	orderRepo := converter.OrderToRepo(*o)

	buildUpdateOne := sq.Update(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Set("user_id", orderRepo.UserID).
		Set("part_ids", orderRepo.PartIDs).
		Set("total_price", orderRepo.TotalPrice).
		Set("transaction_id", orderRepo.TransactionID).
		Set("payment_method", orderRepo.PaymentMethod).
		Set("status", orderRepo.Status).
		Where(sq.Eq{"order_id": orderRepo.OrderID})

	query, args, err := buildUpdateOne.ToSql()
	if err != nil {
		slog.Error("repo error build sql", "error ", err)
		return model.ErrBuildQuery
	}

	ct, err := e.Exec(ctx, query, args...)
	if err != nil {
		slog.Error("repo error update order", "id", orderRepo.OrderID, "error ", err)
		return model.ErrExecQuery
	}
	if ct.RowsAffected() != 1 {
		return model.ErrNotFound
	}

	return nil
}
