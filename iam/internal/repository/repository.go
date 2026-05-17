package repository

import (
	"context"
	"time"

	"github.com/ArchibaldKronin/microservices_test/iam/internal/model"
)

const TABLE_NAME = "users"

type UserRepository interface {
	Create(ctx context.Context, user model.User, passwordHash string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetCredentials(ctx context.Context, login string) (id, pw string, err error)
}

type SessionRepository interface {
	CreateUserSession(ctx context.Context, sessionID string, user model.User, ttl time.Duration) error
	GetUserBySession(ctx context.Context, sessionID string) (*model.User, error)
	AddSessionToUserSet(ctx context.Context, userID, sessionID string, ttl time.Duration) error
}

/*


 */
