# Claude Code Prompt: Kotoba Backend API (kotoba-api)

## Project Context
This is the backend API for Kotoba, a Japanese vocabulary learning app that uses home screen widgets for passive vocabulary reinforcement. The full project context is in `vocab-widget-app-context.md` - please read it first.

## Repository: kotoba-api
Go backend API that serves vocabulary data, manages user authentication, tracks progress, and powers the iOS app.

---

## Development Goals

### MVP Scope
Build a fully functional REST API that:
- Handles user authentication (register, login, JWT)
- Serves vocabulary data organized by JLPT levels
- Implements vocabulary rotation logic (sequential → review cycles)
- Tracks user progress and learning statistics
- Provides placement test functionality
- Supports the "skip ahead" feature

### Non-Goals for MVP
- ❌ Admin dashboard
- ❌ Social features
- ❌ Advanced analytics
- ❌ Multiple language support (Japanese only for now)
- ❌ Email verification
- ❌ Password reset flow

---

## Phase 1: Project Setup & Foundation

### Task 1.1: Initialize Go Project
```bash
# Create project structure
kotoba-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── handlers/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── services/
│   └── utils/
├── migrations/
├── pkg/
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
├── Dockerfile
└── README.md
```

**Tasks:**
1. Initialize Go modules: `go mod init github.com/erwin/kotoba-api`
2. Create directory structure following best practices
3. Set up .gitignore for Go projects
4. Create .env.example with configuration variables:
   ```
   DATABASE_URL=postgresql://user:pass@localhost:5432/kotoba_db
   JWT_SECRET=your-secret-key-change-in-production
   PORT=8080
   ENV=development
   ```
5. Create README.md with:
   - Project description
   - Setup instructions
   - API endpoint documentation
   - Development workflow

### Task 1.2: Dependencies
Install and configure required Go packages:
```bash
# Core dependencies
go get github.com/gin-gonic/gin                    # Web framework
go get github.com/lib/pq                           # PostgreSQL driver
go get github.com/golang-jwt/jwt/v5                # JWT tokens
go get golang.org/x/crypto/bcrypt                  # Password hashing
go get github.com/joho/godotenv                    # Environment variables
go get github.com/golang-migrate/migrate/v4        # Database migrations

# Testing
go get github.com/stretchr/testify                 # Test assertions
```

### Task 1.3: Database Setup
1. Choose database: **PostgreSQL** (recommended) or SQLite for development
2. Create database schema (see next section)
3. Set up migration tool (golang-migrate)
4. Create initial migrations

---

## Phase 2: Database Schema & Migrations

### Task 2.1: Database Schema Design

