package auth

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *authService) Whoami(ctx context.Context, sessionID string) (*model.User, error) {
	user, err := s.sessionRepo.GetUserBySession(ctx, sessionID)
	if err != nil {
		logger.Error(ctx, "error geting user", zap.Error(err))
		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, model.ErrInternal
	}

	return user, nil
}
