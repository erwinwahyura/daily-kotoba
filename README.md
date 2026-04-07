# Kotoba API 🍃

A REST API for Kotoba, a Japanese vocabulary learning app with JLPT-based progression, home screen widgets, and intelligent grammar pattern learning.

**Live API**: https://kotoba.erwarx.com

## Features

- 🔐 User authentication with JWT
- 📚 JLPT N5-N1 vocabulary progression with sample data
- 📖 Daily word rotation with sequential learning + review cycles
- 📖 N3-N1 grammar patterns with detailed pedagogy (examples, nuances, common mistakes)
- 📊 Progress tracking and statistics
- ✅ Placement test for level assignment
- ⏭️ Skip/mark words as known functionality
- 🔄 Automatic streak tracking
- 🎯 Intelligent level assignment based on test performance

## Tech Stack

- **Language**: Go 1.25+
- **Framework**: Gin
- **Database**: SQLite (single-node) / PostgreSQL (scalable)
- **Authentication**: JWT tokens
- **Password Hashing**: bcrypt (cost 12)
- **Deployment**: Docker + Hetzner VPS

## Quick Start

### Using the Live API

The API is deployed at **https://kotoba.erwarx.com**

```bash
# Health check
curl https://kotoba.erwarx.com/health

# Register
curl -X POST https://kotoba.erwarx.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"yourpassword","name":"Your Name"}'

# Login (returns JWT token)
curl -X POST https://kotoba.erwarx.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"your@email.com","password":"yourpassword"}'

# Get daily vocabulary (protected)
curl https://kotoba.erwarx.com/api/vocab/daily \
  -H "Authorization: Bearer <your-jwt-token>"

# Get daily grammar pattern (protected)
curl https://kotoba.erwarx.com/api/grammar/daily \
  -H "Authorization: Bearer <your-jwt-token>"
```

### Local Development

#### 1. Clone and Setup

```bash
git clone https://github.com/erwinwahyura/daily-kotoba.git
cd daily-kotoba
cp .env.example .env
```

#### 2. Start with Docker (recommended)

```bash
# Start with SQLite (no external DB needed)
make docker-build
make docker-up

# Or with PostgreSQL
docker-compose -f docker-compose.prod.yml up -d
```

#### 3. The server auto-migrates and auto-seeds on startup

Migrations run automatically. Sample data (vocabulary, grammar, placement questions) is loaded from `seeds/` directory on first start.

#### 4. API available at `http://localhost:8080`

## Environment Variables

See `.env.example` for all configuration options.

### Key Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENV` | Environment mode | `development` |
| `JWT_SECRET` | JWT signing secret | (required) |
| `DB_DRIVER` | `sqlite` or `postgres` | `postgres` |
| `SQLITE_PATH` | Path for SQLite file | `./kotoba.db` |
| `MIGRATIONS_DIR` | Migration files location | `./migrations` |
| `SEEDS_DIR` | Seed data files location | `./seeds` |

## API Documentation

Full API documentation is available in [API.md](./API.md).

### Authentication (All endpoints return JWT)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/auth/register` | No | Create new account |
| `POST` | `/api/auth/login` | No | Login, get token |
| `GET` | `/api/auth/me` | Yes | Get current user |

### Vocabulary Learning

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/vocab/daily` | Yes | Get today's word |
| `GET` | `/api/vocab/:id` | Yes | Get specific word |
| `POST` | `/api/vocab/:id/skip` | Yes | Skip/mark known |
| `GET` | `/api/vocab/level/:level` | Yes | Get words by JLPT level |

### Grammar Patterns (N3-N1)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/grammar/daily` | Yes | Get today's grammar pattern |
| `GET` | `/api/grammar/:id` | Yes | Get specific pattern |
| `GET` | `/api/grammar/level/:level` | Yes | Browse by level |

### Progress & Stats

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/progress` | Yes | Current progress |
| `GET` | `/api/progress/stats` | Yes | Detailed statistics |

### Placement Test

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `GET` | `/api/placement-test` | No | Get 20 test questions |
| `POST` | `/api/placement-test/submit` | Yes | Submit answers |
| `GET` | `/api/placement-test/result` | Yes | View result & assigned level |

### Test Account

```bash
# Login with test account
curl -X POST https://kotoba.erwarx.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"hiru@erwarx.com","password":"kotoba2024"}'
```

## Project Structure

```
daily-kotoba/
├── cmd/
│   ├── api/
│   │   └── main.go              # Application entry
│   └── seed/
│       ├── main.go              # Vocabulary seeding
│       ├── seed_placement.go    # Placement questions
│       └── seed_n3_grammar.go   # Grammar patterns
├── internal/
│   ├── config/                  # Configuration
│   ├── db/                      # Database + migrations + seeding
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Auth middleware
│   ├── models/                  # Data models
│   ├── repository/              # Data access layer
│   ├── services/                # Business logic
│   └── utils/                   # Utilities
├── migrations/                  # Database migrations
│   ├── *.up.sql                # PostgreSQL versions
│   └── *.sqlite.up.sql         # SQLite versions
├── seeds/                       # Seed data (JSON)
│   ├── 001_sample_n5_vocab.json
│   ├── 002_sample_n3_grammar.json
│   └── 003_sample_placement.json
├── docker-compose.hetzner.yml   # Production deployment
├── Dockerfile.hetzner           # Production build
└── README.md                    # This file
```

## Deployment 🚀

### Hetzner VPS (Current Production)

The API is deployed on Hetzner CX23 with Docker.

**Key Features:**
- **Auto-migrations**: Database schema updates on startup
- **Auto-seeding**: Sample data loaded from `seeds/` folder
- **SQLite with WAL mode**: Fast, single-node deployment
- **Persistent volume**: Data stored at `/mnt/apps-data`

#### Deploy Changes

```bash
ssh deploy@46.224.127.221
cd /opt/daily-kotoba
git pull origin main
docker-compose -f docker-compose.hetzner.yml down
docker-compose -f docker-compose.hetzner.yml up -d --build
```

### First-Time Setup

```bash
# On server:
mkdir -p /opt/daily-kotoba
cd /opt/daily-kotoba
git clone https://github.com/erwinwahyura/daily-kotoba.git .
cp .env.hetzner .env
# Edit .env with secure JWT_SECRET
docker-compose -f docker-compose.hetzner.yml up -d --build
```

### Local Production Test

```bash
# Build and run with production settings
make docker-build
docker-compose -f docker-compose.prod.yml up -d
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

### Makefile Commands

```bash
make help            # Show all commands
make docker-build    # Build Docker image
make docker-up       # Start containers
make docker-down     # Stop containers
make logs            # View container logs
```

## Project Architecture

### Auto-Migration System

On startup, the API automatically:
1. Creates `schema_migrations` tracking table
2. Runs pending `.up.sql` migrations (PostgreSQL or SQLite variants)
3. Skips already-applied migrations

### Auto-Seeding System

On startup, after migrations:
1. Creates `schema_seeds` tracking table  
2. Loads `.json` files from `seeds/` directory
3. Inserts data only once per seed file
4. Tracks applied seeds to avoid duplicates

**Adding Data:** Simply add new `.json` files to `seeds/` and redeploy.

### Database Support

| Feature | PostgreSQL | SQLite |
|---------|------------|--------|
| Auto-migrations | ✅ | ✅ |
| Auto-seeding | ✅ | ✅ |
| WAL mode | N/A | ✅ (performance) |
| Single-node | ❌ | ✅ |
| Production scaling | ✅ | ⚠️ (limited) |

**Current Production**: SQLite (simple, fast, single-node)

## License

MIT