**Users Table:**
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    current_jlpt_level VARCHAR(10) NOT NULL DEFAULT 'N5', -- N5, N4, N3, N2, N1
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
```

**Vocabulary Table:**
```sql
CREATE TABLE vocabulary (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    word VARCHAR(100) NOT NULL,                    -- 諦める
    reading VARCHAR(100) NOT NULL,                  -- あきらめる
    short_meaning VARCHAR(255) NOT NULL,            -- to give up
    detailed_explanation TEXT NOT NULL,
    example_sentences JSONB,                        -- Array of example sentences
    usage_notes TEXT,
    jlpt_level VARCHAR(10) NOT NULL,                -- N5, N4, N3, N2, N1
    index_position INTEGER NOT NULL,                -- 0-199 for N4, etc.
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vocab_level ON vocabulary(jlpt_level);
CREATE INDEX idx_vocab_level_index ON vocabulary(jlpt_level, index_position);
```

**User Progress Table:**
```sql
CREATE TABLE user_progress (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_vocab_index INTEGER NOT NULL DEFAULT 0,
    last_word_id UUID REFERENCES vocabulary(id),
    streak_days INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE,
    words_learned_count INTEGER NOT NULL DEFAULT 0,
    words_skipped_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**User Vocabulary Status Table:**
```sql
CREATE TABLE user_vocab_status (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocab_id UUID NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,                    -- 'learning', 'known', 'skipped'
    marked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, vocab_id)
);

CREATE INDEX idx_user_vocab_status ON user_vocab_status(user_id, status);
```

**Placement Test Results Table:**
```sql
CREATE TABLE placement_test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    test_score INTEGER NOT NULL,
    assigned_level VARCHAR(10) NOT NULL,
    completed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_user ON placement_test_results(user_id);
```

**Placement Test Questions Table (optional - can be hardcoded):**
```sql
CREATE TABLE placement_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_text TEXT NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    wrong_answers JSONB NOT NULL,                   -- Array of 3 wrong answers
    difficulty_level VARCHAR(10) NOT NULL,          -- N5, N4, N3, N2, N1
    order_index INTEGER NOT NULL
);
```

### Task 2.2: Create Migration Files
Create migration files in `migrations/` directory:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`
- `000002_create_vocabulary_table.up.sql`
- `000002_create_vocabulary_table.down.sql`
- etc.

### Task 2.3: Migration Commands
Add migration scripts to README:
```bash
# Run migrations
migrate -path migrations -database "postgresql://localhost:5432/kotoba_db" up

# Rollback migrations
migrate -path migrations -database "postgresql://localhost:5432/kotoba_db" down
```

---

## Phase 3: Core API Implementation

### Task 3.1: Project Structure

**Models (internal/models/):**
```go
// user.go
type User struct {
    ID              string    `json:"id"`
    Email           string    `json:"email"`
    PasswordHash    string    `json:"-"`
    CurrentLevel    string    `json:"current_jlpt_level"`
    CreatedAt       time.Time `json:"created_at"`
    UpdatedAt       time.Time `json:"updated_at"`
}

// vocabulary.go
type Vocabulary struct {
    ID                  string          `json:"id"`
    Word                string          `json:"word"`
    Reading             string          `json:"reading"`
    ShortMeaning        string          `json:"short_meaning"`
    DetailedExplanation string          `json:"detailed_explanation"`
    ExampleSentences    []string        `json:"example_sentences"`
    UsageNotes          string          `json:"usage_notes"`
    JLPTLevel           string          `json:"jlpt_level"`
    IndexPosition       int             `json:"index_position"`
    CreatedAt           time.Time       `json:"created_at"`
}

// progress.go
type UserProgress struct {
    UserID            string    `json:"user_id"`
    CurrentVocabIndex int       `json:"current_vocab_index"`
    LastWordID        string    `json:"last_word_id"`
    StreakDays        int       `json:"streak_days"`
    LastStudyDate     time.Time `json:"last_study_date"`
    WordsLearnedCount int       `json:"words_learned_count"`
    WordsSkippedCount int       `json:"words_skipped_count"`
    UpdatedAt         time.Time `json:"updated_at"`
}

// Add other models as needed
```

### Task 3.2: Authentication Endpoints

**POST /api/auth/register**
```go
// Request body
{
    "email": "user@example.com",
    "password": "securepassword123"
}

// Response
{
    "user": {
        "id": "uuid",
        "email": "user@example.com",
        "current_jlpt_level": "N5"
    },
    "token": "jwt-token-here"
}
```

**POST /api/auth/login**
```go
// Request body
{
    "email": "user@example.com",
    "password": "securepassword123"
}

// Response
{
    "user": { /* user object */ },
    "token": "jwt-token-here"
}
```

**GET /api/auth/me** (Protected)
```go
// Headers: Authorization: Bearer <jwt-token>

// Response
{
    "user": { /* user object */ }
}
```

**Implementation tasks:**
1. Create auth handlers in `internal/handlers/auth.go`
2. Implement password hashing with bcrypt
3. Implement JWT token generation and validation
4. Create auth middleware for protected routes
5. Add input validation (email format, password length)
6. Error handling (user exists, invalid credentials, etc.)

### Task 3.3: Vocabulary Endpoints

**GET /api/vocab/daily** (Protected)
```go
// Returns the current word for the user based on their progress

// Response
{
    "vocabulary": {
        "id": "uuid",
        "word": "諦める",
        "reading": "あきらめる",
        "short_meaning": "to give up",
        "detailed_explanation": "...",
        "example_sentences": ["..."],
        "usage_notes": "...",
        "jlpt_level": "N4",
        "index_position": 45
    },
    "progress": {
        "current_index": 45,
        "total_words_in_level": 200,
        "words_learned": 44,
        "streak_days": 5
    }
}
```

**GET /api/vocab/:id** (Protected)
```go
// Get detailed information about a specific vocabulary word

// Response
{
    "vocabulary": { /* full vocab object */ }
}
```

**POST /api/vocab/:id/skip** (Protected)
```go
// Mark word as known and get next word

// Request body (optional)
{
    "status": "known"  // or "skipped"
}

// Response
{
    "next_vocabulary": { /* next vocab object */ },
    "progress": { /* updated progress */ }
}
```

**GET /api/vocab/level/:level** (Protected)
```go
// Get all vocabulary for a specific JLPT level (paginated)
// Example: /api/vocab/level/N4?page=1&limit=20

// Response
{
    "vocabulary": [ /* array of vocab objects */ ],
    "pagination": {
        "page": 1,
        "limit": 20,
        "total": 200,
        "total_pages": 10
    }
}
```

**Implementation tasks:**
1. Create vocab handlers in `internal/handlers/vocabulary.go`
2. Implement vocabulary rotation logic (see Task 3.6)
3. Handle edge cases (end of level, user level change)
4. Add caching for frequently accessed words (optional)

### Task 3.4: Progress Endpoints

**GET /api/progress** (Protected)
```go
// Get user's current progress

// Response
{
    "progress": {
        "current_vocab_index": 45,
        "current_level": "N4",
        "streak_days": 5,
        "last_study_date": "2026-01-17",
        "words_learned": 44,
        "words_skipped": 3,
        "total_words_in_level": 200
    }
}
```

**GET /api/progress/stats** (Protected)
```go
// Get detailed statistics

// Response
{
    "stats": {
        "total_days_active": 45,
        "current_streak": 5,
        "longest_streak": 12,
        "words_learned_by_level": {
            "N5": 150,
            "N4": 44
        },
        "average_words_per_day": 4.3,
        "total_words_learned": 194
    }
}
```

### Task 3.5: Placement Test Endpoints

**GET /api/placement-test**
```go
// Get placement test questions (20 questions)

// Response
{
    "questions": [
        {
            "id": "uuid",
            "question": "Choose the correct reading for 食べる",
            "options": ["たべる", "のべる", "かべる", "よべる"],
            "difficulty": "N5"
        },
        // ... 19 more questions
    ]
}
```

**POST /api/placement-test/submit** (Protected)
```go
// Submit test answers and get assigned level

// Request body
{
    "answers": {
        "question_id_1": "たべる",
        "question_id_2": "...",
        // ... 20 answers
    }
}

// Response
{
    "score": 16,
    "total_questions": 20,
    "assigned_level": "N4",
    "breakdown": {
        "N5": 5,  // correct answers at N5 level
        "N4": 8,
        "N3": 3,
        "N2": 0,
        "N1": 0
    }
}
```

**Implementation tasks:**
1. Create placement test questions (can hardcode 20 questions for MVP)
2. Implement scoring algorithm:
   - 15-20 correct → N3
   - 10-14 correct → N4
   - 5-9 correct → N5
   - <5 correct → N5
3. Update user's current_jlpt_level after test
4. Initialize user_progress with correct starting index

### Task 3.6: Vocabulary Rotation Logic

**Core Business Logic (internal/services/vocab_service.go):**

```go
// GetDailyWord returns the current word for a user
func (s *VocabService) GetDailyWord(userID string) (*Vocabulary, error) {
    // 1. Get user's progress
    progress, err := s.repo.GetUserProgress(userID)
    if err != nil {
        return nil, err
    }
    
    // 2. Check if it's a new day (update streak)
    s.updateStreak(userID, progress)
    
    // 3. Get vocabulary at current index for user's level
    vocab, err := s.repo.GetVocabByLevelAndIndex(
        progress.CurrentLevel, 
        progress.CurrentVocabIndex,
    )
    if err != nil {
        return nil, err
    }
    
    return vocab, nil
}

// SkipToNextWord marks current word as known/skipped and returns next word
func (s *VocabService) SkipToNextWord(userID, vocabID, status string) (*Vocabulary, error) {
    // 1. Mark current word with status
    err := s.repo.MarkVocabStatus(userID, vocabID, status)
    if err != nil {
        return nil, err
    }
    
    // 2. Increment user's vocab index
    progress, err := s.repo.IncrementVocabIndex(userID)
    if err != nil {
        return nil, err
    }
    
    // 3. Check if we've reached end of level (index >= total words)
    // If yes, cycle back to 0 (Phase 1 logic)
    if progress.CurrentVocabIndex >= s.getTotalWordsForLevel(progress.CurrentLevel) {
        progress.CurrentVocabIndex = 0
        s.repo.UpdateUserProgress(progress)
    }
    
    // 4. Get next vocabulary word
    nextVocab, err := s.repo.GetVocabByLevelAndIndex(
        progress.CurrentLevel,
        progress.CurrentVocabIndex,
    )
    if err != nil {
        return nil, err
    }
    
    return nextVocab, nil
}

// Helper: Get total words for a JLPT level
func (s *VocabService) getTotalWordsForLevel(level string) int {
    // Can be configured or queried from DB
    counts := map[string]int{
        "N5": 800,
        "N4": 300,   // Start with smaller set for MVP
        "N3": 650,
        "N2": 1000,
        "N1": 2000,
    }
    return counts[level]
}
```

**Phase 2 Logic (Post-MVP):**
After completing one cycle through all words:
- Mix 80% new words (from same level or next level) + 20% review words
- Implement spaced repetition (words user marked as "learning")

---

## Phase 4: Middleware & Utilities

### Task 4.1: JWT Middleware
```go
// internal/middleware/auth.go
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extract token from Authorization header
        // 2. Validate JWT token
        // 3. Extract user ID from token claims
        // 4. Set user ID in context
        // 5. Call next handler or abort with 401
    }
}
```

### Task 4.2: CORS Middleware
```go
// Allow iOS app to make requests
func CORSMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
        
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }
        
        c.Next()
    }
}
```

### Task 4.3: Error Handling
```go
// internal/utils/errors.go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func HandleError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, APIError{
        Code:    statusCode,
        Message: message,
    })
}
```

### Task 4.4: Logging
Set up structured logging:
```go
// Use standard log package or zerolog/logrus
log.Printf("[INFO] User %s fetched daily word", userID)
log.Printf("[ERROR] Failed to fetch vocabulary: %v", err)
```

---

## Phase 5: Seed Data & Testing

### Task 5.1: Vocabulary Seed Data
Create seed script (`cmd/seed/main.go`) that populates:

**Minimum for MVP:**
- 50 N4 vocabulary entries (can expand later)
- Example format:
```go
var n4Vocab = []Vocabulary{
    {
        Word: "諦める",
        Reading: "あきらめる",
        ShortMeaning: "to give up",
        DetailedExplanation: "Used when abandoning an effort or goal. Often implies acceptance of failure or impossibility.",
        ExampleSentences: []string{
            "試験に合格するのを諦めた - I gave up on passing the exam",
            "彼女は夢を諦めなかった - She didn't give up on her dream",
        },
        UsageNotes: "Commonly used in both casual and formal contexts",
        JLPTLevel: "N4",
        IndexPosition: 0,
    },
    // ... 49 more entries
}
```

**Where to get vocabulary data:**
- JLPT official word lists
- Jisho.org API (if available)
- WaniKani API (requires subscription)
- Manually curate from textbooks (Genki, Minna no Nihongo)

### Task 5.2: Placement Test Questions
Create 20 test questions:
- 5 N5 level questions
- 5 N4 level questions
- 5 N3 level questions
- 3 N2 level questions
- 2 N1 level questions

Example:
```go
var placementQuestions = []PlacementQuestion{
    {
        QuestionText: "Choose the correct reading for 食べる",
        CorrectAnswer: "たべる",
        WrongAnswers: []string{"のべる", "かべる", "よべる"},
        DifficultyLevel: "N5",
        OrderIndex: 1,
    },
    // ... 19 more
}
```

### Task 5.3: Unit Tests
Write tests for critical functions:

```go
// internal/services/vocab_service_test.go
func TestGetDailyWord(t *testing.T) {
    // Test: Returns correct word for user's current index
}

