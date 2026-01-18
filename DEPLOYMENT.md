# Kotoba API - Deployment Guide

This guide covers deploying the Kotoba API to a production server (Hetzner VPS).

## Prerequisites

- Docker and Docker Compose installed on the server
- Domain name (optional, for HTTPS)
- SSH access to the server

## Quick Start (Local Testing)

1. Build and start the production stack:
```bash
make docker-build
make docker-up
```

2. Check logs:
```bash
make logs
```

3. Test the API:
```bash
curl http://localhost:8080/health
```

## Production Deployment on Hetzner VPS

### Step 1: Server Setup

1. SSH into your Hetzner VPS:
```bash
ssh root@your-server-ip
```

2. Install Docker and Docker Compose:
```bash
# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
apt install docker-compose -y

# Verify installation
docker --version
docker-compose --version
```

3. Create application directory:
```bash
mkdir -p /opt/kotoba
cd /opt/kotoba
```

### Step 2: Deploy Application

1. Clone the repository (or upload files):
```bash
git clone https://github.com/yourusername/kotoba-api.git .
```

2. Create production environment file:
```bash
cp .env.production.example .env.production
nano .env.production
```

Edit the file and set secure values:
```env
DB_USER=kotoba
DB_PASSWORD=<generate-secure-password>
DB_NAME=kotoba_db
DB_PORT=5432

JWT_SECRET=<generate-long-random-string-minimum-32-chars>
JWT_EXPIRATION_HOURS=24

API_PORT=8080
SERVER_ENV=production
```

Generate secure values:
```bash
# Generate secure DB password
openssl rand -base64 32

# Generate JWT secret
openssl rand -base64 64
```

3. Build and start the application:
```bash
docker-compose -f docker-compose.prod.yml --env-file .env.production up -d --build
```

4. Verify services are running:
```bash
docker-compose -f docker-compose.prod.yml ps
```

5. Check logs:
```bash
docker-compose -f docker-compose.prod.yml logs -f api
```

### Step 3: Database Setup

1. Wait for PostgreSQL to be ready (check logs)

2. Run migrations:
```bash
# Access the API container
docker exec -it kotoba-api-prod sh

# Migrations are copied to /root/migrations
# They will need to be run manually or via a migration tool
```

For now, run migrations directly on PostgreSQL:
```bash
docker exec -it kotoba-postgres-prod psql -U kotoba -d kotoba_db

-- Then run each migration file content manually
```

3. Seed initial data (from your local machine):
```bash
# Make sure DATABASE_URL points to production
go run cmd/seed/main.go
go run cmd/seed/seed_placement.go
```

### Step 4: Firewall Configuration

1. Configure UFW firewall:
```bash
# Allow SSH
ufw allow 22/tcp

# Allow HTTP
ufw allow 80/tcp

# Allow HTTPS
ufw allow 443/tcp

# Allow API port (if not using reverse proxy)
ufw allow 8080/tcp

# Enable firewall
ufw enable
```

### Step 5: Reverse Proxy with Nginx (Optional but Recommended)

1. Install Nginx:
```bash
apt install nginx -y
```

2. Create Nginx configuration:
```bash
nano /etc/nginx/sites-available/kotoba-api
```

```nginx
server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

3. Enable the site:
```bash
ln -s /etc/nginx/sites-available/kotoba-api /etc/nginx/sites-enabled/
nginx -t
systemctl restart nginx
```

4. Install SSL with Let's Encrypt:
```bash
apt install certbot python3-certbot-nginx -y
certbot --nginx -d api.yourdomain.com
```

### Step 6: Monitoring and Maintenance

1. View logs:
```bash
docker-compose -f docker-compose.prod.yml logs -f
```

2. Restart services:
```bash
docker-compose -f docker-compose.prod.yml restart
```

3. Update application:
```bash
git pull
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d --build
```

4. Backup database:
```bash
docker exec kotoba-postgres-prod pg_dump -U kotoba kotoba_db > backup_$(date +%Y%m%d_%H%M%S).sql
```

5. Restore database:
```bash
cat backup_file.sql | docker exec -i kotoba-postgres-prod psql -U kotoba -d kotoba_db
```

## Environment Variables Reference

| Variable | Description | Example |
|----------|-------------|---------|
| DB_USER | PostgreSQL username | kotoba |
| DB_PASSWORD | PostgreSQL password | <secure-password> |
| DB_NAME | Database name | kotoba_db |
| DB_PORT | PostgreSQL port | 5432 |
| JWT_SECRET | JWT signing secret | <random-64-char-string> |
| JWT_EXPIRATION_HOURS | JWT token expiration | 24 |
| API_PORT | API server port | 8080 |
| SERVER_ENV | Environment mode | production |

## Troubleshooting

### API container won't start
```bash
# Check logs
docker logs kotoba-api-prod

# Common issues:
# - Database not ready: wait a few seconds and check postgres logs
# - Environment variables: verify .env.production file
```

### Database connection issues
```bash
# Check if PostgreSQL is running
docker exec kotoba-postgres-prod pg_isready -U kotoba

# Test connection from API container
docker exec -it kotoba-api-prod sh
# Try connecting to postgres
```

### Port conflicts
```bash
# Check what's using port 8080
lsof -i :8080

# Change API_PORT in .env.production if needed
```

## Security Checklist

- [ ] Changed default database password
- [ ] Generated secure JWT secret (minimum 32 characters)
- [ ] Firewall configured (UFW)
- [ ] SSL certificate installed (Let's Encrypt)
- [ ] Database backups scheduled
- [ ] Container logs configured
- [ ] Non-root user for application (if not using Docker)
- [ ] Rate limiting configured (Nginx)
- [ ] CORS properly configured

## Performance Optimization

1. Database connection pooling (already configured in code)
2. Enable gzip compression in Nginx
3. Add caching headers for static responses
4. Monitor memory usage and scale as needed

## Next Steps

After deployment:
1. Test all API endpoints
2. Set up monitoring (Prometheus/Grafana)
3. Configure automated backups
4. Set up CI/CD pipeline
5. Add rate limiting
6. Configure log aggregation
