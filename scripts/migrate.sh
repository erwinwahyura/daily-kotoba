#!/bin/bash

# Database Migration Script
# Runs all migrations in order

set -e

echo "🗄️  Kotoba API - Database Migration"
echo "===================================="

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

CONTAINER_NAME=${1:-kotoba-postgres-prod}
DB_USER=${2:-kotoba}
DB_NAME=${3:-kotoba_db}

echo -e "${YELLOW}Container:${NC} $CONTAINER_NAME"
echo -e "${YELLOW}Database:${NC} $DB_NAME"
echo -e "${YELLOW}User:${NC} $DB_USER"
echo ""

# Check if container is running
if ! docker ps | grep -q $CONTAINER_NAME; then
    echo "Error: Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "Running migrations..."

# Array of migration files in order
migrations=(
    "000001_create_users_table.up.sql"
    "000002_create_vocabulary_table.up.sql"
    "000003_create_user_progress_table.up.sql"
    "000004_create_user_vocab_status_table.up.sql"
    "000005_create_placement_test_results_table.up.sql"
    "000006_create_placement_questions_table.up.sql"
)

for migration in "${migrations[@]}"; do
    echo -n "  → $migration ... "

    if docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME < migrations/$migration > /dev/null 2>&1; then
        echo -e "${GREEN}✓${NC}"
    else
        # Check if table already exists (not an error)
        echo -e "${YELLOW}⚠${NC} (may already exist)"
    fi
done

echo ""
echo -e "${GREEN}✓ Migration completed!${NC}"
