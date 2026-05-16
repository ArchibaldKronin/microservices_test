package session

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository/converter"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	redigo "github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

func (r *sessionRepository) GetUserBySession(ctx context.Context, sessionID string) (*model.User, error) {
	cacheKey := getSessionCacheKay(sessionID)

	values, err := r.repo.HGetAll(ctx, cacheKey)
	if err != nil {
		logger.Error(ctx, "error get user by sessionID", zap.String("sessionID", sessionID), zap.Error(err))
		if errors.Is(err, redigo.ErrNil) {
			return nil, repoModel.ErrNotFound
		}
		return nil, err
	}

	if len(values) == 0 {
		return nil, repoModel.ErrNotFound
	}

	var userRedisView repoModel.UserRedisView
	err = redigo.ScanStruct(values, &userRedisView)
	if err != nil {
		logger.Error(ctx, "error scan user structure", zap.String("sessionID", sessionID), zap.Error(err))
		return nil, err
	}

	user, err := converter.UserFromRedisView(userRedisView)
	if err != nil {
		logger.Error(ctx, "error create domain model", zap.Error(err))
		return nil, err
	}

	return &user, nil
}
