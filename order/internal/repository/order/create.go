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
	"github.com/jackc/pgx/v5/pgconn"
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
		slog.Error("repo error build sql", "error ", err)
		return model.ErrBuildQuery
	}

	_, err = r.pgRepo.Exec(ctx, query, args...)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				slog.Error("repo error duplicate", "id", orderRepo.OrderID, "error ", err)
				return model.ErrDuplicate
			}
		}

		slog.Error("repo error create order", "id", orderRepo.OrderID, "error ", err)
		return model.ErrExecQuery
	}

	return nil
}
