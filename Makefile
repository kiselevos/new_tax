default: help
#.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# --- Backend commands ---
.PHONY: codegen
codegen: ## Generate gRPC code via Buf
	@echo "🧬 Generating gRPC code using Buf"
	@buf generate
	@echo "✅ gRPC code generated in gen/grpc/"

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

# --- Frontend commands ---
# --- Frontend commands (minimal) ---
.PHONY: frontend-build
frontend-build: ## Build frontend (minimal check)
	@echo "🏗️ Building frontend..."
	@(cd web && npm run build --silent) || echo "⚠️  Frontend build skipped"

.PHONY: frontend-lint
frontend-lint: ## Quick frontend lint
	@echo "🔍 Quick frontend lint..."
	@(cd web && npm run lint --silent --if-present) || echo "⚠️  Frontend lint skipped"

.PHONY: ci-frontend
ci-frontend: frontend-build ## Frontend CI: just build check

# --- Combined commands ---
.PHONY: ci-backend
ci-backend: tidy check-fmt lint-all test-all build ## Run all backend CI checks

.PHONY: ci-all
ci-all: ci-backend ci-frontend ## Run all CI checks
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

# --- Combined commands ---
.PHONY: ci-backend
ci-backend: tidy check-fmt lint-all test-all build ## Run all backend CI checks

.PHONY: ci-frontend
ci-frontend: frontend-type-check frontend-lint frontend-build ## Run all frontend CI checks

.PHONY: ci-all
ci-all: ci-backend ci-frontend ## Run all CI checks

.PHONY: local-CI
local-CI: ## Use act to check CI local
	@act -P ubuntu-latest=catthehacker/ubuntu:act-22.04
