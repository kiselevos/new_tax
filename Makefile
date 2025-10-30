# ============================================================
# 🧩 Version Variables
# ============================================================
PROTOC_GO_VERSION := v1.34.2
PROTOC_GO_GRPC_VERSION := v1.4.0
GOLANGCI_LINT_VERSION := v1.60.3

TOOLS_DIR := $(CURDIR)/.tools/bin
export PATH := $(TOOLS_DIR):$(PATH)

# ============================================================
# 🧠 Default & Help
# ============================================================
default: help

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# ============================================================
# ⚙️ Setup Tools
# ============================================================
.PHONY: setup-tools
setup-tools: setup-lint-tools setup-proto-tools ## Install all dev tools (linters + proto)
	@echo "✅ All tools installed in $(TOOLS_DIR)"

.PHONY: setup-lint-tools
setup-lint-tools: ## Install only linting tools
	@echo "🔍 Installing lint tools..."
	@mkdir -p $(TOOLS_DIR)
	@GOBIN=$(TOOLS_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "✅ Lint tools installed in $(TOOLS_DIR)"

.PHONY: setup-proto-tools
setup-proto-tools: ## Install proto/gRPC tools for code generation
	@echo "🧬 Installing proto tools..."
	@mkdir -p $(TOOLS_DIR)
	@GOBIN=$(TOOLS_DIR) go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)
	@GOBIN=$(TOOLS_DIR) go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GO_GRPC_VERSION)
	@echo "✅ Proto tools installed in $(TOOLS_DIR)"

# ============================================================
# ⚙️ Setup Commands
# ============================================================
.PHONY: setup-backend
setup-backend: ## Устанавливает зависимости для backend (CLI-инструменты фиксируются через tools.go)
	@echo "⚙️ Installing Backend..."
	@go mod tidy
	@go mod download	
	@echo "✅ Backend install"

.PHONY: setup
setup: setup-backend
	@cd web && make setup-frontend
	@echo "✅ Backend & Frontend were installed successfully"


# ============================================================
# 🧬 Code Generation
# ============================================================
.PHONY: codegen
codegen: ## Generate backend Go gRPC code via protoc
	@echo "🧬 Generating Go gRPC code..."
	@mkdir -p gen/grpc/api
	@protoc \
	  --proto_path=./docs/grpc \
	  --plugin=protoc-gen-go=$(shell which protoc-gen-go) \
	  --plugin=protoc-gen-go-grpc=$(shell which protoc-gen-go-grpc) \
	  --go_out=gen/grpc --go_opt=paths=source_relative \
	  --go-grpc_out=gen/grpc --go-grpc_opt=paths=source_relative \
	  $(shell find ./docs/grpc -name "*.proto" -type f)
	@mv gen/grpc/tax*.pb.go gen/grpc/api/ 2>/dev/null || true
	@echo "✅ gRPC code generated → gen/grpc/api"


# ============================================================
# 🧱 Build & Run
# ============================================================
.PHONY: run
run: ## Run backend application
	@go run ./cmd/main.go

.PHONY: build
build: ## Compile Go binary
	@echo "🏗️  Building binary..."
	@go build -o bin/tax ./cmd
	@echo "✅ Binary built at bin/tax"


# ============================================================
# 🧪 Tests & Checks
# ============================================================
.PHONY: test-all
test-all: ## Run all tests
	@go test -v ./...

.PHONY: tidy
tidy: ## Check go.mod/go.sum
	@go mod tidy
	@git diff --exit-code || (echo "::error::go.mod or go.sum is out of sync" && exit 1)
	@go mod download

.PHONY: lint-all
lint-all: ## Run all linters
	@go vet ./...
	@golangci-lint run ./...

.PHONY: gofmt
gofmt: ## Format code
	@gofmt -s -w .

.PHONY: check-generated
check-generated: ## Check git diff for generated files
	@git diff --exit-code gen || (echo "❌ Generated files are not committed" && exit 1)

.PHONY: check-fmt
check-fmt: ## Check formatting (CI)
	@gofmt -l . | grep -q . && (echo "❌ Files need formatting (run 'make gofmt')"; exit 1) || true


# ============================================================
# 🐳 Docker Commands
# ============================================================
.PHONY: docker-up
docker-up: ## Start docker containers without build
	@docker compose up -d

.PHONY: docker-build
docker-build: ## Build and start docker containers
	@docker compose up --build

.PHONY: docker-down
docker-down: ## Stop and remove docker containers
	@docker compose down

.PHONY: docker-prune
docker-prune: ## Cleane docker cash
	@docker builder prune -f

# ============================================================
# 🔬 CI & Local Utilities
# ============================================================
.PHONY: local-CI
local-CI: ## Run GitHub Actions locally via act
	@act -P ubuntu-latest=catthehacker/ubuntu:act-22.04


# ============================================================
# 🔬 Work with frontend
# ============================================================
.PHONY: run-all
run-all: ## Run backend and frontend together
	@echo "🚀 Starting backend and frontend..."
	@trap 'kill 0' EXIT; \
	go run -tags=tools ./cmd/main.go & \
	cd web && go run ./cmd/web.go
