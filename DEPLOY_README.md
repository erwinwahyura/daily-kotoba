# 🐹 Kotoba Kanji Backend — Deploy Instructions

## Status: READY TO DEPLOY ✅

**What's included:**
- 11 N5 kanji with full stroke data (日, 月, 火, 水, 木, 金, 人, 大, 小, 上, 下)
- Complete REST API with 7 endpoints
- Stroke comparison algorithm with scoring
- Practice session tracking
- CORS enabled for Cloudflare Workers

## Quick Deploy (Hetzner VPS)

```bash
# 1. SSH to server
ssh deploy@46.224.127.221

# 2. Pull latest code
cd ~/daily-kotoba
git pull origin main

# 3. Build
export PORT=8090
go build -o kotoba-api ./cmd/api/

# 4. Stop old, start new
pkill -f "./kotoba-api" || true
nohup ./kotoba-api > kotoba.log 2>&1 &

# 5. Verify
curl http://localhost:8090/health
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/kanji/character/:char` | Get kanji with stroke data |
| GET | `/api/kanji/level/:level` | List kanji by JLPT level |
| POST | `/api/kanji/practice/start` | Start practice session |
| POST | `/api/kanji/practice/compare` | Submit stroke, get accuracy |
| GET | `/api/kanji/practice/:id` | Get session progress |
| GET | `/api/kanji/stats` | User statistics |
| POST | `/api/kanji/seed` | Admin: seed kanji data |

## Seed Kanji Data

```bash
curl -X POST http://localhost:8090/api/kanji/seed \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Frontend Integration

```javascript
const API_BASE = 'https://kotoba.erwarx.com'; // or localhost:8090

// Start practice
fetch(`${API_BASE}/api/kanji/practice/start`, {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: JSON.stringify({ kanji_char: '日' })
});

// Submit stroke
fetch(`${API_BASE}/api/kanji/practice/compare`, {
  method: 'POST',
  headers: { 'Authorization': `Bearer ${token}` },
  body: JSON.stringify({
    session_id: 'xxx',
    stroke_num: 1,
    user_path: [{x: 10, y: 50}, {x: 90, y: 50}]
  })
});
```

## CORS Allowed Origins

- `https://kotoba-web.erwinwahyura.workers.dev`
- `http://localhost:3000`
- `http://localhost:5173`

---
Deployed by: Go 🐹⚡
Date: 2026-04-20
