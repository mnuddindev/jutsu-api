<div align="center">
[![Jutsu](https://i.imgur.com/WlFJPl1.png)](#jutsu---a-anime-api-made-with-golang)
# 🌀 Jutsu API

> **High-Performance Anime Streaming API built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge)](https://gofiber.io/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)

**A blazing-fast, production-ready RESTful API for anime streaming platforms**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [API Reference](#-api-reference) • [Contributing](#-contributing)

</div>

---

## ✨ Features

<div align="center">

|        🚀 Performance        |      🔒 Security      |    📦 Architecture    | 🛠️ Developer Experience |
| :--------------------------: | :-------------------: | :-------------------: | :---------------------: |
| ⚡ Ultra-fast response times | 🔐 JWT Authentication | 🏗️ Clean Architecture |  📝 Comprehensive Docs  |
|       💾 Redis Caching       |  🛡️ Input Validation  |   🧩 Modular Design   |      🧪 Unit Tests      |
|   🔄 Concurrent Processing   |  🔒 CORS Protection   | 📊 Structured Logging |    🐳 Docker Support    |

</div>

### 🎯 Core Capabilities

- **📺 Complete Anime Data** - Episodes, metadata, streaming sources, characters, voice actors
- **🔍 Advanced Search** - Filter by genre, type, status, rating, and more
- **📅 Schedule Management** - Daily schedules and next episode notifications
- **🎲 Random Discovery** - Get random anime recommendations
- **👥 User Features** - Watchlist support and personalized recommendations
- **⚡ High Performance** - Optimized caching with intelligent TTL management
- **📊 Real-time Updates** - Fresh data with efficient cache invalidation

---

## 🧱 Tech Stack

<div align="center">

|   Component    |                              Technology                               | Purpose                  |
| :------------: | :-------------------------------------------------------------------: | :----------------------- |
|  **Language**  |                    [Go 1.21+](https://golang.org/)                    | High-performance backend |
| **Framework**  |                    [Fiber v2](https://gofiber.io/)                    | Fast HTTP framework      |
|   **Cache**    |                      [Redis](https://redis.io/)                       | Distributed caching      |
|   **Logger**   |                 [Zap](https://github.com/uber-go/zap)                 | Structured logging       |
|   **Config**   |                [Viper](https://github.com/spf13/viper)                | Configuration management |
| **Validation** | [go-playground/validator](https://github.com/go-playground/validator) | Request validation       |
|    **Docs**    |                    [Swagger](https://swagger.io/)                     | API documentation        |

</div>

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21 or higher
- **Redis** 6.0+ (optional, for caching)
- **Docker** & **Docker Compose** (optional, for containerized setup)

### Installation

#### 1. Clone the Repository

```bash
git clone https://github.com/mnuddindev/jutsu-api.git
cd jutsu-api
```

#### 2. Install Dependencies

```bash
make install
# or
go mod download
```

#### 3. Configure Environment

```bash
cp env.example .env
# Edit .env with your configuration
```

#### 4. Start Services (Docker)

```bash
make docker-up
# or
docker-compose up -d
```

#### 5. Run the Application

```bash
make run
# or
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

> 📖 For detailed setup instructions, see [SETUP.md](SETUP.md)

---

## 📁 Project Structure

```
jutsu-api/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── config/           # Configuration management
│   ├── infrastructure/   # Cache, logger, database
│   ├── interface/        # HTTP handlers, middleware, routes
│   └── usecase/          # Business logic layer
├── pkg/
│   ├── extractors/       # Data extraction logic
│   ├── parsers/          # Data parsing & decryption
│   ├── helper/           # Helper functions & cache constants
│   ├── utils/            # Utility functions
│   ├── httpclient/       # HTTP client wrapper
│   └── scrape/           # Web scraping utilities
├── tests/                # Test files
├── docs/                 # Documentation
│   └── API.md           # Complete API documentation
├── docker-compose.yml    # Docker services configuration
├── Makefile             # Build & development commands
└── README.md            # This file
```

---

## 📚 Documentation

### API Reference

Complete API documentation is available in [docs/API.md](docs/API.md)

### Quick Examples

#### Get Home Page Data

```bash
curl http://localhost:8080/api
```

#### Get Anime Information

```bash
curl "http://localhost:8080/api/info?id=frieren-beyond-journeys-end-18542"
```

#### Search Anime

```bash
curl "http://localhost:8080/api/search?keyword=naruto&type=tv"
```

#### Get Streaming Info

```bash
curl "http://localhost:8080/api/stream?id=anime-id?ep=episode-id&server=hd-1&type=sub"
```

### Swagger Documentation

Once the server is running, visit:

```
http://localhost:8080/swagger/index.html
```

---

## 🎯 API Endpoints

### Core Endpoints

| Method |      Endpoint       | Description           |
| :----: | :-----------------: | :-------------------- |
| `GET`  |       `/api`        | Home page data        |
| `GET`  |     `/api/info`     | Anime details         |
| `GET`  | `/api/episodes/:id` | Episode list          |
| `GET`  |    `/api/stream`    | Streaming information |
| `GET`  |    `/api/search`    | Search anime          |
| `GET`  |    `/api/filter`    | Filter anime          |

### Category Endpoints

| Method |       Endpoint       | Description        |
| :----: | :------------------: | :----------------- |
| `GET`  | `/api/genre/{genre}` | Genre listings     |
| `GET`  |  `/api/top-airing`   | Top airing anime   |
| `GET`  | `/api/most-popular`  | Most popular anime |
| `GET`  |    `/api/top-ten`    | Top 10 anime       |

### Additional Endpoints

| Method |       Endpoint       | Description         |
| :----: | :------------------: | :------------------ |
| `GET`  |   `/api/schedule`    | Daily schedule      |
| `GET`  |    `/api/random`     | Random anime        |
| `GET`  | `/api/character/:id` | Character details   |
| `GET`  |  `/api/actors/:id`   | Voice actor details |
| `GET`  | `/api/producer/:id`  | Producer listings   |
| `GET`  |  `/api/studio/:id`   | Studio listings     |

> 📖 See [docs/API.md](docs/API.md) for complete endpoint documentation with examples

---

## ⚙️ Configuration

Configuration is managed through environment variables. Key settings:

```env
# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Logger
LOG_LEVEL=info
LOG_ENCODING=console
```

See `env.example` for all available options.

---

## 🧪 Testing

Run all tests:

```bash
make test
# or
go test ./...
```

Run tests with coverage:

```bash
make test-coverage
```

Run linting:

```bash
make lint
# or
golangci-lint run ./...
```

---

## 🐳 Docker

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

### Build Docker Image

```bash
docker build -t jutsu-api .
docker run -p 8080:8080 jutsu-api
```

---

## 📊 Performance

### Caching Strategy

The API implements intelligent caching with optimized TTLs:

|   Data Type   | Cache Duration | Reason             |
| :-----------: | :------------: | :----------------- |
|   Home Info   |   15 minutes   | Frequently updated |
|  Anime Info   |     1 hour     | Relatively stable  |
|  Categories   |   30 minutes   | Moderate updates   |
| Search/Filter |   5 minutes    | Dynamic results    |
|   Streaming   |   5 minutes    | Links may expire   |
|   Schedule    |   10 minutes   | Daily updates      |

### Benchmarks

- **Response Time**: < 50ms (cached)
- **Throughput**: 10,000+ req/s
- **Memory**: < 100MB (idle)

---

## 🔒 Security

- ✅ Input validation on all endpoints
- ✅ CORS protection
- ✅ Rate limiting (configurable)
- ✅ Secure error handling
- ✅ No sensitive data exposure

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and conventions
- Write tests for new features
- Update documentation
- Ensure all tests pass
- Run linter before submitting

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## ⚠️ Disclaimer

1. **Jutsu API** does not store any files. It only links to media hosted on 3rd party services.
2. This API is explicitly made for **educational purposes only** and not for commercial usage.
3. This repository will not be responsible for any misuse of it.

---

## 🙏 Acknowledgments

- Built with [Go Fiber](https://gofiber.io/)
- Powered by the anime community

---

<div align="center">

**Made with ❤️ using Go**

[⭐ Star this repo](https://github.com/mnuddindev/jutsu-api) • [🐛 Report Bug](https://github.com/mnuddindev/jutsu-api/issues) • [💡 Request Feature](https://github.com/mnuddindev/jutsu-api/issues)

</div>
