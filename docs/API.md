# 🎌 Jutsu API Documentation

> **Complete RESTful API documentation for Jutsu Anime API**

## 📋 Table of Contents

- [Base URL](#base-url)
- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Home Info](#home-info)
  - [Anime Information](#anime-information)
  - [Categories & Genres](#categories--genres)
  - [Search & Filter](#search--filter)
  - [Streaming](#streaming)
  - [Episodes](#episodes)
  - [Schedule](#schedule)
  - [Characters & Voice Actors](#characters--voice-actors)
  - [Random](#random)
  - [Watchlist](#watchlist)

---

## Base URL

```
http://localhost:8080/api
```

---

## Authentication

Currently, no authentication is required for public endpoints.

---

## Endpoints

### Home Info

#### `GET /api` or `GET /api/`

Get home page data including spotlights, trending, top 10, schedule, and category previews.

**Example Request:**
```bash
curl http://localhost:8080/api
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "spotlights": [...],
    "trending": [...],
    "topTen": {
      "today": [...],
      "week": [...],
      "month": [...]
    },
    "today": {
      "schedule": [...]
    },
    "topAiring": {...},
    "mostPopular": {...},
    "genres": ["action", "adventure", ...]
  }
}
```

---

### Anime Information

#### `GET /api/info`

Get detailed information about a specific anime including seasons.

**Query Parameters:**
- `id` (required) - Anime ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/info?id=frieren-beyond-journeys-end-18542"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "data": {
      "id": "frieren-beyond-journeys-end-18542",
      "title": "Frieren: Beyond Journey's End",
      "poster": "https://...",
      "description": "...",
      "genres": [...],
      "episodes": [...]
    },
    "seasons": [...]
  }
}
```

---

### Categories & Genres

#### `GET /api/{routeType}`

Get paginated anime listings by category or genre.

**Route Types:**
- Genres: `genre/action`, `genre/adventure`, `genre/comedy`, etc.
- Categories: `top-airing`, `most-popular`, `most-favorite`, `completed`, `recently-updated`, `top-upcoming`, `recently-added`
- Types: `movie`, `tv`, `ova`, `ona`, `special`
- A-Z Lists: `az-list`, `az-list/a`, `az-list/b`, etc.

**Query Parameters:**
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/genre/action?page=1"
curl "http://localhost:8080/api/top-airing?page=1"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "totalPages": 10,
    "data": [
      {
        "id": "anime-id",
        "title": "Anime Title",
        "poster": "https://...",
        "description": "...",
        "tvInfo": {...}
      }
    ]
  }
}
```

#### `GET /api/top-ten`

Get top 10 anime for today, week, and month.

**Example Request:**
```bash
curl http://localhost:8080/api/top-ten
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "today": [...],
    "week": [...],
    "month": [...]
  }
}
```

#### `GET /api/producer/:id`

Get anime by producer with pagination.

**Path Parameters:**
- `id` (required) - Producer ID or slug

**Query Parameters:**
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/producer/studio-pierrot?page=1"
```

#### `GET /api/studio/:id`

Get anime by studio with pagination.

**Path Parameters:**
- `id` (required) - Studio ID or slug

**Query Parameters:**
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/studio/studio-pierrot?page=1"
```

---

### Search & Filter

#### `GET /api/search`

Search anime with various filters.

**Query Parameters:**
- `keyword` (optional) - Search keyword
- `type` (optional) - Anime type (tv, movie, ova, etc.)
- `status` (optional) - Status (ongoing, completed, etc.)
- `rated` (optional) - Rating filter
- `score` (optional) - Minimum score
- `season` (optional) - Season filter
- `language` (optional) - Language (sub, dub)
- `genres` (optional) - Comma-separated genre IDs
- `sort` (optional) - Sort order
- `sy`, `sm`, `sd` (optional) - Start year, month, day
- `ey`, `em`, `ed` (optional) - End year, month, day
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/search?keyword=naruto&type=tv&page=1"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "data": [...],
    "totalPage": 5
  }
}
```

#### `GET /api/filter`

Filter anime with advanced criteria (same parameters as search).

**Example Request:**
```bash
curl "http://localhost:8080/api/filter?type=tv&status=ongoing&genres=1,2,3&page=1"
```

#### `GET /api/search/suggest`

Get search suggestions for a keyword.

**Query Parameters:**
- `keyword` (required) - Search keyword

**Example Request:**
```bash
curl "http://localhost:8080/api/search/suggest?keyword=naru"
```

**Example Response:**
```json
{
  "success": true,
  "results": [
    {
      "id": "naruto",
      "title": "Naruto",
      "poster": "https://..."
    }
  ]
}
```

#### `GET /api/top-search`

Get top search keywords.

**Example Request:**
```bash
curl http://localhost:8080/api/top-search
```

---

### Streaming

#### `GET /api/stream/{id}`

Get streaming information for an episode.

**Path Parameters:**
- `id` (required) – Anime slug (e.g., `frieren-beyond-journeys-end-18542`)

**Query Parameters:**
- `ep` (required) – Episode ID
- `server` (optional) – Server name (e.g., `hd-1`)
- `type` (optional) – Stream type (`sub` or `dub`)

**Example Request:**
```bash
curl "http://localhost:8080/api/stream/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "streamingLink": [
      {
        "id": 1,
        "type": "sub",
        "link": {
          "file": "https://...",
          "type": "hls"
        },
        "tracks": [...],
        "intro": {...},
        "outro": {...},
        "server": "hd-1"
      }
    ],
    "servers": [...]
  }
}
```

#### `GET /api/stream/fallback/{id}`

Get fallback streaming information (same parameters as `/api/stream/{id}`).

**Example Request:**
```bash
curl "http://localhost:8080/api/stream/fallback/frieren-beyond-journeys-end-18542?ep=107257&server=hd-1&type=sub"
```

#### `GET /api/servers`

Get available streaming servers for an episode.

**Query Parameters:**
- `ep` (required) – Episode ID

**Example Request:**
```bash
curl "http://localhost:8080/api/servers?ep=124260"
```

**Example Response:**
```json
{
  "success": true,
  "results": [
    {
      "type": "sub",
      "data_id": 123,
      "server_id": 456,
      "serverName": "Vidcloud"
    }
  ]
}
```

---

### Episodes

#### `GET /api/episodes/:id`

Get list of episodes for an anime.

**Path Parameters:**
- `id` (required) - Anime ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/episodes/frieren-beyond-journeys-end-18542"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "episodes": [
      {
        "id": "episode-id",
        "number": 1,
        "title": "Episode 1",
        "url": "https://..."
      }
    ]
  }
}
```

---

### Schedule

#### `GET /api/schedule`

Get anime schedule for a specific date.

**Query Parameters:**
- `date` (optional) - Date in YYYY-MM-DD format (default: today)
- `tzOffset` (optional) - Timezone offset in minutes (default: -330)

**Example Request:**
```bash
curl "http://localhost:8080/api/schedule?date=2025-01-18"
```

**Example Response:**
```json
{
  "success": true,
  "results": [
    {
      "id": "anime-id",
      "title": "Anime Title",
      "episode": "Episode 1",
      "time": "12:00",
      "poster": "https://..."
    }
  ]
}
```

#### `GET /api/schedule/:id`

Get next episode schedule for a specific anime.

**Path Parameters:**
- `id` (required) - Anime ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/schedule/frieren-beyond-journeys-end-18542"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "nextEpisodeSchedule": "2025-01-25 12:00"
  }
}
```

---

### Characters & Voice Actors

#### `GET /api/character/list/:id`

Get paginated list of characters with voice actors for an anime.

**Path Parameters:**
- `id` (required) - Anime ID or slug

**Query Parameters:**
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/character/list/frieren-beyond-journeys-end-18542?page=1"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "currentPage": 1,
    "totalPages": 5,
    "data": [
      {
        "character": {
          "id": "character-id",
          "poster": "https://...",
          "name": "Character Name",
          "cast": "Main"
        },
        "voiceActors": [...]
      }
    ]
  }
}
```

