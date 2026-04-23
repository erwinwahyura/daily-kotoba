# Kotoba API Documentation - Expanded

## Table of Contents
1. [Quick Start](#quick-start)
2. [Authentication](#authentication)
3. [Kanji API](#kanji-api)
4. [Nichijou Conversation API](#nichijou-conversation-api)
5. [Vocabulary API](#vocabulary-api)
6. [Grammar API](#grammar-api)
7. [SRS API](#srs-api)
8. [Progress API](#progress-api)
9. [Client SDK Examples](#client-sdk-examples)
10. [Error Handling](#error-handling)

---

## Quick Start

### Base URLs
```
Production: http://46.224.127.221:8090/api
Health Check: http://46.224.127.221:8090/health
```

### Authentication
All endpoints require a Bearer token except `/auth/login` and `/auth/register`:
```
Authorization: Bearer <jwt_token>
```

### CORS Origins Allowed
- `https://kotoba-web.erwinwahyura.workers.dev`
- `http://localhost:3000`
- `http://localhost:5173`
- `http://localhost:8080`

---

## Authentication

### Register
```http
POST /auth/register
Content-Type: application/json

{
  "name": "Tanaka Yuki",
  "email": "tanaka@example.com",
  "password": "securepassword123"
}
```

**Response:**
```json
{
  "message": "User registered successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Tanaka Yuki",
    "email": "tanaka@example.com",
    "token": "eyJhbGciOiJIUzI1NiIs..."
  }
}
```

### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "tanaka@example.com",
  "password": "securepassword123"
}
```

### Get Current User
```http
GET /auth/me
Authorization: Bearer <token>
```

---

## Kanji API

### Get Kanji by Character
```http
GET /kanji/character/%E6%97%A5
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "character": "日",
    "jlpt_level": "N5",
    "stroke_count": 4,
    "meaning": "day, sun, Japan",
    "readings": ["にち", "ひ", "か"],
    "strokes": [
      {
        "stroke_num": 1,
        "path": [
          {"x": 50, "y": 20},
          {"x": 50, "y": 80}
        ],
        "direction": {
          "start": {"x": 50, "y": 20},
          "end": {"x": 50, "y": 80}
        }
      },
      {
        "stroke_num": 2,
        "path": [
          {"x": 20, "y": 50},
          {"x": 80, "y": 50}
        ],
        "direction": {
          "start": {"x": 20, "y": 50},
          "end": {"x": 80, "y": 50}
        }
      }
    ]
  }
}
```

### List Kanji by JLPT Level
```http
GET /kanji/level/N5
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": [
    {
      "id": "...",
      "character": "日",
      "jlpt_level": "N5",
      "stroke_count": 4,
      "meaning": "day, sun, Japan"
    },
    {
      "id": "...",
      "character": "月",
      "jlpt_level": "N5",
      "stroke_count": 4,
      "meaning": "month, moon"
    }
  ],
  "total_count": 11,
  "level": "N5"
}
```

### Start Practice Session
```http
POST /kanji/practice/start
Authorization: Bearer <token>
Content-Type: application/json

{
  "kanji_char": "日"
}
```

**Response:**
```json
{
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440002",
    "kanji": "日",
    "total_strokes": 4,
    "started_at": "2026-04-22T10:30:00Z"
  }
}
```

### Submit Stroke for Comparison
```http
POST /kanji/practice/compare
Authorization: Bearer <token>
Content-Type: application/json

{
  "session_id": "550e8400-e29b-41d4-a716-446655440002",
  "stroke_num": 1,
  "user_path": [
    {"x": 48, "y": 22},
    {"x": 49, "y": 35},
    {"x": 50, "y": 50},
    {"x": 51, "y": 65},
    {"x": 52, "y": 78}
  ]
}
```

**Response:**
```json
{
  "data": {
    "accuracy": 87.5,
    "is_correct": true,
    "feedback": "Great stroke! Slightly off center at the end.",
    "stroke_matched": 1,
    "next_stroke": 2,
    "direction": "vertical"
  }
}
```

### Get Session Progress
```http
GET /kanji/practice/550e8400-e29b-41d4-a716-446655440002
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440002",
    "kanji": "日",
    "current_stroke": 3,
    "total_strokes": 4,
    "completed": false,
    "accuracy_scores": [92.0, 87.5, 90.0],
    "remaining_strokes": 1
  }
}
```

### Get User Statistics
```http
GET /kanji/stats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "total_practiced": 25,
    "mastered": 12,
    "by_level": {
      "N5": {"practiced": 25, "mastered": 12, "accuracy_avg": 85.5},
      "N4": {"practiced": 0, "mastered": 0, "accuracy_avg": 0}
    },
    "recent_sessions": [
      {
        "kanji": "日",
        "accuracy": 87.5,
        "completed_at": "2026-04-22T10:35:00Z"
      }
    ],
    "streak_days": 5
  }
}
```

---

## Nichijou Conversation API

### Get Available Scenarios
```http
GET /nichijou/scenarios?level=5&category=daily
Authorization: Bearer <token>
```

**Query Parameters:**
- `level` (optional): JLPT level (1-5)
- `category` (optional): daily, social, business, emergency

**Response:**
```json
{
  "data": {
    "scenarios": [
      {
        "id": "konbini",
        "name": "Convenience Store",
        "name_ja": "コンビニ",
        "emoji": "🏪",
        "description": "Practice ordering, asking about point cards, and small talk at a Japanese convenience store",
        "jlpt_level": 5,
        "category": "daily",
        "estimated_duration": "3-5 min",
        "difficulty": "beginner"
      },
      {
        "id": "restaurant",
        "name": "Restaurant",
        "name_ja": "レストラン",
        "emoji": "🍜",
        "description": "Order food, ask for recommendations, and handle payments",
        "jlpt_level": 5,
        "category": "daily",
        "estimated_duration": "5-7 min",
        "difficulty": "beginner"
      },
      {
        "id": "directions",
        "name": "Asking Directions",
        "name_ja": "道を聞く",
        "emoji": "🗺️",
        "description": "Ask for and give directions to landmarks",
        "jlpt_level": 4,
        "category": "daily",
        "estimated_duration": "4-6 min",
        "difficulty": "elementary"
      },
      {
        "id": "weather_chat",
        "name": "Weather Small Talk",
        "name_ja": "天気の話",
        "emoji": "☀️",
        "description": "Casual conversation about weather and seasons",
        "jlpt_level": 5,
        "category": "social",
        "estimated_duration": "2-3 min",
        "difficulty": "beginner"
      },
      {
        "id": "weekend_plans",
        "name": "Weekend Plans",
        "name_ja": "週末の予定",
        "emoji": "🎉",
        "description": "Discuss what you did or will do on the weekend",
        "jlpt_level": 4,
        "category": "social",
        "estimated_duration": "4-5 min",
        "difficulty": "elementary"
      }
    ]
  }
}
```

### Start Conversation Session
```http
POST /nichijou/chat/start
Authorization: Bearer <token>
Content-Type: application/json

