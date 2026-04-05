# SQLite Migration Guide for Hetzner

## Why SQLite on Hetzner?

With a **persistent volume** on Hetzner, SQLite is perfect:
- ✅ Single file database survives redeploys
- ✅ No external database service needed
- ✅ Simple backups (just copy the file)
- ✅ Faster than network roundtrips to Postgres

## Deployment Steps

### 1. Prepare Hetzner Server

Create a persistent volume (if not already done):
```bash
# On Hetzner Cloud Console or CLI
hcloud volume create --name kotoba-data --size 10 --server <your-server-id>
```

Mount the volume:
```bash
# SSH to your server
ssh root@<your-server-ip>

# Create mount point
mkdir -p /mnt/kotoba-data

# Mount the volume (adjust device name as needed)
mount /dev/disk/by-id/scsi-0HC_Volume_<volume-id> /mnt/kotoba-data

# Add to fstab for persistence across reboots
echo '/dev/disk/by-id/scsi-0HC_Volume_<volume-id> /mnt/kotoba-data ext4 defaults 0 0' >> /etc/fstab
```

### 2. Deploy with Docker Compose

```bash
# Clone and build
git clone <your-repo>
cd daily-kotoba

# Copy environment
cp .env.hetzner .env

# IMPORTANT: Edit JWT_SECRET in .env!
nano .env

# Deploy
docker-compose -f docker-compose.hetzner.yml up -d --build
```

### 3. Initialize Database

```bash
# Run migrations (creates SQLite schema)
docker-compose -f docker-compose.hetzner.yml exec kotoba-api ./kotoba-api migrate

# Or manually with SQLite CLI
docker-compose -f docker-compose.hetzner.yml exec kotoba-api sh
sqlite3 /data/kotoba.db < /app/migrations_sqlite/001_initial_schema.sql
```

### 4. Verify Persistence

```bash
# Check database exists on volume
docker-compose -f docker-compose.hetzner.yml exec kotoba-api ls -la /data/

# Redeploy test
docker-compose -f docker-compose.hetzner.yml down
docker-compose -f docker-compose.hetzner.yml up -d

# Data should still be there!
docker-compose -f docker-compose.hetzner.yml exec kotoba-api sqlite3 /data/kotoba.db "SELECT COUNT(*) FROM vocabulary;"
```

## Data Migration (PostgreSQL → SQLite)

If you have existing PostgreSQL data to preserve:

```bash
# 1. Export from PostgreSQL
pg_dump --data-only --inserts kotoba_db > postgres_data.sql

# 2. Transform and import (use the migration script)
./scripts/migrate_postgres_to_sqlite.sh

# 3. Copy to Hetzner volume
scp kotoba.db root@<hetzner-ip>:/mnt/kotoba-data/
```

## Backup Strategy

Since SQLite is a single file, backups are trivial:

```bash
# Automated backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
cp /mnt/kotoba-data/kotoba.db /mnt/kotoba-backups/kotoba_backup_$DATE.db
# Keep only last 30 days
find /mnt/kotoba-backups -name "kotoba_backup_*.db" -mtime +30 -delete
```

Or use Litestream for continuous backup to S3:
```bash
# Litestream config
litestream replicate /data/kotoba.db s3://your-bucket/kotoba.db
```

## Volume Management

### Resize volume (if needed):
```bash
hcloud volume resize --size 20 kotoba-data
# Then resize filesystem on server
resize2fs /dev/disk/by-id/scsi-0HC_Volume_<volume-id>
```

### Move to new server:
```bash
# Detach from old server
hcloud volume detach kotoba-data

# Attach to new server
hcloud volume attach kotoba-data <new-server-id>
```

## Troubleshooting

**Issue: Permission denied on /data**
- Solution: Ensure volume is mounted with correct permissions
- `chown -R 1000:1000 /mnt/kotoba-data` (Docker runs as non-root)

**Issue: WAL mode errors**
- SQLite WAL requires write permissions on directory
- Ensure /data is writable by container user

**Issue: Database locked**
- Check if multiple processes accessing same file
- SQLite allows multiple readers but single writer
- Use connection pooling carefully (limit to 1 write connection)

## Configuration Reference

| Variable | Description | Hetzner Value |
|----------|-------------|---------------|
| DB_DRIVER | Database type | `sqlite` |
| SQLITE_PATH | Database file location | `/data/kotoba.db` |
| JWT_SECRET | Signing secret | Change this! |

## Why This Works on Hetzner

Unlike serverless platforms (Vercel, Cloudflare Workers), Hetzner provides:
- Persistent block storage volumes
- Long-running containers
- Filesystem persistence across deploys

Your SQLite file lives at `/data/kotoba.db` on a mounted volume that survives:
- Container restarts
- Docker image updates
- Server reboots (with fstab entry)

**The key:** The volume is separate from the container filesystem.