#### `GET /api/character/:id`

Get detailed information about a character.

**Path Parameters:**
- `id` (required) - Character ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/character/asta-340"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "data": [{
      "id": "asta-340",
      "name": "Asta",
      "profile": "https://...",
      "japaneseName": "アスタ",
      "about": {
        "description": "...",
        "style": "<p>...</p>"
      },
      "voiceActors": [...],
      "animeography": [...]
    }]
  }
}
```

#### `GET /api/actors/:id`

Get detailed information about a voice actor.

**Path Parameters:**
- `id` (required) - Voice actor ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/actors/gakuto-kajiwara-534"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "data": [{
      "id": "gakuto-kajiwara-534",
      "name": "Kajiwara, Gakuto",
      "profile": "https://...",
      "japaneseName": "梶原岳人",
      "about": {...},
      "roles": [...]
    }]
  }
}
```

---

### Random

#### `GET /api/random`

Get a random anime with full details.

**Example Request:**
```bash
curl http://localhost:8080/api/random
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "id": "random-anime-id",
    "title": "Random Anime",
    "poster": "https://...",
    ...
  }
}
```

#### `GET /api/random/id`

Get only a random anime ID.

**Example Request:**
```bash
curl http://localhost:8080/api/random/id
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "id": "random-anime-id"
  }
}
```

