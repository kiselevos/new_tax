default: help
#.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: codegen
codegen: ## Generate gRPC code
	@echo "Generating gRPC code"
	@mkdir -p gen/grpc/
	@cd api && protoc \
		-I . \
		-I /usr/include \
		-I /usr/local/include \
		-I $(shell go list -f '{{ .Dir }}' -m google.golang.org/protobuf)/../ \
		--experimental_allow_proto3_optional \
		--go_out=../gen/grpc --go_opt=paths=source_relative \
		--go-grpc_out=../gen/grpc --go-grpc_opt=paths=source_relative \
		tax.proto
	@echo "✅ gRPC code generated"

.PHONY: run
run: ## Run application
	@go run ./cmd/main.go

.PHONY: tidy
tidy: ## Check go.mod/go.sum
	@go mod tidy
	@git diff --exit-code || (echo "::error::go.mod or go.sum is out of sync" && exit 1)

.PHONY: test-all
test-all: ## Run tests
	@go test -v ./...

.PHONY: build
build: ## Compile Go binary
	@echo "Building binary..."
	@go build -o bin/tax ./cmd
	@echo "✅ Binary built at bin/tax"

.PHONY: docker-up
docker-up: ## Start docker containers without build
	@docker compose up -d

.PHONY: docker-build
docker-build: ## Build and start docker containers
	@docker compose up --build

.PHONY: docker-down
docker-down: ## Down docker conteiners
	@docker compose down

.PHONY: lint-all
lint-all: ## Run all linters
	@go vet ./...
	@golangci-lint run ./...

.PHONY: gofmt
gofmt: ## Format code
	@gofmt -s -w . 

.PHONY: check-generated
check-generated: ## Just check git diff
	@git diff --exit-code gen || (echo "❌ Generated files are not committed" && exit 1)

.PHONY: check-fmt
check-fmt: ## Check formatting (CI)
	@gofmt -l . | grep -q . && (echo "❌ Files need formatting (run 'make fmt')"; exit 1) || true

.PHONY: local-CI
local-CI: ## Use act to check CI local
	@act -P ubuntu-latest=catthehacker/ubuntu:act-22.04