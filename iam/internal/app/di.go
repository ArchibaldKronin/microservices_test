package app

import (
	"context"

	authV1Api "github.com/ArchibaldKronin/microservices_test/iam/internal/api/auth/v1"
	userV1Api "github.com/ArchibaldKronin/microservices_test/iam/internal/api/user/v1"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/config"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/repository"
	authRepository "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/session"
	userRepository "github.com/ArchibaldKronin/microservices_test/iam/internal/repository/user"
	"github.com/ArchibaldKronin/microservices_test/iam/internal/service"
	authService "github.com/ArchibaldKronin/microservices_test/iam/internal/service/auth"
	userService "github.com/ArchibaldKronin/microservices_test/iam/internal/service/user"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/cache"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/cache/redis"
	"github.com/ArchibaldKronin/microservices_test/platform/pkg/logger"
	auth_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/auth/v1"
	user_v1 "github.com/ArchibaldKronin/microservices_test/shared/pkg/proto/user/v1"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type diContainer struct {
	userV1Api      user_v1.UserServiceServer
	userService    service.UserService
	userRepository repository.UserRepository

	userPgPool *pgxpool.Pool

	authV1Api             auth_v1.AuthServiceServer
	authService           service.AuthService
	authSessionRepository repository.SessionRepository

	authRedisClient cache.RedisClient
	authRedisPool   *redigo.Pool
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) UserV1Api(ctx context.Context) (user_v1.UserServiceServer, error) {
	if d.userV1Api == nil {
		service, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}

		d.userV1Api = userV1Api.NewApiUser(service)
	}
	return d.userV1Api, nil
}

func (d *diContainer) UserService(ctx context.Context) (service.UserService, error) {
	if d.userService == nil {
		repo, err := d.UserRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.userService = userService.NewService(repo)
	}
	return d.userService, nil
}

func (d *diContainer) UserRepository(ctx context.Context) (repository.UserRepository, error) {
	if d.userRepository == nil {
		pool, err := d.UserPgPool(ctx)
		if err != nil {
			return nil, err
		}

		d.userRepository = userRepository.NewRepository(pool)
	}
	return d.userRepository, nil
}

func (d *diContainer) UserPgPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.userPgPool == nil {
		pool, err := pgxpool.New(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			return nil, err
		}

		if err = pool.Ping(ctx); err != nil {
			return nil, err
		}

		d.userPgPool = pool
	}
	return d.userPgPool, nil
}

func (d *diContainer) AuthV1Api(ctx context.Context) (auth_v1.AuthServiceServer, error) {
	if d.authV1Api == nil {
		service, err := d.AuthService(ctx)
		if err != nil {
			return nil, err
		}

		d.authV1Api = authV1Api.NewApiAuth(service)
	}
	return d.authV1Api, nil
}

func (d *diContainer) AuthService(ctx context.Context) (service.AuthService, error) {
	if d.authService == nil {
		repo, err := d.AuthSessionRepository(ctx)
		if err != nil {
			return nil, err
		}

		userServ, err := d.UserService(ctx)
		if err != nil {
			return nil, err
		}

		ttl := config.AppConfig().Redis.CacheTTL()

		d.authService = authService.NewAuthService(repo, userServ, ttl)
	}
	return d.authService, nil
}

func (d *diContainer) AuthSessionRepository(ctx context.Context) (repository.SessionRepository, error) {
	if d.authSessionRepository == nil {
		client, err := d.AuthRedisClient(ctx)
		if err != nil {
			return nil, err
		}

		d.authSessionRepository = authRepository.NewSessionRepository(client)
	}
	return d.authSessionRepository, nil
}

func (d *diContainer) AuthRedisClient(ctx context.Context) (cache.RedisClient, error) {
	if d.authRedisClient == nil {
		pool, err := d.AuthRedisPool(ctx)
		if err != nil {
			return nil, err
		}

		client := redis.NewClient(pool, logger.Logger(), config.AppConfig().Redis.ConnectionTimeout())
		if err = client.Ping(ctx); err != nil {
			logger.Error(ctx, "Faild to PING Redis", zap.Error(err))
			return nil, err
		}
		logger.Info(ctx, "Successfuly connect to Redis", zap.String("Address", config.AppConfig().Redis.Address()))

		d.authRedisClient = client
	}
	return d.authRedisClient, nil
}

func (d *diContainer) AuthRedisPool(ctx context.Context) (*redigo.Pool, error) {
	if d.authRedisPool == nil {
		d.authRedisPool = &redigo.Pool{
			DialContext: func(ctx context.Context) (redigo.Conn, error) {
				return redigo.DialContext(ctx, "tcp", config.AppConfig().Redis.Address())
			},
			MaxIdle:     config.AppConfig().Redis.MaxIdle(),
			IdleTimeout: config.AppConfig().Redis.IdleTimeout(),
		}
	}
	return d.authRedisPool, nil
}
