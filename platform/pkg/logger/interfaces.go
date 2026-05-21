package logger

type initConfig interface {
	OTLPAddress() string
	ServiceName() string
	ServiceEnvironment() string
	AsJSON() bool
	EnableOTLP() bool
}
