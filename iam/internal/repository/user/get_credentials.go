package user

import (
	"context"
	"errors"

	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/jackc/pgx/v5"

	def "github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *repository) GetCredentials(ctx context.Context, login string) (id string, pw string, err error) {
	buildSelectOne := sq.Select(
		"user_id",
		"password_hash",
	).
		From(def.TABLE_NAME).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"login": login}).
		Limit(1)

	query, args, err := buildSelectOne.ToSql()
	if err != nil {
		logger.Error(ctx, "repo error build sql", zap.Error(err))
		return "", "", err
	}

	var lg repoModel.LoginCredentials
	err = r.pool.QueryRow(ctx, query, args...).Scan(
		&lg.Id,
		&lg.Pw,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", repoModel.ErrNotFound
		}
		logger.Error(ctx, "repo error get credentials", zap.String("login", login), zap.Error(err))
		return "", "", repoModel.ErrSelectQuery
	}

	return lg.Id, lg.Pw, nil
}