---

### Watchlist

#### `GET /api/watchlist/:userId` or `GET /api/watchlist/:userId/:page`

Get user's watchlist with pagination.

**Path Parameters:**
- `userId` (required) - User ID
- `page` (optional) - Page number (default: 1)

**Example Request:**
```bash
curl "http://localhost:8080/api/watchlist/user123?page=1"
curl "http://localhost:8080/api/watchlist/user123/1"
```

**Example Response:**
```json
{
  "success": true,
  "results": {
    "totalPages": 3,
    "data": [
      {
        "id": "anime-id",
        "poster": "https://...",
        "title": "Anime Title",
        "tvInfo": {...}
      }
    ]
  }
}
```

---

### Additional Endpoints

#### `GET /api/qtip/:id`

Get qtip (quick tip) information for an anime.

**Path Parameters:**
- `id` (required) - Anime ID or slug

**Example Request:**
```bash
curl "http://localhost:8080/api/qtip/frieren-beyond-journeys-end-18542"
```

---

### Health & Platform

#### `GET /health`

Comprehensive health report including cache and upstream dependency checks.

**Example Response:**
```json
{
  "status": "ok",
  "service": "Jutsu API",
  "version": "1.0.0",
  "uptime": "72h34m",
  "checks": {
    "cache": {
      "status": "ok",
      "response_time_ms": 2
    },
    "external_sources": {
      "anime_provider": {
        "status": "ok"
      }
    }
  }
}
```

#### `GET /ready`

Returns readiness state used by orchestrators / load balancers.

**Example Response:**
```json
{
  "status": "ready",
  "timestamp": "2025-01-18T12:00:00Z"
}
```

#### `GET /live`

Simple liveness probe to confirm the process is running.

**Example Response:**
```json
{
  "status": "alive",
  "timestamp": "2025-01-18T12:00:00Z"
}
```

---

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "message": "Error message here"
}
```

**HTTP Status Codes:**
- `200` - Success
- `400` - Bad Request (missing/invalid parameters)
- `404` - Not Found
- `500` - Internal Server Error
- `502` - Bad Gateway (upstream service error)
- `503` - Service Unavailable

---

## Rate Limiting

Currently, no rate limiting is implemented. Please use the API responsibly.

---

## Caching

The API uses Redis caching for improved performance. Cache TTLs are optimized based on data volatility:

- **Home Info**: 15 minutes
- **Anime Info**: 1 hour
- **Categories**: 30 minutes
- **Search/Filter**: 5 minutes
- **Streaming Info**: 5 minutes
- **Schedule**: 10 minutes

---

## Support

For issues, questions, or contributions, please open an issue on GitHub.

---

**Made with ❤️ using Go**

