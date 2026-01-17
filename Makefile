.PHONY: help dev build test clean docker-up docker-down migrate-up migrate-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

dev: ## Run the development server
	go run cmd/api/main.go

build: ## Build the application binary
	go build -o bin/kotoba-api cmd/api/main.go

test: ## Run all tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	go clean

docker-up: ## Start PostgreSQL with Docker Compose
	docker-compose up -d

docker-down: ## Stop PostgreSQL container
	docker-compose down

docker-logs: ## View PostgreSQL logs
	docker-compose logs -f postgres

migrate-up: ## Run database migrations (to be implemented in Phase 2)
	@echo "Migrations will be added in Phase 2"

migrate-down: ## Rollback database migrations (to be implemented in Phase 2)
	@echo "Migrations will be added in Phase 2"

.DEFAULT_GOAL := help
