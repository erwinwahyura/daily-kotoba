#!/bin/bash
# Deploy script for Kotoba Kanji backend
# Run on Hetzner VPS as deploy user

set -e

echo "🐹 Deploying Kotoba Kanji Backend..."

# Navigate to project
cd ~/daily-kotoba || exit 1

# Pull latest changes (includes 11 kanji seed)
git pull origin main

# Build binary
echo "Building Go binary..."
go build -o kotoba-api ./cmd/api/

# Check if build succeeded
if [ ! -f "./kotoba-api" ]; then
    echo "❌ Build failed!"
    exit 1
fi

echo "✅ Build successful"

# Stop existing service (if running)
echo "Stopping existing service..."
pkill -f "./kotoba-api" || true

# Start new service
echo "Starting Kotoba API..."
export PORT=8090
export DB_DRIVER=sqlite
export DB_CONNECTION="./kotoba.db"
export JWT_SECRET="${JWT_SECRET:-your-jwt-secret-here}"
export JWT_EXPIRATION_HOURS=24

nohup ./kotoba-api > kotoba.log 2>&1 &

# Wait for startup
sleep 3

# Check health
echo "Checking health..."
if curl -s http://localhost:8090/health | grep -q "ok"; then
    echo "✅ API is running!"
    echo ""
    echo "🎉 DEPLOY SUCCESSFUL"
    echo "API URL: http://localhost:8090"
    echo "Health: http://localhost:8090/health"
    echo ""
    echo "Next: Seed kanji data"
    echo "curl -X POST http://localhost:8090/api/kanji/seed \"
    echo "  -H \"Authorization: Bearer YOUR_TOKEN\""
else
    echo "❌ Health check failed. Check kotoba.log"
    exit 1
fi
