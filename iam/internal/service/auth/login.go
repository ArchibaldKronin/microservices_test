package auth

import (
	"context"
	"errors"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (s *authService) Login(ctx context.Context, login string, pw string) (string, error) {
	// сверить пароль
	userID, hashed, err := s.userService.GetCredentials(ctx, login)
	if err != nil {
		logger.Error(ctx, "error geting credentials", zap.Error(err))
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrAutheticationData
		}
		return "", model.ErrInternal
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashed), []byte(pw))
	if err != nil {
		logger.Error(ctx, "error login", zap.Error(err))
		return "", model.ErrAutheticationData
	}

	//sessionID
	sessionID := uuid.NewString()
	user, err := s.userService.GetUserByID(ctx, userID)
	if err != nil {
		logger.Error(ctx, "error getting user by ID", zap.Error(err))
		return "", model.ErrInternal
	}

	//create session in Redis
	err = s.sessionRepo.CreateUserSession(ctx, sessionID, *user, s.cacheTTL)
	if err != nil {
		logger.Error(ctx, "error creating session", zap.Error(err))
		return "", model.ErrInternal
	}

	//add session in User's session set
	_ = s.sessionRepo.AddSessionToUserSet(ctx, user.UserUUID, sessionID, s.cacheTTL)

	//return sessionID
	return sessionID, nil
}
