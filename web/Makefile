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
.PHONY: tidy
tidy: ## Check go.mod/go.sum
	@go mod tidy
	@git diff --exit-code || (echo "::error::go.mod or go.sum is out of sync" && exit 1)
	

.PHONY: lint-all
lint-all: ## Run linter
	@go vet ./...
	@golangci-lint run ./...
	@echo "✅ Linters check"

.PHONY: test-all
test-all: ## Run tests
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
