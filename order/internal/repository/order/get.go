package order

import (
	"context"
	"errors"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func (r *repository) GetOrder(ctx context.Context, id string) (*serviceModel.Order, error) {
	buildSelectOne := sq.Select(def.RepoFields...).
		From(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"order_id": id}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		logger.Error(ctx, "repo error build sql", zap.Error(err))
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
		logger.Error(ctx, "repo error get order", zap.String("order_ID", order.OrderID), zap.Error(err))
		return nil, model.ErrSelectQuery
	}

	return lo.ToPtr(converter.OrderToDomain(order)), nil
}
