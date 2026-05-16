package session

import (
	"context"
	"time"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository/converter"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (r *sessionRepository) CreateUserSession(ctx context.Context, sessionID string, user model.User, ttl time.Duration) error {
	cacheKey := getSessionCacheKay(sessionID)
	redisView, err := converter.UserToRedisView(user)
	if err != nil {
		logger.Error(ctx, "error create redis model", zap.Error(err))
		return err
	}
	err = r.repo.HashSet(ctx, cacheKey, redisView)
	if err != nil {
		logger.Error(ctx, "error set model into redis", zap.Error(err))
		return err
	}

	err = r.repo.Expire(ctx, cacheKey, ttl)
	if err != nil {
		logger.Error(ctx, "error set TTL", zap.Error(err))
		return err
	}

	return nil
}
