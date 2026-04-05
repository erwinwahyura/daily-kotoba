
---

## Grammar Pattern Endpoints

### GET /api/grammar/daily

Get the current grammar pattern for the user's level.

**Authentication:** Required

**Response (200):**
```json
{
  "code": 200,
  "message": "Grammar pattern retrieved successfully",
  "data": {
    "pattern": {
      "id": "uuid",
      "pattern": "〜わけにはいかない",
      "plain_form": "わけにはいかない",
      "meaning": "cannot afford to; must not",
      "detailed_explanation": "Used when expressing that something cannot be done due to circumstances, obligations, or moral constraints. Stronger than 〜ないわけにはいかない (which is incorrect grammar).",
      "conjugation_rules": "Verb (dictionary form) + わけにはいかない",
      "usage_examples": [
        {
          "japanese": "仕事があるから、遊びに行くわけにはいかない",
          "reading": "しごとがあるから、あそびにいくわけにはいかない",
          "meaning": "I have work, so I cannot afford to go out and play",
          "nuance": "Expresses an obligation preventing the action",
          "context": "Telling friends you can't hang out",
          "alternative": "〜なければならない (if expressing obligation instead)"
        }
      ],
      "nuance_notes": "This pattern implies external circumstances make an action impossible. It's often used for moral obligations or unavoidable situations. Compare with 〜わけがない (completely different meaning!).",
      "jlpt_level": "N3",
      "related_patterns": [
        {
          "pattern": "〜わけがない",
          "relationship": "often confused with",
          "key_difference": "〜わけにはいかない = cannot afford to; 〜わけがない = no way that"
        }
      ],
      "common_mistakes": "Don't confuse with 〜わけがない. Also note: you CANNOT use this with adjectives directly (must use 〜くて/〜で form first).",
      "index_position": 0,
      "created_at": "2026-01-18T00:00:00Z"
    },
    "progress": {
      "current_index": 0,
      "total_patterns": 50,
      "patterns_learned": 0
    }
  }
}
```

**Errors:**
- `401`: Unauthorized
- `500`: Server error

---

### GET /api/grammar/:id

Get a specific grammar pattern by ID.

**Authentication:** Required

**Parameters:**
- `id` (path): Grammar pattern UUID

**Response (200):**
```json
{
  "code": 200,
  "message": "Grammar pattern retrieved successfully",
  "data": {
    "pattern": {
      "id": "uuid",
      "pattern": "〜ものだ",
      "plain_form": "ものだ",
      "meaning": "should; ought to; express emotion or general truth",
      "detailed_explanation": "Used to express: (1) emotional response to something, (2) general truths or common sense, (3) advice/suggestions (soft).",
      "conjugation_rules": "Verb (plain past/present) + ものだ / もんだ (casual)",
      "usage_examples": [
        {
          "japanese": "時間が経つのは早いものだ",
          "reading": "じかんがたつのははやいものだ",
          "meaning": "Time really does fly (emotional reflection)",
          "nuance": "Expresses sentiment about a universal truth",
          "context": "Reminiscing with an old friend",
          "alternative": "〜ね (less emotional)"
        }
      ],
      "nuance_notes": "Contrast with 〜べきだ (stronger obligation/advice) and 〜わけだ (explanation). ものだ is softer, more reflective.",
      "jlpt_level": "N3",
      "related_patterns": [
        {
          "pattern": "〜ものではない",
          "relationship": "similar form, opposite meaning",
          "key_difference": "〜ものだ = general truth/emotion; 〜ものではない = should never (strong prohibition)"
        }
      ],
      "common_mistakes": "Often confused with 〜ものではない (prohibition). Also: 〜もんだ is casual contraction, not a different grammar point.",
      "index_position": 5,
      "created_at": "2026-01-18T00:00:00Z"
    }
  }
}
```

**Errors:**
- `401`: Unauthorized
- `404`: Grammar pattern not found

---

### GET /api/grammar/level/:level

Get grammar patterns for a specific JLPT level.

**Authentication:** Required

**Parameters:**
- `level` (path): JLPT level (N5, N4, N3, N2, N1)

