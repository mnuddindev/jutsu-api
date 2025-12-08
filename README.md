[![Jutsu](https://i.imgur.com/WlFJPl1.png)](#jutsu---a-anime-api-made-with-golang)

<div align="center">

# 🌀 Jutsu API

> **High-Performance Anime Streaming API built with Go**

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge)](https://gofiber.io/)
[![Redis](https://img.shields.io/badge/Redis-6+-DC382D?style=for-the-badge&logo=redis)](https://redis.io/)

**A blazing-fast, production-ready RESTful API for anime streaming platforms**

[Features](#-features) • [Quick Start](#-quick-start) • [Documentation](#-documentation) • [API Reference](#-documentation) • [Contributing](#-contributing)

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

Complete API documentation is available in [docs/API.md](docs/API.md)

## 📚 Documentation

The HTTP interface is fully described via Swagger:

- UI: `http://localhost:8080/docs`
- Files: `docs/swagger.json`, `docs/swagger.yaml`, `docs/docs.go`

### `GET` Home info

```bash
GET /api/
```

### Endpoint

```bash
/api/
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "spotlights": [
      {
        "id": "string",
        "data_id": "string",
        "title": "Frieren: Beyond Journey's End",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "description": "string",
        "tvInfo": {
          "showType": "string",
          "duration": "string",
          "releaseDate": "string",
          "quality": "string",
          "episodeInfo": [
            {}
          ]
        }
      }
    ],
    "trending": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "topTen": {
      "today": [
        {
          "id": "string",
          "data_id": "string",
          "number": 0,
          "name": "string",
          "poster": "https://example.com/poster.jpg",
          "tvInfo": {}
        }
      ],
      "week": [
        {
          "id": "string",
          "data_id": "string",
          "number": 0,
          "name": "string",
          "poster": "https://example.com/poster.jpg",
          "tvInfo": {}
        }
      ],
      "month": [
        {
          "id": "string",
          "data_id": "string",
          "number": 0,
          "name": "string",
          "poster": "https://example.com/poster.jpg",
          "tvInfo": {}
        }
      ]
    },
    "today": {
      "schedule": [
        {
          "id": "string",
          "data_id": "string",
          "title": "string",
          "japanese_title": "string",
          "releaseDate": "string",
          "time": "string",
          "episode_no": 0
        }
      ]
    },
    "genres": [
      "string"
    ]
  }
}
```

### `GET` Top 10 anime's info

```bash
GET /api/top-ten
```

### Endpoint

```bash
/api/top-ten
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/top-ten"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "today": [
      {
        "id": "string",
        "data_id": "string",
        "number": 0,
        "name": "string",
        "poster": "https://example.com/poster.jpg",
        "tvInfo": {}
      }
    ],
    "week": [
      {
        "id": "string",
        "data_id": "string",
        "number": 0,
        "name": "string",
        "poster": "https://example.com/poster.jpg",
        "tvInfo": {}
      }
    ],
    "month": [
      {
        "id": "string",
        "data_id": "string",
        "number": 0,
        "name": "string",
        "poster": "https://example.com/poster.jpg",
        "tvInfo": {}
      }
    ]
  }
}
```

### `GET` Top Search

```bash
GET /api/top-search
```

### Endpoint

```bash
/api/top-search
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/top-search"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": [
    {
      "title": "string",
      "link": "string"
    }
  ]
}
```

### `GET` Specified anime's info

```bash
GET /api/info
```

### Endpoint

```bash
/api/info?id={string}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
|   `id`    |    `query`     |  string   |  anime-id   |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/info?id=frieren-beyond-journeys-end-18542"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": {
      "id": "string",
      "title": "string",
      "japanese_title": "string",
      "poster": "https://example.com/poster.jpg",
      "description": "string",
      "stats": {},
      "genres": [
        "string"
      ]
    },
    "seasons": [
      {}
    ],
    "relatedAnime": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "recommended": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ]
  }
}
```

### `GET` Random anime's info

```bash
GET /api/random
```

### Endpoint

```bash
/api/random
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/random"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": {
      "id": "string",
      "title": "string",
      "japanese_title": "string",
      "poster": "https://example.com/poster.jpg",
      "description": "string",
      "stats": {},
      "genres": [
        "string"
      ]
    },
    "seasons": [
      {}
    ],
    "relatedAnime": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "recommended": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ]
  }
}
```

### `GET` Categories info

```bash
GET /api/{category}
```

### Endpoint

```bash
/api/{category}?page={number}
```

#### Parameters

| Parameter  | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :--------: | :------------: | :-------: | :---------: | :---------: | :-----: |
| `category` |     `path`     | `string`  | `Category`  |   Yes ✔️    |   --    |
|   `page`   |    `query`     | `number`  | `Page-no.`  |    No ❌    |   `1`   |

## Available Categories

### 📊 Status & Popularity Filters

- `top-airing`
- `most-popular`
- `most-favorite`
- `completed`
- `recently-updated`
- `recently-added`
- `top-upcoming`

### 🗣️ Language Filters

- `subbed-anime`
- `dubbed-anime`

### 🎭 Genre Filters

- `genre/action`
- `genre/adventure`
- `genre/cars`
- `genre/comedy`
- `genre/dementia`
- `genre/demons`
- `genre/drama`
- `genre/ecchi`
- `genre/fantasy`
- `genre/game`
- `genre/harem`
- `genre/historical`
- `genre/horror`
- `genre/isekai`
- `genre/josei`
- `genre/kids`
- `genre/magic`
- `genre/martial-arts`
- `genre/mecha`
- `genre/military`
- `genre/music`
- `genre/mystery`
- `genre/parody`
- `genre/police`
- `genre/psychological`
- `genre/romance`
- `genre/samurai`
- `genre/school`
- `genre/sci-fi`
- `genre/seinen`
- `genre/shoujo`
- `genre/shoujo-ai`
- `genre/shounen`
- `genre/shounen-ai`
- `genre/slice-of-life`
- `genre/space`
- `genre/sports`
- `genre/super-power`
- `genre/supernatural`
- `genre/thriller`
- `genre/vampire`

### 🔠 Alphabetical (A-Z) Listing

- `az-list`
- `az-list/other`
- `az-list/0-9`
- `az-list/a`
- `az-list/b`
- `az-list/c`
- `az-list/d`
- `az-list/e`
- `az-list/f`
- `az-list/g`
- `az-list/h`
- `az-list/i`
- `az-list/j`
- `az-list/k`
- `az-list/l`
- `az-list/m`
- `az-list/n`
- `az-list/o`
- `az-list/p`
- `az-list/q`
- `az-list/r`
- `az-list/s`
- `az-list/t`
- `az-list/u`
- `az-list/v`
- `az-list/w`
- `az-list/x`
- `az-list/y`
- `az-list/z`

### 🎞️ Format/Type Filters

- `movie`
- `special`
- `ova`
- `ona`
- `tv`
- `music`

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/most-popular?page=1"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "totalPages": 0
  }
}
```

### `GET` Anime of specific producers or studio

```bash
GET /api/producer/{id}
```

### Endpoint

```bash
/api/producer/{id}?page={number}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type |  Description  | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :-----------: | :---------: | :-----: |
|   `id`    |     `path`     | `string`  | `Producer ID` |   Yes ✔️    |   --    |
|  `page`   |    `query`     | `number`  |  `Page-no.`   |    No ❌    |   `1`   |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/producer/1?page=1"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "totalPages": 0
  }
}
```

### `GET` Search result's info

```bash
GET /api/search
```

### Endpoint

```bash
/api/search?keyword={string}&page={number}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
| `keyword` |    `query`     | `string`  |  `keyword`  |    No ❌    |   --    |
|  `page`   |    `query`     | `number`  | `Page-no.`  |    No ❌    |   `1`   |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/search?keyword=one%20punch%20man&page=1"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "totalPage": 0,
    "page": 0,
    "hasNext": true
  }
}
```

### `GET` Search suggestions

```bash
GET /api/search/suggest
```

### Endpoint

```bash
/api/search/suggest?keyword={string}
```

#### Parameters

| Parameter | Parameter-Type |   Type   | Description | Mandatory ? | Default |
| :-------: | :------------: | :------: | :---------: | :---------: | :-----: |
| `keyword` |    `query`     | `string` |  `keyword`  |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/search/suggest?keyword=demon"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": [
    {
      "title": "string",
      "link": "string"
    }
  ]
}
```

