#!/bin/bash

# Kotoba API - Deployment Script
# This script helps deploy the application to a production server

set -e

echo "🚀 Kotoba API Deployment Script"
echo "================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if .env.production exists
if [ ! -f .env.production ]; then
    echo -e "${RED}Error: .env.production file not found!${NC}"
    echo "Please create .env.production from .env.production.example"
    exit 1
fi

# Load environment variables
set -a
source .env.production
set +a

echo -e "${GREEN}✓${NC} Environment file loaded"

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker is not installed!${NC}"
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}Error: Docker Compose is not installed!${NC}"
    exit 1
fi

echo -e "${GREEN}✓${NC} Docker and Docker Compose detected"

# Build Docker image
echo ""
echo "📦 Building Docker image..."
docker build -t kotoba-api:latest .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Docker image built successfully"
else
    echo -e "${RED}✗${NC} Failed to build Docker image"
    exit 1
fi

# Stop existing containers
echo ""
echo "🛑 Stopping existing containers..."
docker-compose -f docker-compose.prod.yml down

# Start new containers
echo ""
echo "🚀 Starting containers..."
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Containers started successfully"
else
    echo -e "${RED}✗${NC} Failed to start containers"
    exit 1
fi

# Wait for services to be healthy
echo ""
echo "⏳ Waiting for services to be ready..."
sleep 10

# Check health
echo ""
echo "🏥 Checking API health..."
HEALTH_CHECK=$(curl -s http://localhost:${API_PORT:-8080}/health || echo "failed")

if [[ $HEALTH_CHECK == *"ok"* ]]; then
    echo -e "${GREEN}✓${NC} API is healthy!"
else
    echo -e "${RED}✗${NC} API health check failed"
    echo "Check logs with: docker-compose -f docker-compose.prod.yml logs api"
    exit 1
fi

echo ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Deployment completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo ""
echo "API is running at: http://localhost:${API_PORT:-8080}"
echo ""
echo "Useful commands:"
echo "  - View logs:    docker-compose -f docker-compose.prod.yml logs -f"
echo "  - Stop:         docker-compose -f docker-compose.prod.yml down"
echo "  - Restart:      docker-compose -f docker-compose.prod.yml restart"
echo ""
