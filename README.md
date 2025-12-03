[![Jutsu](https://i.imgur.com/WlFJPl1.png)](#jutsu---a-anime-api-made-with-golang)

<div align="center">

# 🌀 Jutsu API

> **High-Performance Anime Streaming API built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge)](https://gofiber.io/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)

**A blazing-fast, production-ready RESTful API for anime streaming platforms**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [API Reference](#-api-reference) • [Contributing](#-contributing)

</div>

## ⚠️ Disclaimer

1. This API does **not** store any media files; it only links to content hosted by 3rd‑party services.
2. This API is built **for educational purposes only** and not for commercial use.
3. The maintainers of this repository are **not responsible** for any misuse.

---

## ✨ Features

- **Home info** – spotlights, trending lists, top‑ten, schedule and category previews
- **Anime info** – seasons, related shows, recommendations, characters and sidebar “Qtip” data
- **Episodes & streaming** – episode lists, available servers, primary + fallback streaming info
- **Search & filter** – keyword search, advanced filter options, suggestions and top search keywords
- **People & studios** – character details, voice actors, producer and studio catalogues
- **Schedule & random** – daily schedule plus next‑episode info and random anime endpoints
- **Caching & observability** – per‑route TTLs backed by Redis (optional) and Zap structured logs

All public HTTP routes are documented via Swagger and kept in sync with the codebase.

---

## 🧱 Tech Stack

|   Component   | Technology | Purpose                                       |
| :-----------: | ---------- | --------------------------------------------- |
| **Language**  | Go 1.25    | High‑performance backend runtime              |
| **Framework** | Fiber v2   | Fast HTTP framework with minimal overhead     |
|   **Cache**   | Redis      | Optional response caching per endpoint        |
|  **Logger**   | Zap        | Structured JSON/console logging               |
|  **Config**   | Viper      | Configuration and `.env` loading              |
|   **Docs**    | Swag + UI  | Swagger/OpenAPI generation and UI integration |

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

## 📚 Documentation

### API Reference

Complete API documentation is available in http://localhost:8080/docs

## 📚 Documentation

The HTTP interface is fully described via Swagger:

- UI: `http://localhost:8080/swagger/index.html`
- Files: `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go`

### GET Home info

```bash
curl "http://localhost:8080/api/"
```

### GET Top 10 anime's info

```bash
curl "http://localhost:8080/api/top-ten"
```

### GET Top Search

```bash
curl "http://localhost:8080/api/top-search"
```

### GET Specified anime's info

```bash
curl "http://localhost:8080/api/info?id=frieren-beyond-journeys-end-18542"
```

### GET Random anime's info

```bash
curl "http://localhost:8080/api/random"
```

### GET Categories info

```bash
curl "http://localhost:8080/api/genre/action?page=1"
curl "http://localhost:8080/api/top-airing?page=1"
```

### GET Producer's & studio's anime

```bash
curl "http://localhost:8080/api/producer/1?page=1"
curl "http://localhost:8080/api/studio/1?page=1"
```

### GET Search results info

```bash
curl "http://localhost:8080/api/search?keyword=naruto&type=tv&page=1"
```

### GET Search suggestions

```bash
curl "http://localhost:8080/api/search/suggest?keyword=nar"
```

### GET Anime schedule

```bash
curl "http://localhost:8080/api/schedule?date=2025-01-18&tzOffset=-330"
```

### GET Anime's next episode schedule

```bash
curl "http://localhost:8080/api/schedule/frieren-beyond-journeys-end-18542"
```

### GET Anime Qtip info

```bash
curl "http://localhost:8080/api/qtip/frieren-beyond-journeys-end-18542"
```

### GET Anime characters

```bash
curl "http://localhost:8080/api/character/list/frieren-beyond-journeys-end-18542?page=1"
```

### GET Character details

```bash
curl "http://localhost:8080/api/character/asta-340"
```

### GET Voice actor details

```bash
curl "http://localhost:8080/api/actors/gakuto-kajiwara-534"
```

### GET Anime episodes

```bash
curl "http://localhost:8080/api/episodes/frieren-beyond-journeys-end-18542"
```

### GET Anime episode's available servers

```bash
curl "http://localhost:8080/api/servers?ep=107257"
```

### GET Anime stream info

```bash
curl "http://localhost:8080/api/stream/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

### GET Fallback streaming info

```bash
curl "http://localhost:8080/api/stream/fallback/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

### GET User watchlist

```bash
curl "http://localhost:8080/api/watchlist/user123"
curl "http://localhost:8080/api/watchlist/user123/2"
```

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

## 📄 License

This project is licensed under the **MIT License** – see `LICENSE` for details.

---

## 🙏 Acknowledgements

- Powered by the anime community
- Built with [Go](https://go.dev/) and [Fiber](https://gofiber.io/)

**Made with Go and Fiber.** If you find this useful, consider starring the repo. ⭐
