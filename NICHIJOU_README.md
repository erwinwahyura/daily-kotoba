# Nichijou (日常) Conversation Features

## Vision
Transform Kotoba from a study tool into a living Japanese conversation practice platform. Focus on **natural, everyday interactions** — the kind of small talk, casual exchanges, and practical scenarios that actually happen in Japan.

---

## Core Features

### 1. AI Conversation Partner (Solo Practice)

**Scenario-Based Chat Simulation**
- Text-based chat with AI conversation partner
- Pre-defined scenarios: ordering at restaurant, asking directions, konbini small talk, weather chat
- JLPT-level appropriate responses (N5 simple → N3 natural)
- Real-time grammar and vocabulary corrections

**Key Mechanics:**
- User sends message in Japanese
- AI responds naturally + provides correction/feedback panel
- Tracks "naturalness score" (0-100) based on appropriateness for context
- Suggests alternative phrasings (polite vs casual)

**Tech Notes:**
- Requires LLM integration (OpenAI/Claude API)
- Conversation state stored per session
- Correction logic: grammar check + context appropriateness

---

### 2. Voice Shadowing & Pronunciation

**Shadowing Practice**
- Native audio dialogue (TTS or recorded)
- User records their attempt
- Playback comparison: native → user → native
- Pitch accent visualization (optional v2)

**Pronunciation Check**
- Speech-to-text accuracy
- Rhythm/pause pattern matching
- Score on clarity and natural flow

**Tech Notes:**
- Web Audio API for recording
- ElevenLabs or similar for quality TTS
- Storage: temporary audio blobs (don't persist)

---

### 3. Peer Practice (Social)

**Quick Match System**
- 5-minute conversation slots
- Matched by JLPT level + interests
- Topic prompts provided ("今日の天気", "週末の予定")
- Text or voice (user choice)

**Safety/Moderation:**
- Report/block functionality
- Pre-built conversation starters to reduce awkwardness
- Optional "teacher mode" where one user is more advanced

**Tech Notes:**
- Requires WebSocket for real-time chat
- Matchmaking queue (Redis/simple in-memory)
- Optional: turn-based if async preferred

---

### 4. Micro-Interactions (Nichijou Vibe)

**30-Second Casual Exchanges**
- Designed for "in-between" moments (waiting for train, coffee break)
- Single exchange: greeting + response + closing
- Focus on **reaction patterns** and **filler words**

**Examples:**
```
AI: 「あ、雨降ってきたね」
User options: 
- 「うわ、本当だ。傘持ってない…」
- 「えー、やだなぁ」
- 「へえ、そうなんだ」
```

**Register Training:**
- です/ます ↔ だ/よね switching
- Filler words: えーと、なんか、まぁ
- Natural pauses and back-channeling (うん、へえ、そうなんだ)

---

## Implementation Phases

### Phase 1: Solo AI Chat (MVP)
- [ ] LLM integration endpoint
- [ ] Simple chat UI (text only)
- [ ] 5 starter scenarios
- [ ] Basic correction display

### Phase 2: Voice Features
- [ ] Audio recording component
- [ ] TTS integration for prompts
- [ ] Shadowing playback UI
- [ ] Speech-to-text check

### Phase 3: Social Features
- [ ] WebSocket server
- [ ] Matchmaking logic
- [ ] Peer chat UI
- [ ] Topic prompt system

### Phase 4: Polish
- [ ] Micro-interaction library (20+ exchanges)
- [ ] Register/style switching lessons
- [ ] Reaction training ("I didn't catch that" responses)
- [ ] Progress tracking across conversation modes

---

## Database Schema (Suggested)

```sql
-- Conversation sessions
CREATE TABLE conversation_sessions (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    mode VARCHAR(20), -- 'ai', 'peer', 'shadowing'
    scenario_id VARCHAR(50), -- 'restaurant', 'directions', etc.
    level INTEGER, -- JLPT level
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    naturalness_avg INTEGER -- 0-100 average score
);

-- Messages within sessions
CREATE TABLE conversation_messages (
    id UUID PRIMARY KEY,
    session_id UUID REFERENCES conversation_sessions(id),
    sender VARCHAR(10), -- 'user', 'ai', 'peer'
    content TEXT,
    correction TEXT, -- AI feedback on user message
    naturalness_score INTEGER, -- 0-100
    created_at TIMESTAMP
);

-- Shadowing attempts
CREATE TABLE shadowing_attempts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    prompt_id VARCHAR(50),
    audio_url TEXT, -- temporary storage
    accuracy_score INTEGER,
    created_at TIMESTAMP
);
```

---

## API Endpoints (Suggested)

```
POST /api/nichijou/chat/start      -- Start AI conversation session
POST /api/nichijou/chat/message    -- Send message, get AI response
GET  /api/nichijou/scenarios       -- List available scenarios
POST /api/nichijou/shadow/start    -- Start shadowing session
POST /api/nichijou/shadow/submit   -- Submit recording for analysis
POST /api/nichijou/match/join      -- Join peer matchmaking queue
POST /api/nichijou/match/cancel    -- Leave queue
```

---

## Design Notes

**Tone:** Casual, encouraging, not "test-like"
- AI should feel like a patient friend, not a teacher
- Corrections should be gentle suggestions, not red marks
- Celebrate small wins: 「いいね！自然な日本語だよ」

**Visual Style:**
- Chat bubble interface (LINE/WhatsApp feel)
- Soft colors, rounded corners
- Typing indicators for AI/peer
- Reaction buttons (👍, 😂, 🤔) for quick responses

---

## Success Metrics

- Session completion rate (target: 80%)
- User return rate for conversation features
- Naturalness score improvement over time
- Peer match satisfaction (post-chat rating)

---

## Future Extensions

- **Cultural context lessons:** Why お疲れ is tricky, when to use やばい
- **Dialect exposure:** Osaka-ben, Hokkaido-ben variations
- **Business keigo track:** Formal workplace scenarios
- **Anime integration:** Practice with actual dialogue clips

---

*Let's make Japanese feel alive, not just studied.* 🎋
