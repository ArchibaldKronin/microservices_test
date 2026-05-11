module github.com/ArchibaldKronin/microservices_test/notification

go 1.26.1

require (
	github.com/ArchibaldKronin/microservices_test/platform v0.0.0-00010101000000-000000000000
	github.com/ArchibaldKronin/microservices_test/shared v0.0.0-00010101000000-000000000000
	github.com/IBM/sarama v1.48.0
	github.com/caarlos0/env/v11 v11.4.1
	github.com/go-telegram/bot v1.20.0
	github.com/joho/godotenv v1.5.1
	go.uber.org/zap v1.28.0
	google.golang.org/protobuf v1.36.11
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/eapache/go-resiliency v1.7.0 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.18.5 // indirect
	github.com/pierrec/lz4/v4 v4.1.26 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20250401214520-65e299d6c5c9 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/crypto v0.50.0 // indirect
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
)

replace github.com/ArchibaldKronin/microservices_test/shared => ../shared

replace github.com/ArchibaldKronin/microservices_test/platform => ../platform
