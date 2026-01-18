# Kotoba API Documentation

Base URL: `http://localhost:8080/api`

## Table of Contents

- [Authentication](#authentication)
- [Endpoints](#endpoints)
  - [Health Check](#health-check)
  - [Authentication Endpoints](#authentication-endpoints)
  - [Placement Test Endpoints](#placement-test-endpoints)
  - [Vocabulary Endpoints](#vocabulary-endpoints)
  - [Progress Endpoints](#progress-endpoints)
- [Error Handling](#error-handling)
- [Response Format](#response-format)

## Authentication

Most endpoints require JWT authentication. Include the JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

Obtain a token by registering or logging in through the auth endpoints.

## Response Format

All responses follow this structure:

### Success Response
```json
{
  "code": 200,
  "message": "Success message",
  "data": {
    // Response data
  }
}
```

### Error Response
```json
{
  "code": 400,
  "message": "Error message"
}
```

## Endpoints

### Health Check

#### GET /health

Check if the API is running.

**Authentication:** Not required

**Response:**
```json
{
  "status": "ok",
  "message": "Kotoba API is running"
}
```

---

## Authentication Endpoints

### POST /api/auth/register

Register a new user account.

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (201):**
```json
{
  "code": 201,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "current_jlpt_level": "N5",
      "created_at": "2026-01-18T00:00:00Z",
      "updated_at": "2026-01-18T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Errors:**
- `409`: User with this email already exists
- `400`: Invalid request body

---

### POST /api/auth/login

Log in to an existing account.

**Authentication:** Not required

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "current_jlpt_level": "N4",
      "created_at": "2026-01-18T00:00:00Z",
      "updated_at": "2026-01-18T00:00:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Errors:**
- `401`: Invalid credentials
- `400`: Invalid request body

---

### GET /api/auth/me

Get current user information.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "User retrieved successfully",
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "current_jlpt_level": "N4",
      "created_at": "2026-01-18T00:00:00Z",
      "updated_at": "2026-01-18T00:00:00Z"
    }
  }
}
```

**Errors:**
- `401`: Unauthorized (invalid or expired token)

---

## Placement Test Endpoints

### GET /api/placement-test

Get all placement test questions.

**Authentication:** Not required

**Response (200):**
```json
{
  "code": 200,
  "message": "Placement test retrieved successfully",
  "data": {
    "questions": [
      {
        "id": "uuid",
        "question": "Choose the correct reading for 食べる",
        "options": ["たべる", "のべる", "かべる", "よべる"],
        "difficulty": "N5",
        "order_index": 1
      },
      // ... 19 more questions
    ]
  }
}
```

**Notes:**
- Returns 20 questions across all JLPT levels
- Answer options are shuffled
- Correct answers are not included in the response

---

### POST /api/placement-test/submit

Submit placement test answers and get assigned JLPT level.

**Authentication:** Required

**Request Body:**
```json
{
  "answers": {
    "question-uuid-1": "たべる",
    "question-uuid-2": "to see, to watch",
    // ... answers for all 20 questions
  }
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Placement test submitted successfully",
  "data": {
    "score": 16,
    "total_questions": 20,
    "assigned_level": "N4",
    "breakdown": {
      "N5": 5,
      "N4": 8,
      "N3": 3,
      "N2": 0,
      "N1": 0
    }
  }
}
```

**Level Assignment Logic:**
- N1: 2+ correct N1 questions
- N2: 2+ correct N2 questions
- N3: 3+ correct N3 questions
- N4: 3+ correct N4 questions
- N5: Default (beginner)

**Errors:**
- `401`: Unauthorized
- `400`: No answers provided

---

### GET /api/placement-test/result

Get user's most recent placement test result.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "Placement test result retrieved successfully",
  "data": {
    "id": "uuid",
    "user_id": "uuid",
    "test_score": 16,
    "assigned_level": "N4",
    "completed_at": "2026-01-18T00:00:00Z"
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: No placement test result found

---

## Vocabulary Endpoints

### GET /api/vocab/daily

Get the current vocabulary word for the user based on their progress.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "Daily word retrieved successfully",
  "data": {
    "vocabulary": {
      "id": "uuid",
      "word": "諦める",
      "reading": "あきらめる",
      "short_meaning": "to give up",
      "detailed_explanation": "Used when abandoning an effort or goal...",
      "example_sentences": [
        "試験に合格するのを諦めた - I gave up on passing the exam",
        "彼女は夢を諦めなかった - She didn't give up on her dream",
        "もう諦めたほうがいい - You should give up already"
      ],
      "usage_notes": "Commonly used in both casual and formal contexts...",
      "jlpt_level": "N4",
      "index_position": 0,
      "created_at": "2026-01-18T00:00:00Z"
    },
    "progress": {
      "current_index": 0,
      "total_words_in_level": 50,
      "words_learned": 0,
      "streak_days": 0
    }
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: No vocabulary found for user's level

---

### GET /api/vocab/:id

Get detailed information about a specific vocabulary word.

**Authentication:** Required

**Parameters:**
- `id` (path): Vocabulary word UUID

**Response (200):**
```json
{
  "code": 200,
  "message": "Vocabulary retrieved successfully",
  "data": {
    "vocabulary": {
      "id": "uuid",
      "word": "諦める",
      "reading": "あきらめる",
      "short_meaning": "to give up",
      "detailed_explanation": "...",
      "example_sentences": ["..."],
      "usage_notes": "...",
      "jlpt_level": "N4",
      "index_position": 0,
      "created_at": "2026-01-18T00:00:00Z"
    }
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: Vocabulary not found

---

### POST /api/vocab/:id/skip

Mark current word as known/skipped and advance to the next word.

**Authentication:** Required

**Parameters:**
- `id` (path): Vocabulary word UUID

**Request Body (optional):**
```json
{
  "status": "known"  // or "skipped"
}
```

**Response (200):**
```json
{
  "code": 200,
  "message": "Moved to next word successfully",
  "data": {
    "vocabulary": {
      // Next vocabulary word details
    },
    "progress": {
      "current_index": 1,
      "total_words_in_level": 50,
      "words_learned": 0,
      "streak_days": 1
    }
  }
}
```

**Notes:**
- Automatically advances user to the next word
- Updates user progress and streak
- Cycles back to index 0 when reaching the end of level

**Errors:**
- `401`: Unauthorized
- `404`: Vocabulary not found

---

### GET /api/vocab/level/:level

Get all vocabulary words for a specific JLPT level (paginated).

**Authentication:** Required

**Parameters:**
- `level` (path): JLPT level (N5, N4, N3, N2, N1)

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Example:** `/api/vocab/level/N4?page=1&limit=10`

**Response (200):**
```json
{
  "code": 200,
  "message": "Vocabulary list retrieved successfully",
  "data": {
    "vocabulary": [
      {
        "id": "uuid",
        "word": "諦める",
        "reading": "あきらめる",
        "short_meaning": "to give up",
        "detailed_explanation": "...",
        "example_sentences": ["..."],
        "usage_notes": "...",
        "jlpt_level": "N4",
        "index_position": 0,
        "created_at": "2026-01-18T00:00:00Z"
      },
      // ... more words
    ],
    "pagination": {
      "page": 1,
      "limit": 10,
      "total": 50,
      "total_pages": 5
    }
  }
}
```

**Errors:**
- `401`: Unauthorized
- `400`: Invalid level or pagination parameters

---

## Progress Endpoints

### GET /api/progress

Get user's current learning progress.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "Progress retrieved successfully",
  "data": {
    "progress": {
      "current_vocab_index": 15,
      "current_level": "N4",
      "streak_days": 7,
      "last_study_date": "2026-01-18",
      "words_learned": 12,
      "words_skipped": 3,
      "total_words_in_level": 50
    }
  }
}
```

**Errors:**
- `401`: Unauthorized

---

### GET /api/progress/stats

Get detailed learning statistics.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "Statistics retrieved successfully",
  "data": {
    "stats": {
      "total_days_active": 45,
      "current_streak": 7,
      "longest_streak": 15,
      "words_learned_count": 194,
      "words_skipped_count": 12,
      "current_level": "N4",
      "current_index": 15
    }
  }
}
```

**Errors:**
- `401`: Unauthorized

---

## Error Handling

### Common HTTP Status Codes

| Code | Meaning |
|------|---------|
| 200 | Success |
| 201 | Created (successful registration) |
| 400 | Bad Request (invalid input) |
| 401 | Unauthorized (missing or invalid token) |
| 404 | Not Found |
| 409 | Conflict (duplicate email) |
| 500 | Internal Server Error |

### Example Error Response

```json
{
  "code": 401,
  "message": "Invalid or expired token"
}
```

---

## Rate Limiting

Currently, no rate limiting is implemented. This will be added in future versions.

---

## CORS

CORS is configured to allow all origins in development. In production, configure specific allowed origins in the environment.

---

## Changelog

### v1.0.0 (2026-01-18)
- Initial API release
- Authentication system
- Placement test (20 questions)
- Vocabulary system (50 N4 words)
- Progress tracking
- All CRUD endpoints

---

## Support

For issues or questions:
- GitHub Issues: https://github.com/yourusername/kotoba-api/issues
- Email: support@example.com

---

**Last Updated:** January 18, 2026
