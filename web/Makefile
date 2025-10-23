# === Web Service Makefile ===
.PHONY: run
run: ## Run web server locally
	go run cmd/web.go

.PHONY: build
build: ## Build binary
	go build -o bin/web cmd/web.go

.PHONY: fmt
fmt: ## Format code
	go fmt ./...

.PHONY: codegen
codegen: ## Generate Go & Connect code from proto
	@protoc \
		--proto_path=./docs/grpc \
		--go_out=./gen/grpc/api --go_opt=paths=source_relative \
		--connect-go_out=./gen/grpc/api --connect-go_opt=paths=source_relative \
		docs/grpc/tax.proto