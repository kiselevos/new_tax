# ============================================================
# 🌐 Web Service Makefile
# ============================================================

# Версии инструментов
PROTOC_GEN_GO_VERSION := v1.34.2
PROTOC_GEN_CONNECT_GO_VERSION := v1.16.0
GOLANGCI_LINT_VERSION := v1.60.3

# ============================================================
# 🧠 Default & Help
# ============================================================
default: help

.PHONY: help
help: ## Display available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'


# ============================================================
# ⚙️ Setup & Dependencies
# ============================================================
.PHONY: setup-frontend
setup-frontend: ## Устанавливает зависимости для frontend-приложения
	@echo "🎨 Installing Frontend..."
	@go mod tidy
	@go mod download
	@echo "✅ Frontend dependencies installed"


# ============================================================
# 🧬 Code Generation
# ============================================================
.PHONY: codegen
codegen: ## Generate ConnectRPC client code from .proto
	@echo "🧬 Generating ConnectRPC client code..."
	@mkdir -p ./gen/grpc/api
	@protoc \
	  --proto_path=../docs/grpc \
	  --plugin=protoc-gen-go=$(shell which protoc-gen-go) \
	  --plugin=protoc-gen-connect-go=$(shell which protoc-gen-connect-go) \
	  --go_out=./gen/grpc/api --go_opt=paths=source_relative \
	  --connect-go_out=./gen/grpc/api --connect-go_opt=paths=source_relative \
	  $(shell find ../docs/grpc -name "*.proto" -type f)
	@echo "✅ ConnectRPC client code generated → gen/grpc/api"


# ============================================================
# 🧱 Build & Run
# ============================================================
.PHONY: run
run: ## Run web server locally
	@echo "🚀 Starting web server..."
	@go run ./cmd/web.go

.PHONY: build
build: ## Build binary
	@echo "🏗️  Building web binary..."
	@go build -o bin/web ./cmd/web.go
	@echo "✅ Binary built at bin/web"


# ============================================================
# 🧪 Tests & Checks
# ============================================================
.PHONY: lint
lint: ## Run linter
	@echo "🔍 Running linters..."
	@golangci-lint run ./...

.PHONY: test
test: ## Run tests
	@go test -v ./...

.PHONY: fmt
fmt: ## Format code
	@gofmt -s -w .


# ============================================================
# 🐳 Docker Commands
# ============================================================
.PHONY: docker-build
docker-build: ## Build docker image for web
	@docker build -t tax-web .

.PHONY: docker-run
docker-run: ## Run docker container for web
	@docker run -p 8080:8080 tax-web
