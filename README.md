# Weather API

A simple HTTP API that returns current weather for a fixed set of cities. Responses are cached in Redis.

Project idea: [Weather API (roadmap.sh)](https://roadmap.sh/projects/weather-api-wrapper-service)

## Prerequisites

- **Go** 1.21+ (or the version specified in `go.mod`)
- **Redis** (for response caching)

## Configuration

Edit `config.json` in the project root:

| Field         | Description                    | Example        |
|---------------|--------------------------------|----------------|
| `server_port` | HTTP listen address            | `":8080"`      |
| `redis_url`   | Redis server address           | `"localhost:6379"` |

Example:

```json
{
  "server_port": ":8080",
  "redis_url": "localhost:6379"
}
```

## Running the project

1. **Start Redis** (if not already running):

   ```bash
   redis-server
   ```

   Or with Docker:

   ```bash
   docker run -d -p 6379:6379 redis:alpine
   ```

2. **Install dependencies** (from the project root):

   ```bash
   go mod download
   ```

3. **Run the server**:

   ```bash
   go run ./cmd/server
   ```

   The server reads `config.json` from the current working directory, so run this from the project root.

   To run the built binary instead:

   ```bash
   go build -o weather-api ./cmd/server
   ./weather-api
   ```

## API

| Method | Path     | Description                    |
|--------|----------|--------------------------------|
| GET    | `/health` | Health check. Returns `200 OK`. |
| GET    | `/weather?city=<name>` | Current weather for the given city. |

**Supported cities:** London, Paris, Tokyo, New York, Sydney, Rome, Madrid, Berlin.

**Example:**

```bash
curl "http://localhost:8080/weather?city=London"
```

**Example response:**

```json
{
  "latitude": 51.5074,
  "longitude": -0.1278,
  "current_weather": {
    "temperature": 12.5,
    "windspeed": 5.2,
    "weathercode": 61,
    "time": "2025-02-12T12:00"
  }
}
```

## Project structure

```
weather-api/
├── cmd/server/          # Application entrypoint
│   ├── main.go
│   └── server.go
├── internal/
│   ├── config/          # Configuration loading
│   ├── handler/         # HTTP handlers
│   ├── weather/         # Weather client (Open-Meteo + Redis cache)
│   └── cache/           # Reserved for cache abstractions
├── config.json
├── go.mod
└── README.md
```
