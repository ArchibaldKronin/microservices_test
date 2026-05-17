package user

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	repoModel "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"go.uber.org/zap"
)

func (s *service) GetCredentials(ctx context.Context, login string) (id, pw string, err error) {
	id, pw, err = s.repo.GetCredentials(ctx, login)
	if err != nil {
		logger.Error(ctx, "error getting user", zap.String("login", login), zap.Error(err))

		if errors.Is(err, repoModel.ErrNotFound) {
			return "", "", model.ErrNotFound
		}
		return "", "", model.ErrInternal
	}

	return id, pw, err
}
