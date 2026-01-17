.PHONY: help dev build test clean docker-up docker-down migrate-up migrate-down migrate-create

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

migrate-up: ## Run database migrations
	migrate -path migrations -database "postgresql://kotoba:kotoba_dev_password@localhost:5432/kotoba_db?sslmode=disable" up

migrate-down: ## Rollback database migrations
	migrate -path migrations -database "postgresql://kotoba:kotoba_dev_password@localhost:5432/kotoba_db?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create NAME=your_migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

.DEFAULT_GOAL := help
