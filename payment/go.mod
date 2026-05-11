module github.com/ArchibaldKronin/microservices_test/payment

go 1.26.1

replace github.com/ArchibaldKronin/microservices_test/shared => ../shared

replace github.com/ArchibaldKronin/microservices_test/platform => ../platform

require (
	github.com/ArchibaldKronin/microservices_test/platform v0.0.0-00010101000000-000000000000
	github.com/ArchibaldKronin/microservices_test/shared v0.0.0-00010101000000-000000000000
	github.com/brianvoe/gofakeit/v7 v7.14.1
	github.com/caarlos0/env/v11 v11.4.0
	github.com/google/uuid v1.6.0
	github.com/joho/godotenv v1.5.1
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.27.1
	google.golang.org/grpc v1.80.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/stretchr/objx v0.5.3 // indirect
	go.opentelemetry.io/otel/sdk v1.43.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260420184626-e10c466a9529 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
