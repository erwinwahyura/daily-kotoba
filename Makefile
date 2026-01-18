.PHONY: help build run test clean docker-build docker-up docker-down migrate-up migrate-down seed-vocab seed-placement logs

# Default target
help:
	@echo "Available commands:"
	@echo "  make build           - Build the Go application"
	@echo "  make run             - Run the application locally"
	@echo "  make test            - Run tests"
	@echo "  make clean           - Clean build artifacts"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-up       - Start production containers"
	@echo "  make docker-down     - Stop production containers"
	@echo "  make logs            - View container logs"
	@echo ""
	@echo "Database commands:"
	@echo "  make migrate-up      - Run database migrations"
	@echo "  make migrate-down    - Rollback migrations"
	@echo "  make seed-vocab      - Seed vocabulary data"
	@echo "  make seed-placement  - Seed placement test questions"
	@echo ""
	@echo "Development:"
	@echo "  make dev-up          - Start development environment"
	@echo "  make dev-down        - Stop development environment"

# Build the application
build:
	@echo "Building application..."
	go build -o bin/kotoba-api cmd/api/main.go

# Run the application locally
run:
	@echo "Running application..."
	go run cmd/api/main.go

# Run tests
test:
	@echo "Running tests..."
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	rm -rf bin/
	rm -f main

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t kotoba-api:latest .

# Start production containers
docker-up:
	@echo "Starting production containers..."
	docker-compose -f docker-compose.prod.yml up -d
	@echo "Waiting for services to be ready..."
	sleep 10
	@echo "Services started! API available at http://localhost:8080"

# Stop production containers
docker-down:
	@echo "Stopping production containers..."
	docker-compose -f docker-compose.prod.yml down

# View logs
logs:
	docker-compose -f docker-compose.prod.yml logs -f

# Seed vocabulary data
seed-vocab:
	@echo "Seeding vocabulary data..."
	go run cmd/seed/main.go
	@echo "Vocabulary seeded!"

# Seed placement test questions
seed-placement:
	@echo "Seeding placement test questions..."
	go run cmd/seed/seed_placement.go
	@echo "Placement test questions seeded!"

# Development environment
dev-up:
	@echo "Starting development environment..."
	docker-compose up -d
	@echo "Development database started!"

dev-down:
	@echo "Stopping development environment..."
	docker-compose down