{
  "scenario_id": "konbini",
  "level": 5
}
```

**Response:**
```json
{
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440003",
    "scenario": {
      "id": "konbini",
      "name": "Convenience Store",
      "emoji": "🏪"
    },
    "ai_greeting": "いらっしゃいませ！お弁当をお探しですか？",
    "context": {
      "setting": "You're at a 7-Eleven in Tokyo",
      "your_role": "Customer",
      "ai_role": "Store clerk",
      "goal": "Buy a bento and drinks for lunch"
    },
    "started_at": "2026-04-22T10:40:00Z"
  }
}
```

### Send Message
```http
POST /nichijou/chat/message
Authorization: Bearer <token>
Content-Type: application/json

{
  "session_id": "550e8400-e29b-41d4-a716-446655440003",
  "message": "はい、お弁当をください。これを温めてもらえますか？"
}
```

**Response:**
```json
{
  "data": {
    "user_message": {
      "id": "msg-001",
      "content": "はい、お弁当をください。これを温めてもらえますか？",
      "naturalness_score": 85,
      "correction": {
        "original": "これを温めてもらえますか？",
        "suggestion": "これを温めていただけますか？",
        "explanation": "Using いただけます is more polite in this context",
        "grammar_point": "～ていただけますか (polite request)"
      }
    },
    "ai_response": {
      "id": "msg-002",
      "content": "かしこまりました。お弁当をお温めしますね。レンジで1分30秒ほどお時間いただきます。",
      "translation": "Certainly. I'll warm up the bento. It will take about 1 minute 30 seconds in the microwave.",
      "vocabulary": [
        {
          "word": "かしこまりました",
          "reading": "かしこまりました",
          "meaning": "Certainly/Understood (very polite)"
        },
        {
          "word": "お温め",
          "reading": "おあたため",
          "meaning": "Warming up (honorific)"
        }
      ]
    },
    "session_status": {
      "message_count": 2,
      "naturalness_avg": 85,
      "suggested_responses": [
        "はい、お願いします",
        "わかりました。待ちます",
        "他に飲み物もください"
      ]
    }
  }
}
```

### End Session
```http
POST /nichijou/chat/end/550e8400-e29b-41d4-a716-446655440003
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "session_summary": {
      "session_id": "550e8400-e29b-41d4-a716-446655440003",
      "scenario": "konbini",
      "duration_seconds": 180,
      "message_count": 8,
      "naturalness_avg": 82,
      "completed": true,
      "feedback": "Great job! You successfully ordered and used polite forms well. Try using いただけますか more often for extra politeness."
    },
    "improvements": [
      {
        "type": "grammar",
        "point": "～ていただけますか",
        "example": "これを温めていただけますか？"
      },
      {
        "type": "vocabulary",
        "point": "かしこまりました",
        "example": "Store clerk's response"
      }
    ],
    "xp_earned": 50
  }
}
```

### Get Session History
```http
GET /nichijou/chat/history/550e8400-e29b-41d4-a716-446655440003
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "session_id": "550e8400-e29b-41d4-a716-446655440003",
    "messages": [
      {
        "id": "msg-001",
        "sender": "ai",
        "content": "いらっしゃいませ！お弁当をお探しですか？",
        "timestamp": "2026-04-22T10:40:00Z"
      },
      {
        "id": "msg-002",
        "sender": "user",
        "content": "はい、お弁当をください。",
        "naturalness_score": 90,
        "timestamp": "2026-04-22T10:40:15Z"
      }
    ]
  }
}
```

### Get User Conversation Stats
```http
GET /nichijou/stats
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "total_sessions": 12,
    "total_messages": 86,
    "average_naturalness": 78,
    "scenarios_completed": {
      "konbini": 3,
      "restaurant": 2,
      "directions": 1,
      "weather_chat": 4,
      "weekend_plans": 2
    },
    "improvement_trends": {
      "grammar": "+15%",
      "vocabulary": "+8%",
      "naturalness": "+12%"
    },
    "recent_sessions": [
      {
        "scenario": "konbini",
        "date": "2026-04-22T10:40:00Z",
        "naturalness": 82,
        "duration": "3m"
      }
    ]
  }
}
```

---

## Vocabulary API

### Get Daily Vocabulary
```http
GET /vocab/daily
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "vocabulary": {
      "id": "vocab-001",
      "word": "日本",
      "reading": "にほん",
      "short_meaning": "Japan",
      "detailed_explanation": "The country of Japan, consisting of four main islands",
      "jlpt_level": "N5",
      "example_sentences": [
        {
          "japanese": "日本に行きたいです。",
          "reading": "にほんにいきたいです。",
          "english": "I want to go to Japan."
        }
      ],
      "tags": ["country", "common"]
    },
    "progress": {
      "words_learned": 150,
      "current_streak": 7,
      "total_words": 800
    }
  }
}
```

### Search Vocabulary
```http
GET /vocab/search?q=日本&level=N5&page=1
Authorization: Bearer <token>
```

**Query Parameters:**
- `q` (required): Search query
- `level` (optional): N5, N4, N3, N2, N1
- `page` (optional): Page number (default: 1)
- `limit` (optional): Results per page (default: 20)

### Mark Vocabulary Status
```http
POST /vocab/vocab-001/skip
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "known"
}
```

**Status Options:** `known`, `skipped`, `learning`

---

## Grammar API

### Get Daily Grammar Pattern
```http
GET /grammar/daily
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "grammar": {
      "id": "gram-001",
      "pattern": "～たいです",
      "meaning": "Want to do ~",
      "jlpt_level": "N5",
      "explanation": "Expresses desire to do something. Attach to verb stem.",
      "formation": "Verb (stem) + たいです",
      "examples": [
        {
          "japanese": "日本に行きたいです。",
          "reading": "にほんにいきたいです。",
          "english": "I want to go to Japan."
        },
        {
          "japanese": "寿司を食べたいです。",
          "reading": "すしをたべたいです。",
          "english": "I want to eat sushi."
        }
      ],
      "related_patterns": ["～たくないです", "～たかったです"]
    }
  }
}
```

---

## SRS API

### Initialize SRS Item
```http
POST /srs/init
Authorization: Bearer <token>
Content-Type: application/json

