package repository

import (
	"context"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
)

type UserRepository interface {
	Register(ctx context.Context, user model.User, passwordHash string) (string, error)
	GetUser(ctx context.Context, userID string) (*model.User, error)
}

type SessionRepository interface {
}

const TABLE_NAME = "users"