func TestSkipToNextWord(t *testing.T) {
    // Test: Increments index correctly
    // Test: Cycles back to 0 at end of level
    // Test: Marks vocab status correctly
}

// internal/handlers/auth_test.go
func TestRegister(t *testing.T) {
    // Test: Creates user successfully
    // Test: Rejects duplicate email
    // Test: Validates password requirements
}
```

Run tests:
```bash
go test ./...
```

### Task 5.4: API Testing
Create example curl commands in README:

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Get daily word (replace <token>)
curl -X GET http://localhost:8080/api/vocab/daily \
  -H "Authorization: Bearer <token>"
```

---

## Phase 6: Deployment Preparation

### Task 6.1: Environment Configuration
Create production-ready config:
```go
// internal/config/config.go
type Config struct {
    DatabaseURL string
    JWTSecret   string
    Port        string
    Environment string
}

func Load() (*Config, error) {
    // Load from environment variables
    // Validate required fields
}
```

### Task 6.2: Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd/api

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .
EXPOSE 8080

CMD ["./main"]
```

### Task 6.3: Health Check Endpoint
```go
// GET /health
func HealthCheck(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
        "timestamp": time.Now(),
    })
}
```

### Task 6.4: Database Connection Pooling
```go
db, err := sql.Open("postgres", databaseURL)
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(5)
db.SetConnMaxLifetime(5 * time.Minute)
```

### Task 6.5: Deployment Options
Choose one:
1. **Railway.app** - Easy deployment, PostgreSQL included
2. **Fly.io** - Global edge deployment
3. **Render.com** - Free tier available
4. **VPS** (Hetzner) - More control

probably choose hetzner because im already take it.

For Railway:
```bash
# Install Railway CLI
npm i -g @railway/cli

