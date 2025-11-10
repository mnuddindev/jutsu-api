[![Jutsu](https://i.imgur.com/WlFJPl1.png)](#jutsu---a-anime-api-made-with-golang)

# 🌀 Jutsu - An Anime Streaming API built with GoLang

**Jutsu** is a high-performance backend API built with **Go (Golang)** for powering anime streaming platforms.  
It delivers blazing-fast performance, clean architecture, and flexible integrations — designed for developers who want full control over their anime data, streams, and metadata.

The word _Jutsu_ literally translates to _Art_ in Japanese (**術**). And that's what this API is. ;)

## ✨ Features

- ⚙️ RESTful endpoints for anime, episodes, and streaming sources
- ⚡ Fast response times powered by Go’s concurrency
- 🔐 Secure JWT authentication
- 💾 Support for caching (Redis or in-memory)
- 📦 Modular, clean, and scalable architecture
- 🔍 Search and discover anime by name, genre, or ID
- 📺 Stream and metadata delivery ready

---

## 🧱 Tech Stack

> _Constantly evolving — new modules will be added with every release, this is only inits._

| Component        | Technology                                 |
| ---------------- | ------------------------------------------ |
| **Language**     | Go (Golang)                                |
| **Framework**    | GoFiber v2                                 |
| **Cache**        | Redis                                      |
| **Config**       | Viper + godotenv                           |
| **Logger**       | Zap                                        |
| **Validation**   | go-playground/validator                    |
| **Docs**         | Swagger                                    |
| **Architecture** | Clean Architecture (Domain Driven Design)  |

## 🚀 Quick Start

### Prerequisites
- Go 1.25+
- Redis 6+ (optional)

### Installation

1. Clone the repository
2. Install dependencies: `make install`
3. Setup environment: `cp .env.example .env`
4. Start services: `make docker-up`
5. Run the application: `make run`

See [SETUP.md](SETUP.md) for detailed setup instructions.

## 📁 Project Structure

```
jutsu-api/
├── cmd/api/           # Application entry point
├── internal/
│   ├── config/        # Configuration
│   ├── infrastructure/ # Cache, logger
│   ├── interface/     # HTTP handlers, middleware
│   └── usecase/       # Business logic
├── pkg/utils/         # Utilities
└── docs/              # Swagger documentation
```

## ✨ Features

- ✅ Clean Architecture
- ✅ Configuration management
- ✅ Structured logging
- ✅ Redis caching
- ✅ Request validation
- ✅ CORS & Recovery middleware
- ✅ Swagger documentation
- ✅ Graceful shutdown
- ✅ Health checks
- ✅ Docker support
- ✅ Modular design

---

# DISCLAIMER

- `Jutsu` does not store any files , it only link to the media which is hosted on 3rd party services.
- `Jutsu` is explicitly made for educational purposes only and not for commercial usage. This repo will not be responsible for any misuse of it.
