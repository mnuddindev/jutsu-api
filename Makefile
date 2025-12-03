.PHONY: help setup install run dev build build-all test test-coverage test-integration benchmark lint fmt vet check ci clean docker-build docker-up docker-down docker-logs docker-clean swagger deps-update security-check mod-verify mod-graph todo fixme lines install-tools dev-setup build-prod redis-start redis-stop

# Variables
APP_NAME=jutsu-api
BINARY_NAME=bin/$(APP_NAME)
MAIN_PATH=./cmd/api
DOCKER_IMAGE=jutsu-api
DOCKER_TAG=latest

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
RED=\033[0;31m
NC=\033[0m # No Color

# Redis (from old Makefile)
REDIS_PORT=10120

help: ## Show this help message
	@echo '${CYAN}Available targets:${NC}'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  ${GREEN}%-20s${NC} %s\n", $$1, $$2}'

## 🚀 Setup & Execution

setup: ## Initial setup - copy env, install deps, generate swagger
	@echo "${CYAN}Setting up Jutsu API...${NC}"
	@if [ ! -f .env ]; then cp env.example .env; echo "${GREEN}✓ Created .env file${NC}"; fi
	@$(MAKE) install
	@$(MAKE) swagger
	@echo "${GREEN}✓ Setup complete! Run 'make run' to start the API${NC}"

install: ## Install dependencies
	@echo "${CYAN}Installing dependencies...${NC}"
	go mod download
	go mod tidy
	@echo "${GREEN}✓ Dependencies installed${NC}"

deps-update: ## Update dependencies
	@echo "${CYAN}Updating dependencies...${NC}"
	go get -u ./...
	go mod tidy
	@echo "${GREEN}✓ Dependencies updated${NC}"

run: redis-start ## Run the application (uses old run structure with redis-start, but uses go run if 'air' isn't available)
	@echo "${CYAN}Starting ${APP_NAME}...${NC}"
	@which air > /dev/null && air || (echo "${RED}Air not found, running without hot reload.${NC}"; go run $(MAIN_PATH)/main.go)

dev: ## Run with hot reload (requires air)
	@echo "${CYAN}Starting in development mode with hot reload...${NC}"
	@which air > /dev/null || (echo "${RED}Air not found. Install with: go install github.com/cosmtrek/air@latest${NC}" && exit 1)
	air

build: ## Build the binary
	@echo "${CYAN}Building ${APP_NAME}...${NC}"
	go build -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "${GREEN}✓ Built: $(BINARY_NAME)${NC}"

build-all: ## Build for multiple platforms
	@echo "${CYAN}Building for multiple platforms...${NC}"
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 go build -o bin/$(APP_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 go build -o bin/$(APP_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 go build -o bin/$(APP_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "${GREEN}✓ Built for multiple platforms in ./bin/${NC}"

build-prod: clean ## Production build (copied from old Makefile with path fix)
	@echo "${CYAN}Building for production...${NC}"
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "${GREEN}✓ Production build complete: $(BINARY_NAME)${NC}"

## ✅ Quality & Checks

test: ## Run tests
	@echo "${CYAN}Running tests...${NC}"
	go test -v -race ./...

test-coverage: ## Run tests with coverage
	@echo "${CYAN}Running tests with coverage...${NC}"
	go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "${GREEN}✓ Coverage report: coverage.html${NC}"

test-integration: ## Run integration tests
	@echo "${CYAN}Running integration tests...${NC}"
	go test -v -race -tags=integration ./...

benchmark: ## Run benchmarks
	@echo "${CYAN}Running benchmarks...${NC}"
	go test -bench=. -benchmem ./...

lint: ## Run linter
	@echo "${CYAN}Running linter...${NC}"
	@which golangci-lint > /dev/null || (echo "${RED}golangci-lint not found. Install from https://golangci-lint.run/usage/install/${NC}" && exit 1)
	golangci-lint run ./...
	@echo "${GREEN}✓ Linting complete${NC}"

fmt: ## Format code
	@echo "${CYAN}Formatting code...${NC}"
	go fmt ./...
	gofmt -s -w .
	@echo "${GREEN}✓ Code formatted${NC}"

vet: ## Run go vet
	@echo "${CYAN}Running go vet...${NC}"
	go vet ./...

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)
	@echo "${GREEN}✓ All checks passed${NC}"

ci: ## Run CI checks (for GitHub Actions)
	@echo "${CYAN}Running CI checks...${NC}"
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) build
	@echo "${GREEN}✓ CI checks passed${NC}"

security-check: ## Run security checks
	@echo "${CYAN}Running security checks...${NC}"
	@which gosec > /dev/null || (echo "${RED}gosec not found. Installing...${NC}" && go install github.com/securego/gosec/v2/cmd/gosec@latest)
	gosec ./...
	@echo "${GREEN}✓ Security check complete${NC}"

## 🛠️ Tooling & Utilities

install-tools: ## Install development tools
	@echo "${CYAN}Installing development tools...${NC}"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/cosmtrek/air@latest
	@go install github.com/securego/gosec/v2/cmd/gosec@latest
	@echo "${GREEN}✓ Development tools installed${NC}"

dev-setup: install install-tools ## Development environment setup
	@echo "${GREEN}✓ Development environment setup complete${NC}"

swagger: ## Generate Swagger documentation
	@echo "${CYAN}Generating Swagger docs...${NC}"
	@which swag > /dev/null || (echo "${RED}Swag not found. Installing...${NC}" && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g $(MAIN_PATH)/main.go -o ./docs
	@echo "${GREEN}✓ Swagger docs generated in ./docs${NC}"
	@echo "${CYAN}Access docs at: http://localhost:8080/docs${NC}"
	@echo "${CYAN}Swagger UI at: http://localhost:8080/swagger/index.html${NC}"

mod-verify: ## Verify dependencies
	@echo "${CYAN}Verifying dependencies...${NC}"
	go mod verify
	@echo "${GREEN}✓ Dependencies verified${NC}"

mod-graph: ## Show dependency graph
	go mod graph

todo: ## Show TODO comments in code
	@echo "${CYAN}TODO items:${NC}"
	@grep -r "TODO" --include="*.go" . || echo "No TODOs found"

fixme: ## Show FIXME comments in code
	@echo "${CYAN}FIXME items:${NC}"
	@grep -r "FIXME" --include="*.go" . || echo "No FIXMEs found"

lines: ## Count lines of code
	@echo "${CYAN}Lines of code:${NC}"
	@find . -name "*.go" -not -path "./vendor/*" -not -path "./docs/*" | xargs wc -l | tail -1

clean: ## Clean build artifacts
	@echo "${CYAN}Cleaning...${NC}"
	rm -f $(APP_NAME)
	rm -f $(BINARY_NAME)
	rm -rf bin/
	rm -f coverage.out coverage.html
	rm -rf tmp/
	@echo "${GREEN}✓ Cleaned${NC}"

## 🐳 Docker

docker-build: ## Build Docker image
	@echo "${CYAN}Building Docker image...${NC}"
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "${GREEN}✓ Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)${NC}"

docker-up: ## Start Docker containers
	@echo "${CYAN}Starting Docker containers...${NC}"
	docker-compose up -d
	@echo "${GREEN}✓ Containers started${NC}"

docker-down: ## Stop Docker containers
	@echo "${CYAN}Stopping Docker containers...${NC}"
	docker-compose down
	@echo "${GREEN}✓ Containers stopped${NC}"

docker-logs: ## View Docker logs
	docker-compose logs -f

docker-clean: ## Remove Docker containers and volumes
	@echo "${CYAN}Cleaning Docker containers and volumes...${NC}"
	docker-compose down -v
	@echo "${GREEN}✓ Docker cleaned${NC}"

## 💾 Redis (from old Makefile)

redis-start: ## Start Redis in background (only if not already running)
	@echo "${CYAN}Starting Redis on port $(REDIS_PORT)...${NC}"
	@redis-cli -p $(REDIS_PORT) ping > /dev/null 2>&1 || \
	(redis-server --port $(REDIS_PORT) --daemonize yes && echo "${GREEN}✓ Redis started${NC}") || \
	(echo "${RED}Error: redis-server command not found or failed to start.${NC}" && exit 1)

redis-stop: ## Stop Redis
	@echo "${CYAN}Stopping Redis on port $(REDIS_PORT)...${NC}"
	@redis-cli -p $(REDIS_PORT) shutdown || true
	@echo "${GREEN}✓ Redis stopped${NC}"