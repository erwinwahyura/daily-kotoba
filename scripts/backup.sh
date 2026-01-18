#!/bin/bash

# Database Backup Script

set -e

echo "💾 Kotoba API - Database Backup"
echo "================================"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

CONTAINER_NAME=${1:-kotoba-postgres-prod}
DB_USER=${2:-kotoba}
DB_NAME=${3:-kotoba_db}
BACKUP_DIR=${4:-./backups}

# Create backup directory if it doesn't exist
mkdir -p $BACKUP_DIR

# Generate filename with timestamp
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_FILE="$BACKUP_DIR/kotoba_backup_$TIMESTAMP.sql"

echo -e "${YELLOW}Container:${NC} $CONTAINER_NAME"
echo -e "${YELLOW}Database:${NC} $DB_NAME"
echo -e "${YELLOW}Backup file:${NC} $BACKUP_FILE"
echo ""

# Check if container is running
if ! docker ps | grep -q $CONTAINER_NAME; then
    echo "Error: Container $CONTAINER_NAME is not running"
    exit 1
fi

echo "Creating backup..."

# Create backup
if docker exec $CONTAINER_NAME pg_dump -U $DB_USER $DB_NAME > $BACKUP_FILE; then
    echo -e "${GREEN}✓ Backup created successfully!${NC}"

    # Get file size
    SIZE=$(du -h $BACKUP_FILE | cut -f1)
    echo ""
    echo "Backup details:"
    echo "  File: $BACKUP_FILE"
    echo "  Size: $SIZE"
    echo ""

    # Compress backup
    echo "Compressing backup..."
    gzip $BACKUP_FILE

    COMPRESSED_SIZE=$(du -h $BACKUP_FILE.gz | cut -f1)
    echo -e "${GREEN}✓ Backup compressed!${NC}"
    echo "  Compressed file: $BACKUP_FILE.gz"
    echo "  Compressed size: $COMPRESSED_SIZE"
else
    echo "Error: Failed to create backup"
    exit 1
fi

# List recent backups
echo ""
echo "Recent backups:"
ls -lh $BACKUP_DIR/*.sql.gz 2>/dev/null | tail -5 || echo "  No backups found"

echo ""
echo "To restore from this backup:"
echo "  gunzip $BACKUP_FILE.gz"
echo "  cat $BACKUP_FILE | docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $DB_NAME"
