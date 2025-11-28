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

## `GET` Home info

```bash
GET /api
```

### Endpoint

```bash
/api
```

> #### No parameter required ❌

#### Example of request

```bash
curl -X GET "http://localhost:8080/api"
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

## `GET` Top 10 anime's info

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

## `GET` Top Search

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

## `GET` Specified anime's info

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

## `GET` Random anime's info

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

## `GET` Categories info

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

#### Available Categories

- `/api/top-airing`
- `/api/most-popular`
- `/api/genre/{slug}` (e.g., action, adventure, comedy)
- `/api/producer/{id}`
- `/api/studio/{id}`

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

## `GET` Anime of specific producers or studio

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

## `GET` Search result's info

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

## `GET` Search suggestions

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

## `GET` Filter Anime

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

## `GET` Anime's episode list

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

## `GET` Schedule of upcoming anime

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

## `GET` Schedule of next episode of Anime

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

## `GET` Qtip info

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

## `GET` Characters

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

## `GET` Streaming info

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

## `GET` Fallback Streaming info

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

## `GET` Available servers of anime

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

## `GET` Character Details

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

## `GET` Voice Actor Details

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

## `GET` User Watchlist

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

## `GET` Health Check

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