{
  "item_type": "vocabulary",
  "item_id": "vocab-001"
}
```

### Submit Review
```http
POST /srs/review
Authorization: Bearer <token>
Content-Type: application/json

{
  "item_id": "vocab-001",
  "rating": 4
}
```

**Rating Scale:**
- `1`: Again (failed)
- `2`: Hard
- `3`: Good
- `4`: Easy

### Get Review Queue
```http
GET /srs/queue
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "due_count": 15,
    "new_count": 5,
    "total_reviews": 20,
    "items": [
      {
        "id": "vocab-001",
        "type": "vocabulary",
        "word": "日本",
        "reading": "にほん",
        "due_date": "2026-04-22T10:00:00Z",
        "interval_days": 3
      }
    ]
  }
}
```

---

## Progress API

### Get User Progress
```http
GET /progress
Authorization: Bearer <token>
```

**Response:**
```json
{
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "current_level": "N5",
    "overall_progress": {
      "vocabulary": 18.75,
      "grammar": 12.50,
      "kanji": 10.68,
      "conversation": 5.00
    },
    "study_streak": {
      "current": 7,
      "longest": 14
    },
    "last_study_date": "2026-04-22T10:40:00Z"
  }
}
```

### Get Detailed Statistics
```http
GET /progress/stats
Authorization: Bearer <token>
```

---

## Client SDK Examples

### JavaScript/TypeScript SDK

```typescript
// kotoba-client.ts
const API_BASE = 'http://46.224.127.221:8090/api';

