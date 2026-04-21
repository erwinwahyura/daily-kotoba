# Kotoba API Documentation

## Base URL
- **Production:** `http://46.224.127.221:8090/api`
- **Health Check:** `http://46.224.127.221:8090/health`

## Authentication
All endpoints (except `/auth/login` and `/auth/register`) require a Bearer token in the Authorization header:
```
Authorization: Bearer <token>
```

---

## Endpoints

### Auth

#### POST `/auth/register`
Register a new user.

**Request Body:**
```json
{
  "name": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "token": "jwt_token"
  }
}
```

#### POST `/auth/login`
Login existing user.

**Request Body:**
```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**
```json
{
  "message": "Login successful",
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "token": "jwt_token"
  }
}
```

#### GET `/auth/me`
Get current user info.

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "current_level": "N5|N4|N3|N2|N1"
  }
}
```

---

### Kanji Writing Practice

#### GET `/kanji/character/:char`
Get kanji details with stroke data.

**Parameters:**
- `char` - The kanji character (URL-encoded)

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "character": "日",
    "jlpt_level": "N5",
    "stroke_count": 4,
    "meaning": "day, sun, Japan",
    "readings": ["にち", "ひ", "か"],
    "strokes": [
      {
        "stroke_num": 1,
        "path": [{"x": 50, "y": 20}, {"x": 50, "y": 80}],
        "direction": {"start": {"x": 50, "y": 20}, "end": {"x": 50, "y": 80}}
      }
    ]
  }
}
```

#### GET `/kanji/level/:level`
List kanji by JLPT level.

**Parameters:**
- `level` - JLPT level (N5, N4, N3, N2, N1)

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "character": "日",
      "jlpt_level": "N5",
      "stroke_count": 4,
      "meaning": "day, sun, Japan"
    }
  ]
}
```

#### POST `/kanji/practice/start`
Start a practice session.

**Request Body:**
```json
{
  "kanji_char": "日"
}
```

**Response:**
```json
{
  "data": {
    "session_id": "uuid",
    "kanji": "日",
    "total_strokes": 4,
    "started_at": "2026-04-20T18:00:00Z"
  }
}
```

#### POST `/kanji/practice/compare`
Compare user stroke with reference.

**Request Body:**
```json
{
  "session_id": "uuid",
  "stroke_num": 1,
  "path": [
    {"x": 50, "y": 20},
    {"x": 50, "y": 80}
  ]
}
```

**Response:**
```json
{
  "data": {
    "accuracy": 85.5,
    "is_correct": true,
    "feedback": "Excellent! Perfect stroke.",
    "stroke_matched": 1,
    "next_stroke": 2
  }
}
```

#### GET `/kanji/practice/:id`
Get session progress.

**Parameters:**
- `id` - Session ID

**Response:**
```json
{
  "data": {
    "session_id": "uuid",
    "kanji": "日",
    "current_stroke": 2,
    "total_strokes": 4,
    "completed": false,
    "accuracy_scores": [90.0, 85.5]
  }
}
```

#### GET `/kanji/stats`
Get user kanji statistics.

**Response:**
```json
{
  "data": {
    "total_practiced": 10,
    "mastered": 5,
    "by_level": {
      "N5": {"practiced": 10, "mastered": 5}
    },
    "recent_sessions": []
  }
}
```

#### POST `/kanji/seed`
Admin endpoint to seed kanji data.

**Headers:** Authorization with admin token

**Response:**
```json
{
  "message": "Seeded 11 N5 kanji characters",
  "data": {
    "count": 11,
    "kanji": ["日", "月", "火", "水", "木", "金", "人", "大", "小", "上", "下"]
  }
}
```

---

### Vocabulary

#### GET `/vocab/daily`
Get daily vocabulary word.

**Response:**
```json
{
  "data": {
    "vocabulary": {
      "id": "uuid",
      "word": "日本",
      "reading": "にほん",
      "short_meaning": "Japan",
      "detailed_explanation": "...",
      "jlpt_level": "N5"
    },
    "progress": {
      "words_learned": 10,
      "current_streak": 5
    }
  }
}
```

#### POST `/vocab/:id/skip`
Skip/mark vocabulary.

**Request Body:**
```json
{
  "status": "known|skipped|learning"
}
```

#### GET `/vocab/search?q=query&level=N5`
Search vocabulary.

---

### Grammar

#### GET `/grammar/daily`
Get daily grammar pattern.

#### POST `/grammar/:id/skip`
Mark grammar as studied.

#### GET `/grammar/search?q=query`
Search grammar patterns.

---

### Progress

#### GET `/progress`
Get user progress.

#### GET `/progress/stats`
Get progress statistics.

---

### SRS (Spaced Repetition)

#### POST `/srs/init`
Initialize SRS item.

#### POST `/srs/review`
Submit SRS review.

#### GET `/srs/queue`
Get review queue.

---

## Error Responses

### 400 Bad Request
```json
{
  "error": "Invalid request",
  "message": "Detailed error message"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "error": "Not found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "Something went wrong"
}
```

---

## CORS Configuration

The API allows requests from:
- `kotoba-web.erwinwahyura.workers.dev`
- `localhost:3000`
- `localhost:5173`

---

## Data Types

### Point
```json
{
  "x": 50.0,
  "y": 100.0
}
```

### Stroke
```json
{
  "stroke_num": 1,
  "path": [Point],
  "direction": {
    "start": Point,
    "end": Point
  }
}
```

### Kanji
```json
{
  "id": "uuid",
  "character": "日",
  "jlpt_level": "N5",
  "stroke_count": 4,
  "meaning": "day, sun, Japan",
  "readings": ["にち", "ひ", "か"],
  "strokes": [Stroke]
}
```
