module github.com/kiselevos/new_tax

go 1.23.10

require (
	github.com/joho/godotenv v1.5.1
	github.com/kiselevos/new_tax/gen v0.0.0-00010101000000-000000000000
	github.com/kiselevos/new_tax/pkg/logx v0.0.0
	github.com/oapi-codegen/runtime v1.1.2
	github.com/stretchr/testify v1.11.1
	golang.org/x/time v0.5.0
	google.golang.org/grpc v1.70.0
	google.golang.org/protobuf v1.36.5
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/redis/go-redis/v9 v9.17.2 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	go.opentelemetry.io/otel v1.34.0 // indirect
	golang.org/x/net v0.37.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250219182151-9fdb1cabc7b2 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kiselevos/new_tax/gen => ./gen

replace github.com/kiselevos/new_tax/pkg/logx => ./pkg/logx
