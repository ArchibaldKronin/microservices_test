package auth

import (
	"time"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	def "github.com/ArchibaldKronin/microservices_test/iam/internal/service"
)

var _ def.AuthService = (*authService)(nil)

type authService struct {
	sessionRepo repository.SessionRepository
	userService def.UserService

	cacheTTL time.Duration
}

func NewAuthService(
	sessionRepo repository.SessionRepository,
	userService def.UserService,

	cacheTTL time.Duration,
) *authService {
	return &authService{
		sessionRepo: sessionRepo,
		userService: userService,

		cacheTTL: cacheTTL,
	}
}
