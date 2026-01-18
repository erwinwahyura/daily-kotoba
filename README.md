# Kotoba API

A REST API for Kotoba, a Japanese vocabulary learning app with JLPT-based progression and home screen widgets.

## Features

- 🔐 User authentication with JWT
- 📚 JLPT N5-N1 vocabulary progression (50 N4 words currently)
- 📖 Daily word rotation with sequential learning
- 📊 Progress tracking and statistics
- ✅ Placement test for level assignment (20 questions)
- ⏭️ Skip/mark words as known functionality
- 🔄 Automatic streak tracking
- 🎯 Intelligent level assignment based on test performance

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
# Run all migrations
make migrate-up

# Rollback migrations (if needed)
make migrate-down

# Create a new migration
make migrate-create NAME=your_migration_name
```

### 5. Seed Database

```bash
# Seed N4 vocabulary (50 words)
go run cmd/seed/main.go

# Seed placement test questions (20 questions)
go run cmd/seed/seed_placement.go
```

### 6. Start the Server

```bash
make run
# or
go run cmd/api/main.go
```

The API will be available at `http://localhost:8080`

### 7. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Get placement test
curl http://localhost:8080/api/placement-test

# Register a user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
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

Full API documentation is available in [API.md](./API.md).

### Quick Reference

**Authentication:**
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Get current user (protected)

**Placement Test:**
- `GET /api/placement-test` - Get 20 test questions
- `POST /api/placement-test/submit` - Submit answers (protected)
- `GET /api/placement-test/result` - Get test result (protected)

**Vocabulary:**
- `GET /api/vocab/daily` - Get today's word (protected)
- `GET /api/vocab/:id` - Get specific word (protected)
- `POST /api/vocab/:id/skip` - Skip to next word (protected)
- `GET /api/vocab/level/:level` - Get words by level (protected)

**Progress:**
- `GET /api/progress` - Get user progress (protected)
- `GET /api/progress/stats` - Get detailed stats (protected)

## Deployment

### Production Deployment

Full deployment guide is available in [DEPLOYMENT.md](./DEPLOYMENT.md).

#### Quick Deploy with Docker

```bash
# 1. Create production environment file
cp .env.production.example .env.production
# Edit .env.production with your secure values

# 2. Build and deploy
./scripts/deploy.sh

# Or manually:
docker-compose -f docker-compose.prod.yml up -d --build
```

#### Useful Scripts

```bash
./scripts/deploy.sh      # Deploy to production
./scripts/migrate.sh     # Run database migrations
./scripts/backup.sh      # Backup database
```

### Makefile Commands

```bash
make help            # Show all commands
make docker-build    # Build production Docker image
make docker-up       # Start production containers
make docker-down     # Stop production containers
make logs            # View container logs
make seed-vocab      # Seed vocabulary data
make seed-placement  # Seed placement test questions
```

## License

MIT
