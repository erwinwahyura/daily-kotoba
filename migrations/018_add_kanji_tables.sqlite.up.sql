-- Kanji tables for stroke order practice

-- Kanji characters with stroke data
CREATE TABLE IF NOT EXISTS kanji (
    id TEXT PRIMARY KEY,
    character TEXT UNIQUE NOT NULL,
    jlpt_level TEXT NOT NULL,
    meaning TEXT NOT NULL,
    readings JSONB NOT NULL DEFAULT '[]',
    stroke_count INTEGER NOT NULL,
    stroke_order JSONB NOT NULL DEFAULT '[]',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Kanji practice sessions
CREATE TABLE IF NOT EXISTS kanji_practice_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    kanji_id TEXT NOT NULL REFERENCES kanji(id) ON DELETE CASCADE,
    kanji_char TEXT NOT NULL,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    user_strokes JSONB NOT NULL DEFAULT '[]',
    accuracy REAL DEFAULT 0,
    status TEXT DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Index for faster queries
CREATE INDEX IF NOT EXISTS idx_kanji_level ON kanji(jlpt_level);
CREATE INDEX IF NOT EXISTS idx_kanji_sessions_user ON kanji_practice_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_kanji_sessions_kanji ON kanji_practice_sessions(kanji_id);
