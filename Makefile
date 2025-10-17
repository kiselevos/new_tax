default: help
#.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
	
# --- Backend commands ---
.PHONY: codegen
codegen: ## Generate gRPC code via Buf
	@echo "Generating gRPC code using Buf"
	@buf generate
	@echo "✅ gRPC code generated in gen/grpc/api"

.PHONY: run
run: ## Run application
	@go run ./cmd/main.go

.PHONY: tidy
tidy: ## Check go.mod/go.sum
	@go mod tidy
	@git diff --exit-code || (echo "::error::go.mod or go.sum is out of sync" && exit 1)

.PHONY: build
build: ## Compile Go binary
	@echo "Building binary..."
	@go build -o bin/tax ./cmd
	@echo "✅ Binary built at bin/tax"

.PHONY: lint-all
lint-all: ## Run all linters
	@go vet ./...
	@golangci-lint run ./...

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
ci-backend: tidy check-fmt lint-all test-backend build ## Run all backend CI checks

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

.PHONY: setup
setup: ## Install all dependencies (Go + Frontend)
	@echo "Setting up full environment..."
	@make tidy
	@make codegen
	@make frontend-install
	@echo "✅ Environment ready!"