# Login and deploy
railway login
railway init
railway up
```

---

## Phase 7: Documentation & Final Touches

### Task 7.1: API Documentation
Create comprehensive API docs in `API.md`:
- All endpoints with request/response examples
- Authentication flow
- Error codes and messages
- Rate limiting (if implemented)

### Task 7.2: README Updates
Update README with:
- Project setup instructions
- Environment variables
- Running migrations
- Seeding database
- Running tests
- Deployment guide
- Contributing guidelines (if open source)

### Task 7.3: Security Checklist
- [ ] JWT tokens have reasonable expiration (24 hours)
- [ ] Passwords hashed with bcrypt (cost factor 10+)
- [ ] SQL queries use parameterized statements
- [ ] HTTPS enforced in production
- [ ] Sensitive data not logged
- [ ] Environment variables for secrets
- [ ] CORS configured correctly
- [ ] Input validation on all endpoints

### Task 7.4: Performance Optimization
- [ ] Database indexes on frequently queried columns
- [ ] Connection pooling configured
- [ ] Caching for frequently accessed data (optional)
- [ ] Query optimization (avoid N+1 queries)

---

## Success Criteria

Backend is "deploy-ready" when:
- ✅ All core endpoints functional and tested
- ✅ User authentication works (register, login, JWT)
- ✅ Vocabulary rotation logic implemented
- ✅ Database migrations run successfully
- ✅ At least 50 N4 vocabulary words seeded
- ✅ Placement test questions available
- ✅ API documented with examples
- ✅ Deployed to production environment
- ✅ Health check endpoint responds
- ✅ iOS app can successfully connect and authenticate

---

## Next Steps After MVP

1. **Expand vocabulary database**:
   - Complete N4 (300 words)
   - Add N5, N3, N2, N1 levels

2. **Implement Phase 2 rotation**:
   - 80% new + 20% review algorithm
   - Spaced repetition system

3. **Add features**:
   - Admin API for adding/editing vocabulary
   - User settings (notification preferences, etc.)
   - Analytics endpoints (aggregate user statistics)

4. **Optimization**:
   - Caching layer (Redis)
   - Database query optimization
   - API rate limiting

---

## Development Workflow

1. **Start with database schema and migrations**
2. **Build auth system (register, login, JWT)**
3. **Test auth thoroughly before moving on**
4. **Implement vocabulary endpoints and rotation logic**
5. **Create seed data (50 N4 words minimum)**
6. **Test all endpoints with curl/Postman**
7. **Write unit tests for critical functions**
8. **Deploy to production**
9. **Verify iOS app can connect successfully**

Let's build the Kotoba API! 🚀