### `GET` Filter Anime

```bash
GET /api/filter
```

### Endpoint

```bash
/api/filter?keyword={string}&type={string}&status={string}&rated={string}&score={string}&season={string}&language={string}&genres={string}&sort={string}&sy={string}&sm={string}&sd={string}&ey={string}&em={string}&ed={string}&page={number}
```

#### Parameters

| Parameter  | Parameter-Type | Data-Type | Description                                                  | Mandatory ? |   Default   |
| :--------: | :------------: | :-------: | :----------------------------------------------------------- | :---------: | :---------: |
| `keyword`  |    `query`     |  string   | Search keyword                                               |    No ❌    | `undefined` |
|   `type`   |    `query`     |  string   | Type of anime (e.g., `tv`, `movie`, `ova`, `ona`, `special`) |    No ❌    |    `ALL`    |
|  `status`  |    `query`     |  string   | Status of anime (e.g., `ongoing`, `completed`)               |    No ❌    |    `ALL`    |
|  `rated`   |    `query`     |  string   | Rating of anime                                              |    No ❌    |    `ALL`    |
|  `score`   |    `query`     |  string   | Score rating                                                 |    No ❌    |    `ALL`    |
|  `season`  |    `query`     |  string   | Season of anime                                              |    No ❌    |    `ALL`    |
| `language` |    `query`     |  string   | Language of anime (e.g., `sub`, `dub`)                       |    No ❌    |    `ALL`    |
|  `genres`  |    `query`     |  string   | Comma-separated list of genre IDs                            |    No ❌    |    `ALL`    |
|   `sort`   |    `query`     |  string   | Sorting method                                               |    No ❌    |  `DEFAULT`  |
|    `sy`    |    `query`     |  string   | Start year                                                   |    No ❌    | `undefined` |
|    `sm`    |    `query`     |  string   | Start month                                                  |    No ❌    | `undefined` |
|    `sd`    |    `query`     |  string   | Start day                                                    |    No ❌    | `undefined` |
|    `ey`    |    `query`     |  string   | End year                                                     |    No ❌    | `undefined` |
|    `em`    |    `query`     |  string   | End month                                                    |    No ❌    | `undefined` |
|    `ed`    |    `query`     |  string   | End day                                                      |    No ❌    | `undefined` |
|   `page`   |    `query`     |  number   | Page number for pagination                                   |    No ❌    |     `1`     |

