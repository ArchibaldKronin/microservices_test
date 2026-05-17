package user

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "error getting user", zap.String("userID", userID), zap.Error(err))

		if errors.Is(err, repoModel.ErrNotFound) {
			return nil, model.ErrNotFound
		}
		return nil, model.ErrInternal
	}

	return user, nil
}
