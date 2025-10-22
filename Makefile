# ============================================================
# 🧩 Version Variables
# ============================================================
PROTOC_GO_VERSION := v1.34.2
PROTOC_GO_GRPC_VERSION := v1.4.0
GOLANGCI_LINT_VERSION := v1.60.3

# ============================================================
# 🧠 Default & Help
# ============================================================
default: help

.PHONY: help
help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
	awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


# ============================================================
# ⚙️ Setup Commands
# ============================================================
.PHONY: setup
setup: ## Install all development tools locally
	@echo "📦 Installing development dependencies..."
	@echo "➡️  Installing Go protoc plugins..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@$(PROTOC_GO_VERSION)
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@$(PROTOC_GO_GRPC_VERSION)
	@echo "➡️  Installing linters..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "✅ All tools installed successfully!"


# ============================================================
# 🧬 Code Generation
# ============================================================
.PHONY: codegen
codegen: ## Generate backend Go gRPC code via protoc
	@echo "🧬 Generating Go gRPC code using protoc..."
	@mkdir -p gen/grpc/api
	@protoc \
	  --proto_path=./docs/grpc \
	  --go_out=gen/grpc \
	  --go_opt=paths=source_relative \
	  --go-grpc_out=gen/grpc \
	  --go-grpc_opt=paths=source_relative \
	  $$(find ./docs/grpc -name "*.proto")
	@mv gen/grpc/tax.pb.go gen/grpc/api/ 2>/dev/null || true
	@mv gen/grpc/tax_grpc.pb.go gen/grpc/api/ 2>/dev/null || true
	@echo "✅ Go gRPC code generated in gen/grpc/api"


# ============================================================
# 🧱 Build & Run
# ============================================================
.PHONY: run
run: ## Run application
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
test-all: ## Run tests
	@go test -v ./...

.PHONY: tidy
tidy: ## Check go.mod/go.sum
	@go mod tidy
	@git diff --exit-code || (echo "::error::go.mod or go.sum is out of sync" && exit 1)

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


# ============================================================
# 🔬 CI & Local Utilities
# ============================================================
.PHONY: local-CI
local-CI: ## Run GitHub Actions locally via act
	@act -P ubuntu-latest=catthehacker/ubuntu:act-22.04