class KotobaClient {
  private token: string | null = null;

  setToken(token: string) {
    this.token = token;
    localStorage.setItem('kotoba_token', token);
  }

  private async request(endpoint: string, options: RequestInit = {}) {
    const url = `${API_BASE}${endpoint}`;
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      ...options.headers as Record<string, string>
    };

    if (this.token) {
      headers['Authorization'] = `Bearer ${this.token}`;
    }

    const response = await fetch(url, {
      ...options,
      headers
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'API request failed');
    }

    return response.json();
  }

  // Auth
  async login(email: string, password: string) {
    const result = await this.request('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ email, password })
    });
    this.setToken(result.data.token);
    return result;
  }

  async register(name: string, email: string, password: string) {
    const result = await this.request('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ name, email, password })
    });
    this.setToken(result.data.token);
    return result;
  }

  // Kanji
  async getKanji(character: string) {
    return this.request(`/kanji/character/${encodeURIComponent(character)}`);
  }

  async listKanji(level: string) {
    return this.request(`/kanji/level/${level}`);
  }

  async startKanjiPractice(kanjiChar: string) {
    return this.request('/kanji/practice/start', {
      method: 'POST',
      body: JSON.stringify({ kanji_char: kanjiChar })
    });
  }

  async compareStroke(sessionId: string, strokeNum: number, userPath: Point[]) {
    return this.request('/kanji/practice/compare', {
      method: 'POST',
      body: JSON.stringify({
        session_id: sessionId,
        stroke_num: strokeNum,
        user_path: userPath
      })
    });
  }

  // Nichijou Conversation
  async getScenarios(level?: number, category?: string) {
    const params = new URLSearchParams();
    if (level) params.append('level', level.toString());
    if (category) params.append('category', category);
    return this.request(`/nichijou/scenarios?${params}`);
  }

  async startChat(scenarioId: string, level: number) {
    return this.request('/nichijou/chat/start', {
      method: 'POST',
      body: JSON.stringify({ scenario_id: scenarioId, level })
    });
  }

  async sendMessage(sessionId: string, message: string) {
    return this.request('/nichijou/chat/message', {
      method: 'POST',
      body: JSON.stringify({ session_id: sessionId, message })
    });
  }

  async endChat(sessionId: string) {
    return this.request(`/nichijou/chat/end/${sessionId}`, {
      method: 'POST'
    });
  }

  // Vocabulary
  async getDailyVocab() {
    return this.request('/vocab/daily');
  }

  async searchVocab(query: string, level?: string) {
    const params = new URLSearchParams({ q: query });
    if (level) params.append('level', level);
    return this.request(`/vocab/search?${params}`);
  }

  // SRS
  async getReviewQueue() {
    return this.request('/srs/queue');
  }

  async submitReview(itemId: string, rating: number) {
    return this.request('/srs/review', {
      method: 'POST',
      body: JSON.stringify({ item_id: itemId, rating })
    });
  }
}

