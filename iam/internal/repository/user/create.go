package user

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	def "github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

func (r *repository) Create(ctx context.Context, user model.User, passwordHash string) (string, error) {
	userRepo, err := converter.UserToRepo(user, passwordHash)
	if err != nil {
		logger.Error(ctx, "repo error in user converter", zap.String("userID", user.UserUUID), zap.Error(err))
		return "", repoModel.ErrConverter
	}

	buildInsertOne := sq.Insert(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Columns(
			"user_id",
			"login",
			"email",
			"password_hash",
			"notification_methods",
		).
		Values(
			userRepo.UserUUID,
			userRepo.Login,
			userRepo.Email,
			userRepo.PasswordHash,
			userRepo.NotificationMethods,
		)

	query, args, err := buildInsertOne.ToSql()
	if err != nil {
		logger.Error(ctx, "repo error build sql", zap.Error(err))
		return "", repoModel.ErrBuildQuery
	}

	_, err = r.pool.Exec(ctx, query, args...)
	if err != nil {

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				logger.Error(ctx, "repo error duplicate", zap.String("userID", user.UserUUID), zap.Error(err))
				return "", repoModel.ErrDuplicate
			}
		}

		logger.Error(ctx, "repo error register user", zap.String("userID", user.UserUUID), zap.Error(err))
		return "", repoModel.ErrExecQuery
	}

	return user.UserUUID, nil
}
