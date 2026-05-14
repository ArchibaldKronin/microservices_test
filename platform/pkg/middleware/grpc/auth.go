package grpc

import (
	"context"
	"fmt"

	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	common_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/common/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// SessionUUIDMetadataKey ключ для передачи UUID сессии в gRPC metadata
	SessionUUIDMetadataKey = "session-uuid"
)

type contextKey string

const (
	// userContextKey ключ для хранения пользователя в контексте
	userContextKey contextKey = "user"
	// sessionUUIDContextKey ключ для хранения session UUID в контексте
	sessionUUIDContextKey contextKey = "session-uuid"
)

// IAMClient это алиас для сгенерированного gRPC клиента
type IAMClient = auth_v1.AuthServiceClient

// AuthInterceptor interceptor для аутентификации gRPC запросов
type AuthInterceptor struct {
	iamClient IAMClient
}

func NewAuthInterceptor(iamClient IAMClient) *AuthInterceptor {
	return &AuthInterceptor{iamClient: iamClient}
}

// Unary возвращает unary server interceptor для аутентификации
func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := i.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

// authenticate выполняет аутентификацию и добавляет пользователя в контекст
func (i *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	// Извлекаем metadata из контекста
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	// Получаем session UUID из metadata
	sessionUUIDs := md.Get(SessionUUIDMetadataKey)
	if len(sessionUUIDs) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing session-uuid in metadata")
	}

	sessionUUID := sessionUUIDs[0]
	if sessionUUID == "" {
		return nil, status.Error(codes.Unauthenticated, "empty session-uuid")
	}

	// Валидируем сессию через IAM сервис
	whoamiRes, err := i.iamClient.Whoami(ctx, &auth_v1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid session: %v", err))
	}

	// Добавляем пользователя и session UUID в контекст
	authCtx := context.WithValue(ctx, userContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, sessionUUIDContextKey, sessionUUID)

	return authCtx, nil
}

func GetUserFromContext(ctx context.Context) (*common_v1.User, bool) {
	user, ok := ctx.Value(userContextKey).(*common_v1.User)
	return user, ok
}

// GetUserContextKey возвращает ключ контекста для пользователя
func GetUserContextKey() contextKey {
	return userContextKey
}

// GetSessionUUIDFromContext извлекает session UUID из контекста
func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}

// AddSessionUUIDToContext добавляет session UUID в контекст
func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionUUIDContextKey, sessionUUID)
}

// AddUserToContext добавляет User в контекст
func AddUserToContext(ctx context.Context, user *common_v1.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// ForwardSessionUUIDToGRPC добавляет session UUID из контекста в исходящие gRPC metadata
func ForwardSessionUUIDToGRPC(ctx context.Context) context.Context {
	sessionUUID, ok := GetSessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionUUIDMetadataKey, sessionUUID)
}