// Usage example
const client = new KotobaClient();

// Login
await client.login('user@example.com', 'password');

// Get kanji
const kanji = await client.getKanji('日');
console.log(kanji.data.character); // 日

// Start conversation
const chat = await client.startChat('konbini', 5);
console.log(chat.data.ai_greeting); // いらっしゃいませ！

// Send message
const response = await client.sendMessage(chat.data.session_id, 'お弁当をください');
console.log(response.data.ai_response.content);
```

### React Hook Example

```typescript
// useKotoba.ts
import { useState, useCallback } from 'react';

export function useKotobaChat() {
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [loading, setLoading] = useState(false);

  const startChat = useCallback(async (scenarioId: string, level: number) => {
    setLoading(true);
    try {
      const result = await client.startChat(scenarioId, level);
      setSessionId(result.data.session_id);
      setMessages([
        { sender: 'ai', content: result.data.ai_greeting }
      ]);
      return result;
    } finally {
      setLoading(false);
    }
  }, []);

  const sendMessage = useCallback(async (message: string) => {
    if (!sessionId) return;
    
    setMessages(prev => [...prev, { sender: 'user', content: message }]);
    setLoading(true);
    
    try {
      const result = await client.sendMessage(sessionId, message);
      setMessages(prev => [
        ...prev,
        { 
          sender: 'ai', 
          content: result.data.ai_response.content,
          correction: result.data.user_message.correction
        }
      ]);
      return result;
    } finally {
      setLoading(false);
    }
  }, [sessionId]);

  return { messages, loading, startChat, sendMessage };
}
```

---

## Error Handling

### HTTP Status Codes

| Code | Meaning | Action |
|------|---------|--------|
| 200 | Success | Request completed successfully |
| 201 | Created | Resource created successfully |
| 400 | Bad Request | Check request body/parameters |
| 401 | Unauthorized | Token missing or invalid - re-authenticate |
| 403 | Forbidden | Insufficient permissions |
| 404 | Not Found | Resource doesn't exist |
| 409 | Conflict | Resource already exists |
| 422 | Validation Error | Check request data format |
| 429 | Rate Limited | Too many requests - slow down |
| 500 | Server Error | Backend error - try again later |

### Error Response Format

```json
{
  "error": "validation_error",
  "message": "Invalid kanji character provided",
  "details": {
    "field": "kanji_char",
    "issue": "must_be_single_character"
  }
}
```

### Common Error Patterns

**Token Expired:**
```json
{
  "error": "unauthorized",
  "message": "Token has expired"
}
```
→ Re-authenticate with `/auth/login`

**Invalid Kanji:**
```json
{
  "error": "not_found",
  "message": "Kanji 日 not found in database"
}
```
→ Seed kanji data with `/kanji/seed` (admin only)

**Session Not Found:**
```json
{
  "error": "not_found",
  "message": "Conversation session not found"
}
```
→ Start new session with `/nichijou/chat/start`

---

## Rate Limits

- **Authenticated:** 100 requests per minute
- **Anonymous:** 20 requests per minute

Rate limit headers:
```
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1713781200
```

---

## WebSocket (Future)

Real-time features coming in Phase 3:
```
wss://46.224.127.221:8090/ws
```

Events:
- `peer_matched` - Found conversation partner
- `message_received` - New message from peer
- `typing_indicator` - Peer is typing
- `session_ended` - Partner left session

---

*Last updated: 2026-04-22*