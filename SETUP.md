# Jutsu API - Setup Guide

## Prerequisites

- Go 1.25 or higher
- PostgreSQL 12 or higher
- Redis 6 or higher (optional but recommended)
- Make (optional, for using Makefile commands)

## Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd jutsu-api
   ```

2. **Install dependencies**
   ```bash
   make install
   # or
   go mod download
   go mod tidy
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start Docker services** (PostgreSQL and Redis)
   ```bash
   make docker-up
   # or
   docker-compose up -d
   ```

5. **Run database migrations** (when implemented)
   ```bash
   make migrate-up
   ```

6. **Build the application**
   ```bash
   make build
   # or
   go build -o bin/jutsu-api ./cmd/api
   ```

7. **Run the application**
   ```bash
   make run
   # or
   go run ./cmd/api
   ```

## Configuration

The application uses environment variables for configuration. See `.env.example` for all available options.

### Key Configuration Options

- `APP_ENV`: Application environment (development, production)
- `SERVER_PORT`: Server port (default: 8080)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `REDIS_HOST`: Redis host
- `REDIS_PORT`: Redis port

## Development

### Running in Development Mode

```bash
APP_ENV=development APP_DEBUG=true make run
```

### Generate Swagger Documentation

```bash
make swagger
# or
swag init -g cmd/api/main.go -o docs
```

### Running Tests

```bash
make test
# or
go test -v -race ./...
```

### Linting

```bash
make lint
# or
golangci-lint run ./...
```

## Project Structure

```
jutsu-api/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── infrastructure/      # External concerns
│   │   ├── database/        # Database setup
│   │   ├── cache/           # Redis cache
│   │   └── logger/          # Logger setup
│   ├── interface/           # HTTP layer
│   │   ├── http/
│   │   │   ├── handler/     # HTTP handlers
│   │   │   └── router/      # Route setup
│   │   ├── middleware/      # Middleware
│   │   └── validation/      # Validation
│   └── usecase/             # Business logic (to be implemented)
├── pkg/
│   └── utils/               # Utilities
├── docs/                    # Swagger documentation
├── docker-compose.yml       # Docker services
├── Makefile                 # Make commands
└── go.mod                   # Go dependencies
```

## API Endpoints

### Health Check
- `GET /health` - Health check endpoint
- `GET /ready` - Readiness check endpoint
- `GET /live` - Liveness check endpoint

### Swagger Documentation
- `GET /swagger/*` - Swagger UI

## Architecture

The project follows Clean Architecture principles:

- **Domain Layer**: Core business logic and entities
- **Usecase Layer**: Application-specific business rules
- **Interface Layer**: HTTP handlers, middleware, validation
- **Infrastructure Layer**: Database, cache, logger

## Features

- ✅ Clean Architecture
- ✅ Configuration management (Viper + godotenv)
- ✅ Structured logging (Zap)
- ✅ Database integration (GORM + PostgreSQL)
- ✅ Redis caching with helper functions
- ✅ Request validation
- ✅ CORS middleware
- ✅ Recovery middleware
- ✅ Request logging
- ✅ Swagger documentation
- ✅ Graceful shutdown
- ✅ Health checks
- ✅ Docker support
- ✅ Modular design (systems don't depend on each other)

## Performance Optimizations

- Prefork mode support
- Connection pooling (database and Redis)
- Prepared statements
- Response caching
- Efficient data structures

## Next Steps

1. Implement domain models
2. Implement use cases
3. Implement API endpoints
4. Add authentication/authorization
5. Add rate limiting
6. Add API versioning
7. Add monitoring and metrics

## License

MIT

