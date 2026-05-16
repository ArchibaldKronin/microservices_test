package session

import (
	"fmt"

	def "github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/cache"
)

var _ def.SessionRepository = (*sessionRepository)(nil)

const (
	cacheSessionKeyPrefix = "iam:sessionID:"
	cacheUserKeyPrefix    = "iam:userID:"
)

type sessionRepository struct {
	repo cache.RedisClient
}

func NewSessionRepository(
	repo cache.RedisClient,
) *sessionRepository {
	return &sessionRepository{
		repo: repo,
	}
}

func getSessionCacheKay(uuid string) string {
	return fmt.Sprintf("%s%s", cacheSessionKeyPrefix, uuid)
}

func getUserCacheKay(uuid string) string {
	return fmt.Sprintf("%s%s", cacheUserKeyPrefix, uuid)
}
