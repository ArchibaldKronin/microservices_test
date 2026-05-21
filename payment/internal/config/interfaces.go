package config

type LoggerConfig interface {
	Level() string
	AsJSON() bool
	OTLPAddress() string
	ServiceName() string
	ServiceEnvironment() string
	EnableOTLP() bool
}

type PaymentGRPCConfig interface {
	Address() string
}