#### Example of Request

```bash
curl -X GET "http://localhost:8080/api/filter?type=tv&status=completed&rated=5&sort=default&page=1"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "data_id": "string",
        "title": "string",
        "japanese_title": "string",
        "poster": "https://example.com/poster.jpg",
        "duration": "string",
        "type": "string",
        "rating": "string",
        "episodes": {
          "sub": 0,
          "dub": 0
        }
      }
    ],
    "totalPage": 0,
    "page": 0,
    "hasNext": true
  }
}
```

### `GET` Anime's episode list

```bash
GET /api/episodes/{id}
```

### Endpoint

```bash
/api/episodes/{id}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
|   `id`    |     `path`     |  string   |  anime-id   |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/episodes/one-piece-100"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "episodes": [
      {
        "number": 0,
        "title": "string",
        "episodeId": "string",
        "isFiller": true
      }
    ]
  }
}
```

### `GET` Schedule of upcoming anime

```bash
GET /api/schedule
```

### Endpoint

```bash
/api/schedule?date={string}&tzOffset={number}
```

#### Parameters

| Parameter  | Parameter-Type | Data-Type |   Description   | Mandatory ? | Default |
| :--------: | :------------: | :-------: | :-------------: | :---------: | :-----: |
|   `date`   |     query      |  string   |      date       |    No ❌    |  today  |
| `tzOffset` |     query      |  number   | timezone offset |    No ❌    | `-330`  |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/schedule?date=2024-09-23"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": [
    {
      "id": "string",
      "data_id": "string",
      "title": "string",
      "japanese_title": "string",
      "releaseDate": "string",
      "time": "string",
      "episode_no": 0
    }
  ]
}
```

### `GET` Schedule of next episode of Anime

```bash
GET /api/schedule/{id}
```

### Endpoint

```bash
/api/schedule/{id}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
|   `id`    |     param      |  string   |  anime-id   |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/schedule/one-piece-100"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "nextEpisodeSchedule": {
      "id": "string",
      "data_id": "string",
      "title": "string",
      "japanese_title": "string",
      "releaseDate": "string",
      "time": "string",
      "episode_no": 0
    }
  }
}
```

### `GET` Qtip info

```bash
GET /api/qtip/{id}
```

### Endpoint

```bash
/api/qtip/{id}
```

#### Parameters

| Parameter | Data-Type | Description | Mandatory ? | Default |
| :-------: | :-------: | :---------: | :---------: | :-----: |
|   `id`    | `string`  |  anime-id   |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/qtip/frieren-beyond-journeys-end-18542"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {}
}
```

### `GET` Characters

```bash
GET /api/character/list/{id}
```

### Endpoint

```bash
/api/character/list/{id}?page={number}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
|   `id`    |     `path`     | `string`  |  anime-id   |   Yes ✔️    |   --    |
|  `page`   |    `query`     | `number`  | `Page-no.`  |    No ❌    |   `1`   |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/character/list/one-piece-100?page=1"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "currentPage": 0,
    "totalPages": 0,
    "data": [
      {
        "character": {
          "id": "string",
          "name": "string",
          "poster": "string",
          "cast": "string"
        },
        "voiceActors": [
          {
            "id": "string",
            "name": "string",
            "poster": "string"
          }
        ]
      }
    ]
  }
}
```

### `GET` Streaming info

```bash
GET /api/stream/{id}
```

### Endpoint

```bash
/api/stream/{id}?ep={string}&server={string}&type={string}
```

#### Parameters

