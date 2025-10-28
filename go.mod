module github.com/kiselevos/new_tax

go 1.23.10

require (
	github.com/kiselevos/new_tax/gen v0.0.0-00010101000000-000000000000
	github.com/kiselevos/new_tax/pkg/logx v0.0.0
	github.com/oapi-codegen/runtime v1.1.2
	github.com/stretchr/testify v1.11.1
	google.golang.org/grpc v1.68.0
	google.golang.org/protobuf v1.34.2
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/kiselevos/new_tax/gen => ./gen

replace github.com/kiselevos/new_tax/pkg/logx => ./pkg/logx
