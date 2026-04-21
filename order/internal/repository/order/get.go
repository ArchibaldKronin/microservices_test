package order

import (
	"context"
	"errors"
	"log/slog"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
)

// func (r *repository) GetOrder(ctx context.Context, id string) (*serviceModel.Order, error) {
// 	return r.getOrder(ctx, r.pgRepo, id)
// }

func (r *repository) GetOrder(ctx context.Context, id string) (*serviceModel.Order, error) {
	buildSelectOne := sq.Select(def.RepoFields...).
		From(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_id": id}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		slog.Error("repo error build sql", "error ", err)
		return nil, model.ErrBuildQuery
	}

	order := model.Order{}
	err = r.pgExecuter.QueryRow(ctx, query, args...).Scan(
		&order.OrderID,
		&order.UserID,
		&order.PartIDs,
		&order.TotalPrice,
		&order.TransactionID,
		&order.PaymentMethod,
		&order.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		slog.Error("repo error get order", "id", order.OrderID, "error ", err)
		return nil, model.ErrSelectQuery
	}

	return lo.ToPtr(converter.OrderToDomain(order)), nil
}

// func (r *repository) GetOrderTx(ctx context.Context, tx pgx.Tx, id string) (*serviceModel.Order, error) {
// 	return r.getOrder(ctx, tx, id)
// }

// func (r *repository) getOrder(ctx context.Context, q interface {
// 	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
// }, id string,
// ) (*serviceModel.Order, error) {
// 	buildSelectOne := sq.Select(def.RepoFields...).
// 		From(def.TABLE_NAME).
// 		PlaceholderFormat(sq.Dollar).
// 		Where(sq.Eq{"order_id": id}).
// 		Limit(1)

// 	query, args, err := buildSelectOne.ToSql()
// 	if err != nil {
// 		slog.Error("repo error build sql", "error ", err)
// 		return nil, model.ErrBuildQuery
// 	}

// 	order := model.Order{}
// 	err = q.QueryRow(ctx, query, args...).Scan(
// 		&order.OrderID,
// 		&order.UserID,
// 		&order.PartIDs,
// 		&order.TotalPrice,
// 		&order.TransactionID,
// 		&order.PaymentMethod,
// 		&order.Status,
// 	)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			return nil, model.ErrNotFound
// 		}
// 		slog.Error("repo error get order", "id", order.OrderID, "error ", err)
// 		return nil, model.ErrSelectQuery
// 	}

// 	return lo.ToPtr(converter.OrderToDomain(order)), nil
// }
