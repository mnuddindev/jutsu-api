.PHONY: build run test clean install swagger docker-up docker-down

# Application name
APP_NAME=jutsu-api
BINARY_NAME=bin/$(APP_NAME)

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(BINARY_NAME) ./cmd/api
	@echo "Build complete: $(BINARY_NAME)"

# Run the application
run:
	@echo "Running $(APP_NAME)..."
	@go run ./cmd/api

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies installed"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go -o docs
	@echo "Swagger documentation generated"

# Docker compose up
docker-up:
	@echo "Starting Docker containers..."
	@docker-compose up -d
	@echo "Docker containers started"

# Docker compose down
docker-down:
	@echo "Stopping Docker containers..."
	@docker-compose down
	@echo "Docker containers stopped"

# Lint the code
lint:
	@echo "Linting code..."
	@golangci-lint run ./...
	@echo "Linting complete"

# Format the code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Formatting complete"

# Vet the code
vet:
	@echo "Vetting code..."
	@go vet ./...
	@echo "Vetting complete"

# Install development tools
install-tools:
	@echo "Installing development tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development tools installed"

# Development setup
dev-setup: install install-tools
	@echo "Development environment setup complete"

# Production build
build-prod: clean
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o $(BINARY_NAME) ./cmd/api
	@echo "Production build complete: $(BINARY_NAME)"

