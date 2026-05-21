package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	OTLPAddress() string
	ServiceName() string
	ServiceEnvironment() string
	EnableOTLP() bool
}

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type RedisConfig interface {
	Address() string
	ConnectionTimeout() time.Duration
	MaxIdle() int
	IdleTimeout() time.Duration
	CacheTTL() time.Duration
}
