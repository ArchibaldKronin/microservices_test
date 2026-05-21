package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	OTLPAddress() string
	ServiceName() string
	ServiceEnvironment() string
	EnableOTLP() bool
}

type InventoryGRPCConfig interface {
	Address() string
}

type IamGRPCClientConfig interface {
	Address() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}

type TracingConfig interface {
	CollectorEndpoint() string
	ServiceName() string
	Environment() string
	ServiceVersion() string
}
