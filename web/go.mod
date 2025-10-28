module github.com/kiselevos/new_tax/web

go 1.23.10

require (
	github.com/joho/godotenv v1.5.1
	github.com/kiselevos/new_tax/gen v0.0.0-00010101000000-000000000000
	github.com/kiselevos/new_tax/pkg/logx v0.0.0
	google.golang.org/grpc v1.68.0
	google.golang.org/protobuf v1.34.2
)

require (
	golang.org/x/net v0.29.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
	golang.org/x/text v0.18.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240903143218-8af14fe29dc1 // indirect
)

replace github.com/kiselevos/new_tax/gen => ../gen

replace github.com/kiselevos/new_tax/pkg/logx => ../pkg/logx
