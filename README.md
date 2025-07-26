# Banner API

Simple banner click tracking API built with Go and FastHTTP.

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Clone repository
git clone https://github.com/acom21/banner-api.git
cd banner-api

# Start services
docker-compose up -d

# Test API
curl http://localhost:8080/counter/1
curl -X POST http://localhost:8080/stats/1 \
  -H "Content-Type: application/json" \
  -d '{"from": "2025-01-01T00:00:00", "to": "2025-12-31T23:59:59"}'
```

### Manual Setup

```bash
# 1. Start PostgreSQL
docker run --name postgres \
  -e POSTGRES_DB=banner_db \
  -e POSTGRES_USER=user \
  -e POSTGRES_PASSWORD=pass \
  -p 5432:5432 -d postgres:15-alpine

# 2. Build and run
go build -o banner-api ./cmd
./banner-api
```

## API Endpoints

- **GET** `/counter/{bannerID}` - Register banner click
- **POST** `/stats/{bannerID}` - Get click statistics

### Example Usage

```bash
# Register clicks
curl http://localhost:8080/counter/1
curl http://localhost:8080/counter/2

# Get statistics
curl -X POST http://localhost:8080/stats/1 \
  -H "Content-Type: application/json" \
  -d '{"from": "2025-01-01T00:00:00", "to": "2025-12-31T23:59:59"}'

# Response:
# {"stats":[{"ts":"2025-07-28T02:14:00","v":3}]}
```

## Configuration

Environment variables:
- `HOST` - Server host (default: 0.0.0.0)
- `PORT` - Server port (default: 8080)
- `DB_URL` - PostgreSQL connection URL

See `config.example.env` for example configuration.

## Documentation

📖 **[API Documentation (Swagger)](https://github.com/acom21/banner-api/blob/main/docs/swagger.yaml)**

🔗 **[Swagger JSON](https://github.com/acom21/banner-api/blob/main/docs/swagger.json)**
