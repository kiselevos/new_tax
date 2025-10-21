# Version variables
BUF_VERSION := 1.38.0
PROTOC_GO_VERSION := v1.34.2
PROTOC_CONNECT_GO_VERSION := v1.16.0
PROTOC_ES_VERSION := 1.10.1
PROTOC_CONNECT_ES_VERSION := 1.7.0
PROTOBUF_JS_VERSION := 1.10.1
CONNECT_VERSION := 1.7.0
CONNECT_WEB_VERSION := 1.7.0

default: help
#.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: setup
setup: setup-buf setup-go-tools setup-node-deps ## Install all development tools locally
	@echo "🎉 All tools installed! Run 'make codegen' to generate code."

.PHONY: setup-buf
setup-buf: ## Install Buf locally
	@echo "📦 Installing Buf v$(BUF_VERSION)..."
	@cd web && npm install --save-dev @bufbuild/buf@$(BUF_VERSION)
	@echo "✅ Buf installed"

.PHONY: setup-go-tools
setup-go-tools: ## Install Go tools locally
	@echo "🔧 Installing Go plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)
	@go install connectrpc.com/connect/cmd/protoc-gen-connect-go@$(PROTOC_CONNECT_GO_VERSION)
	@echo "✅ Go plugins installed"

.PHONY: setup-node-deps
setup-node-deps: ## Install Node.js dependencies
	@echo "📦 Installing Node.js dependencies..."
	@cd web && npm install
	@echo "✅ Node.js dependencies installed"


# --- Backend commands ---
.PHONY: codegen
codegen: codegen-backend codegen-frontend ## Generate all GPRC code

.PHONY: codegen-backend
codegen-backend: ## Generate only backend Go code
	@echo "🔨 Generating backend Go code..."
	@mkdir -p gen/grpc/api
	@buf generate --template buf.gen.backend.local.yaml ./api
	@echo "✅ Backend code generated in gen/grpc/api"

.PHONY: codegen-frontend
codegen-frontend: ## Generate only frontend TypeScript code
	@echo "🔨 Generating frontend TypeScript code..."
	@mkdir -p web/src/gen/api
	@cd web && npx buf generate --template buf.gen.frontend.local.yaml ./api
	@echo "✅ Frontend code generated in web/src/gen/api"

.PHONY: run-back
run-back: ## Run application backend
	@go run ./cmd/main.go

.PHONY: tidy
tidy:
	@go mod tidy
	@if git rev-parse --is-inside-work-tree >/dev/null 2>&1; then \
	  git diff --exit-code go.mod go.sum || \
	  (echo "::error::go.mod or go.sum is out of sync" && exit 1); \
	else \
	  echo "⚠️  Skipping git diff — no .git directory"; \
	fi

.PHONY: build
build: ## Compile Go binary
	@echo "Building binary..."
	@go build -buildvcs=false -o bin/tax ./cmd
	@echo "✅ Binary built at bin/tax"

.PHONY: vet
vet: ## Run all linters
	@go vet ./...

.PHONY: gofmt
gofmt: ## Format code
	@gofmt -s -w . 

.PHONY: test-backend
test-backend: ## Run tests
	@go test -v ./...

.PHONY: check-generated
check-generated: ## Just check git diff
	@git diff --exit-code gen || (echo "❌ Generated files are not committed" && exit 1)

.PHONY: check-fmt
check-fmt: ## Check formatting (CI)
	@gofmt -l . | grep -q . && (echo "❌ Files need formatting (run 'make fmt')"; exit 1) || true

.PHONY: ci-backend
ci-backend: tidy check-fmt vet test-backend build ## Run all backend CI checks

# --- Frontend commands ----

.PHONY: frontend-install
frontend-install: ## Install frontend dependencies (clean CI install)
	@echo "📦 Installing frontend dependencies..."
	@(cd web && npm ci)
	@echo "✅ Frontend dependencies installed"

.PHONY: frontend-build
frontend-build: ## Build frontend (minimal check)
	@echo "Building frontend..."
	@(cd web && npm run build --silent) || echo "⚠️  Frontend build skipped"

.PHONY: run-front
run-front: ## Run application frontend
	@cd web && npm run dev

.PHONY: lint
lint: ## Run lint frontend
	@cd web && npm run lint

.PHONY: ci-frontend
ci-frontend: frontend-build ## Frontend CI: currently just build check

.PHONY: ci-all
ci-all: ci-backend frontend-build ## All CI: backend + frontend build only

# --- Docker commands ---
.PHONY: docker-up
docker-up: ## Start docker containers without build
	@docker compose up -d

.PHONY: docker-build
docker-build: ## Build and start docker containers
	@docker compose up --build

.PHONY: docker-down
docker-down: ## Down docker containers
	@docker compose down

.PHONY: local-CI
local-CI: ## Use act to check CI local
	@act -P ubuntu-latest=catthehacker/ubuntu:act-22.04


# --- Fullstek --- 

.PHONY: test-all
test-all: test-backend ## Run tests

.PHONY: run
run:
	@echo "🚀 Starting backend and frontend..."
	@go run ./cmd/main.go & \
	cd web && npm run dev

.PHONY: lint-all
lint-all: ## Run all linters
	@go vet ./...
	@cd web && npm run lint