**Query Parameters:**
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20, max: 100)

**Example:** `/api/grammar/level/N3?page=1&limit=10`

**Response (200):**
```json
{
  "code": 200,
  "message": "Grammar patterns retrieved successfully",
  "data": {
    "patterns": [
      {
        "id": "uuid",
        "pattern": "〜ばかり",
        "plain_form": "ばかり",
        "meaning": "only; nothing but; approximately",
        "detailed_explanation": "Multiple uses: (1) 〜ばかり = only doing (negative connotation), (2) 〜ばかりだ = continuously happening, (3) 〜ばかりか = not only...but also.",
        "conjugation_rules": "Verb (te-form) + いる + ばかり; Noun + ばかり; 〜ばかりか〜も",
        "usage_examples": [
          {
            "japanese": "彼はゲームをしているばかりだ",
            "reading": "かれはゲームをしているばかりだ",
            "meaning": "All he does is play games (negative implication)",
            "nuance": "Criticizes excessive focus on one activity",
            "context": "Complaining about someone's habits",
            "alternative": "〜だけ (neutral 'only')"
          }
        ],
        "nuance_notes": "Unlike 〜だけ (neutral), 〜ばかり often carries negative judgment. Also note: when used with time, means 'approximately' (3時ばかり = around 3 o'clock).",
        "jlpt_level": "N3",
        "related_patterns": [
          {
            "pattern": "〜だけ",
            "relationship": "similar meaning, different nuance",
            "key_difference": "〜だけ = neutral 'only'; 〜ばかり = 'nothing but' (often negative)"
          },
          {
            "pattern": "〜のみ",
            "relationship": "formal equivalent",
            "key_difference": "〜のみ = formal/written 'only'; 〜ばかり = spoken, often critical"
          }
        ],
        "common_mistakes": "Don't use 〜ばかり when you mean neutral 'only' (use 〜だけ). Also: position matters - ばかり食べる (only eating) vs 食べるばかり (eating excessively/continuously).",
        "index_position": 10,
        "created_at": "2026-01-18T00:00:00Z"
      }
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

## Vocabulary Enhancements

Vocabulary responses now include enhanced fields:

- `related_words`: Synonyms, antonyms, and confusable words
- `word_type`: verb, noun, i-adjective, na-adjective, adverb, etc.
- `register`: formal, casual, neutral, slang
- `common_mistakes`: Frequent learner errors

**Example enhanced vocabulary entry:**
```json
{
  "id": "uuid",
  "word": "諦める",
  "reading": "あきらめる",
  "short_meaning": "to give up",
  "detailed_explanation": "Used when abandoning an effort or goal. Often implies acceptance of failure or impossibility. This word carries a sense of finality and resignation.",
  "example_sentences": [...],
  "usage_notes": "Commonly used in both casual and formal contexts. Often paired with のを when followed by a noun phrase.",
  "jlpt_level": "N4",
  "index_position": 0,
  "word_type": "verb (ichidan)",
  "register": "neutral",
  "common_mistakes": "Often confused with 辞める (やめる - to quit a job). 諦める is about accepting defeat; 辞める is about stopping an action/position.",
  "related_words": {
    "synonyms": ["放弃する", "止める"],
    "antonyms": ["頑張る", "諦めない"],
    "confusable": ["辞める (やめる - quit job)"],
    "see_also": ["断念する (formal)"]
  },
  "created_at": "2026-01-18T00:00:00Z"
}
```

---

## Error Handling

All endpoints return consistent error responses:

### Error Response Format
```json
{
  "code": 400,
  "message": "Error description"
}
```

### Common HTTP Status Codes

- `200`: Success
- `201`: Created (registration)
- `400`: Bad request (invalid input)
- `401`: Unauthorized (missing/invalid token)
- `404`: Not found
- `500`: Internal server error

---

## Response Format

### Success Response Structure
```json
{
  "code": 200,
  "message": "Human-readable success message",
  "data": { ... }
}
```

### Error Response Structure
```json
{
  "code": 400,
  "message": "Human-readable error message"
}
```

All timestamps are in ISO 8601 format (UTC).
