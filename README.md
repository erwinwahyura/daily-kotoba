# Kotoba API

A REST API for Kotoba, a Japanese vocabulary learning app with JLPT-based progression and home screen widgets.

## Features

- User authentication with JWT
- JLPT N5-N1 vocabulary progression
- Daily word rotation with sequential learning
- Progress tracking and statistics
- Placement test for level assignment
- Skip/mark words as known functionality

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL 16
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt (cost 12)

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL 16 (via Docker)

## Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd kotoba-api
cp .env.example .env
```

### 2. Start PostgreSQL

```bash
docker-compose up -d
```

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run Database Migrations

```bash
# Migrations will be added in Phase 2
make migrate-up
```

### 5. Start the Server

```bash
make dev
# or
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### 6. Test the API

```bash
# Health check
curl http://localhost:8080/health

# API ping
curl http://localhost:8080/api/ping
```

## Environment Variables

See `.env.example` for all configuration options.

Key variables:
- `PORT`: Server port (default: 8080)
- `DB_HOST`: PostgreSQL host
- `DB_PASSWORD`: Database password (required)
- `JWT_SECRET`: Secret key for JWT tokens (required)

## Project Structure

```
kotoba-api/
├── cmd/
│   └── api/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/                  # Configuration management
│   ├── handlers/                # HTTP request handlers
│   ├── middleware/              # HTTP middleware
│   ├── models/                  # Data models
│   ├── repository/              # Database access layer
│   ├── services/                # Business logic
│   └── utils/                   # Utility functions
├── migrations/                  # Database migrations (coming in Phase 2)
├── docker-compose.yml           # Local PostgreSQL setup
├── .env.example                 # Example environment variables
└── README.md                    # This file
```

## Development

### Run Tests

```bash
make test
# or
go test ./...
```

### Build Binary

```bash
make build
# or
go build -o bin/kotoba-api cmd/api/main.go
```

### Docker Commands

```bash
make docker-up       # Start PostgreSQL
make docker-down     # Stop PostgreSQL
make docker-logs     # View PostgreSQL logs
```

### Available Make Commands

```bash
make help            # Show all available commands
make dev             # Run development server
make build           # Build binary
make test            # Run tests
make clean           # Clean build artifacts
```

## API Documentation

Comprehensive API documentation will be added in Phase 9.

## Deployment

Deployment guide for Hetzner VPS will be added in Phase 8.

## License

MIT
