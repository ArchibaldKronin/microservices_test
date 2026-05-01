package order

import (
	"context"

	serviceModel "github.com/ArchibaldKronin/microservices_test/order/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/order/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/converter"
	"github.com/ArchibaldKronin/microservices_test/order/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) UpdateOrder(ctx context.Context, o *serviceModel.Order) error {
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
		logger.Error(ctx, "repo error build sql", zap.Error(err))
		return model.ErrBuildQuery
	}

	ct, err := r.pgExecuter.Exec(ctx, query, args...)
	if err != nil {
		logger.Error(ctx, "repo error update order", zap.String("order_ID", orderRepo.OrderID), zap.Error(err))
		return model.ErrExecQuery
	}
	if ct.RowsAffected() != 1 {
		return model.ErrNotFound
	}

	return nil
}
