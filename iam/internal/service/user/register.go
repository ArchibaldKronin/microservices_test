package user

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *service) RegisterUser(ctx context.Context, user model.User, pw string) (string, error) {
	userID := uuid.NewString()
	user.UserUUID = userID

	byteHash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(ctx, "error create password hash", zap.Error(err))
		return "", model.ErrHashPassword
	}

	userID, err = s.repo.Create(ctx, user, string(byteHash))
	if err != nil {
		logger.Error(ctx, "error register user", zap.String("userID", userID), zap.Error(err))
		return "", model.ErrInternal
	}

	return userID, nil
}
