package service

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
)

type UserService interface {
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetCredentials(ctx context.Context, login string) (id string, pw string, err error)
	RegisterUser(ctx context.Context, user model.User, pw string) (string, error)
}

type AuthService interface {
	Login(ctx context.Context, login string, pw string) (string, error)
	Whoami(ctx context.Context, sessionID string) (*model.User, error)
}
