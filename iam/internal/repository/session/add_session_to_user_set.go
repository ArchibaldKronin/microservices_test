package session

import (
	"context"
	"time"

	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (r *sessionRepository) AddSessionToUserSet(ctx context.Context, userID, sessionID string, ttl time.Duration) error {
	cacheKey := getUserCacheKay(userID)

	err := r.repo.SAdd(ctx, cacheKey, sessionID)
	if err != nil {
		return err
	}

	err = r.repo.Expire(ctx, cacheKey, ttl)
	if err != nil {
		logger.Error(ctx, "error set TTL", zap.Error(err))
		return err
	}

	return nil

}
