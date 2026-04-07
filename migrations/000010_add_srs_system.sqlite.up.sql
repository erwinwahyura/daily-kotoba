-- SRS (Spaced Repetition) System Tables for SQLite
-- Implements SM-2 algorithm for optimal retention

CREATE TABLE srs_schedules (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id TEXT NOT NULL,  -- Can reference vocabulary.id or grammar_patterns.id
    item_type TEXT NOT NULL CHECK (item_type IN ('vocabulary', 'grammar')),
    
    -- SM-2 Algorithm fields
    interval_days INTEGER DEFAULT 0,
    repetitions INTEGER DEFAULT 0,
    ease_factor REAL DEFAULT 2.5,
    
    -- Review tracking
    last_reviewed_at TIMESTAMP,
    next_review_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Performance stats
    total_reviews INTEGER DEFAULT 0,
    correct_reviews INTEGER DEFAULT 0,
    streak INTEGER DEFAULT 0,
    
    -- Status
    status TEXT DEFAULT 'learning' CHECK (status IN ('learning', 'review', 'mastered', 'lapsed')),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint
    UNIQUE(user_id, item_id, item_type)
);

CREATE INDEX idx_srs_user_next_review ON srs_schedules(user_id, next_review_at);
CREATE INDEX idx_srs_user_status ON srs_schedules(user_id, status);

-- Review history table
CREATE TABLE srs_review_history (
    id TEXT PRIMARY KEY,
    schedule_id TEXT NOT NULL REFERENCES srs_schedules(id) ON DELETE CASCADE,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    quality INTEGER NOT NULL CHECK (quality BETWEEN 0 AND 5),
    response_time_ms INTEGER,
    
    item_type TEXT NOT NULL,
    item_id TEXT NOT NULL,
    
    interval_before INTEGER,
    interval_after INTEGER,
    ease_factor_before REAL,
    ease_factor_after REAL,
    
    reviewed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_srs_history_user ON srs_review_history(user_id, reviewed_at DESC);
CREATE INDEX idx_srs_history_schedule ON srs_review_history(schedule_id, reviewed_at DESC);

-- Daily study stats
CREATE TABLE study_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_date TEXT NOT NULL,  -- YYYY-MM-DD format
    
    new_items INTEGER DEFAULT 0,
    review_items INTEGER DEFAULT 0,
    drill_items INTEGER DEFAULT 0,
    
    total_time_seconds INTEGER DEFAULT 0,
    
    correct_count INTEGER DEFAULT 0,
    total_attempts INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, session_date)
);

CREATE INDEX idx_study_sessions_user_date ON study_sessions(user_id, session_date DESC);