| Parameter | Parameter-Type |   Type   | Mandatory ? | Default |
| :-------: | :------------: | :------: | :---------: | :-----: |
|   `id`    |     `path`     | `string` |   Yes ✔️    |   --    |
|   `ep`    |    `query`     | `string` |   Yes ✔️    |   --    |
| `server`  |    `query`     | `string` |    No ❌    | `hd-1`  |
|  `type`   |    `query`     | `string` |    No ❌    |  `sub`  |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/stream/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "streamingLink": [
      {
        "id": "string",
        "type": "sub",
        "link": {
          "file": "https://example.com/stream.m3u8",
          "type": "string"
        },
        "tracks": [
          {}
        ],
        "intro": {},
        "outro": {},
        "server": "string"
      }
    ],
    "servers": [
      {
        "type": "string",
        "data_id": "string",
        "server_id": "string",
        "serverName": "string"
      }
    ]
  }
}
```

### `GET` Fallback Streaming info

```bash
GET /api/stream/fallback/{id}
```

### Endpoint

```bash
/api/stream/fallback/{id}?ep={string}&server={string}&type={string}
```

#### Parameters

| Parameter | Parameter-Type |   Type   | Mandatory ? | Default |
| :-------: | :------------: | :------: | :---------: | :-----: |
|   `id`    |     `path`     | `string` |   Yes ✔️    |   --    |
|   `ep`    |    `query`     | `string` |   Yes ✔️    |   --    |
| `server`  |    `query`     | `string` |    No ❌    |   --    |
|  `type`   |    `query`     | `string` |    No ❌    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/stream/fallback/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "streamingLink": [
      {
        "id": "string",
        "type": "sub",
        "link": {
          "file": "https://example.com/stream.m3u8",
          "type": "string"
        },
        "tracks": [
          {}
        ],
        "intro": {},
        "outro": {},
        "server": "string"
      }
    ],
    "servers": [
      {
        "type": "string",
        "data_id": "string",
        "server_id": "string",
        "serverName": "string"
      }
    ]
  }
}
```

### `GET` Available servers of anime

```bash
GET /api/servers
```

### Endpoint

```bash
/api/servers?ep={string}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
|   `ep`    |    `query`     | `string`  | episode ID  |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/servers?ep=124260"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": [
    {
      "type": "string",
      "data_id": "string",
      "server_id": "string",
      "serverName": "string"
    }
  ]
}
```

### `GET` Character Details

```bash
GET /api/character/{id}
```

### Endpoint

```bash
/api/character/{id}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description  | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :----------: | :---------: | :-----: |
|   `id`    |     `path`     | `string`  | character-id |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/character/asta-340"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "name": "string",
        "japaneseName": "string",
        "profile": "string",
        "about": {
          "description": "string",
          "style": "string"
        },
        "voiceActors": [
          {}
        ],
        "animeography": [
          {}
        ]
      }
    ]
  }
}
```

### `GET` Voice Actor Details

```bash
GET /api/actors/{id}
```

### Endpoint

```bash
/api/actors/{id}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type |  Description   | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :------------: | :---------: | :-----: |
|   `id`    |     `path`     | `string`  | voice-actor-id |   Yes ✔️    |   --    |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/actors/gakuto-kajiwara-534"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "data": [
      {
        "id": "string",
        "name": "string",
        "japaneseName": "string",
        "profile": "string",
        "about": {
          "description": "string",
          "style": "string"
        },
        "roles": [
          {}
        ]
      }
    ]
  }
}
```

### `GET` User Watchlist

```bash
GET /api/watchlist/{userId}
```

### Endpoint

```bash
/api/watchlist/{userId}
/api/watchlist/{userId}/{page}
```

#### Parameters

| Parameter | Parameter-Type | Data-Type | Description | Mandatory ? | Default |
| :-------: | :------------: | :-------: | :---------: | :---------: | :-----: |
| `userId`  |     `path`     | `string`  |   User ID   |   Yes ✔️    |   --    |
|  `page`   |     `path`     | `integer` | Page number |    No ❌    |   `1`   |

#### Example of request

```bash
curl -X GET "http://localhost:8080/api/watchlist/user123"
```

#### Sample Response

```javascript
{
  "success": true,
  "results": {
    "totalPages": 0,
    "data": [
      {
        "id": "string",
        "title": "string",
        "poster": "string",
        "type": "string",
        "subCount": 0,
        "dubCount": 0
      }
    ]
  }
}
```

### `GET` Health Check

```bash
GET /health
```

### Endpoint

```bash
/health
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/health"
```

#### Sample Response

```javascript
{
  "status": "ok",
  "service": "Jutsu API",
  "version": "1.0.0",
  "uptime": "string",
  "timestamp": "string",
  "checks": {
    "cache": {
      "status": "string",
      "response_time_ms": 0
    },
    "external_sources": {
      "anime_provider": {}
    }
  }
}
```

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
