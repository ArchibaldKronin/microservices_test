module github.com/ArchibaldKronin/microservices_test/inventory

go 1.26.1

replace github.com/ArchibaldKronin/microservices_test/shared => ../shared

require (
	github.com/ArchibaldKronin/microservices_test/shared v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.79.3
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/brianvoe/gofakeit/v7 v7.14.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/stretchr/objx v0.5.2 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
