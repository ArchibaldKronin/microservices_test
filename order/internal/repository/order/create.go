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
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *repository) CreateOrder(ctx context.Context, o *serviceModel.Order) error {
	orderRepo := converter.OrderToRepo(*o)

	buildInsertOne := sq.Insert(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Columns(def.RepoFields...).
		Values(
			orderRepo.OrderID,
			orderRepo.UserID,
			orderRepo.PartIDs,
			orderRepo.TotalPrice,
			orderRepo.TransactionID,
			orderRepo.PaymentMethod,
			orderRepo.Status,
		)

	query, args, err := buildInsertOne.ToSql()
	if err != nil {
		logger.Error(ctx, "repo error build sql", zap.Error(err))
		return model.ErrBuildQuery
	}

	_, err = r.pgExecuter.Exec(ctx, query, args...)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Error(ctx, "repo error duplicate", zap.String("order_ID", orderRepo.OrderID), zap.Error(err))
				return model.ErrDuplicate
			}
		}

		logger.Error(ctx, "repo error create order", zap.String("order_ID", orderRepo.OrderID), zap.Error(err))
		return model.ErrExecQuery
	}

	return nil
}
