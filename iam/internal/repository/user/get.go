package user

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/jackc/pgx/v5"

	def "github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) GetUser(ctx context.Context, userID string) (*model.User, error) {
	buildSelectOne := sq.Select(
		"user_id",
		"login",
		"email",
		"notification_methods",
	).
		From(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"user_id": userID}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		logger.Error(ctx, "repo error build sql", zap.Error(err))
		return nil, err
	}

	user := repoModel.User{}
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&user.UserUUID,
		&user.Login,
		&user.Email,
		&user.NotificationMethods,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoModel.ErrNotFound
		}
		logger.Error(ctx, "repo error get user", zap.String("userID", userID), zap.Error(err))
		return nil, repoModel.ErrSelectQuery
	}

	result, err := converter.UserToDomain(user)
	if err != nil {
		logger.Error(ctx, "repo error in user converter", zap.String("userID", userID), zap.Error(err))
		return nil, repoModel.ErrConverter
	}

	return &result, nil
}
