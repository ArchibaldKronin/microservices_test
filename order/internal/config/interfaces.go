package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJSON() bool
}

type OrderHTTPConfigConfig interface {
	Address() string
	ReadTime() time.Duration
	Port() string
}

type InventoryGRPCClientConfig interface {
	Address() string
}

type PaymentGRPCClientConfig interface {
	Address() string
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}
