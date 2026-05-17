-- Nichijou Conversation Features - Phase 1 (SQLite)
-- AI Chat and Shadowing Tables

-- Conversation sessions
CREATE TABLE IF NOT EXISTS conversation_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mode TEXT NOT NULL CHECK (mode IN ('ai', 'peer', 'shadowing', 'micro')),
    scenario_id TEXT,
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 5),
    status TEXT DEFAULT 'active' CHECK (status IN ('active', 'completed', 'abandoned')),
    naturalness_avg INTEGER CHECK (naturalness_avg BETWEEN 0 AND 100),
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    ended_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Messages within sessions
CREATE TABLE IF NOT EXISTS conversation_messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES conversation_sessions(id) ON DELETE CASCADE,
    sender TEXT NOT NULL CHECK (sender IN ('user', 'ai', 'peer', 'system')),
    content TEXT NOT NULL,
    correction TEXT,
    naturalness_score INTEGER CHECK (naturalness_score BETWEEN 0 AND 100),
    alternatives TEXT, -- JSON array of alternative phrasings
    metadata TEXT, -- JSON additional data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Shadowing attempts
CREATE TABLE IF NOT EXISTS shadowing_attempts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt_id TEXT NOT NULL,
    scenario_id TEXT,
    native_audio_url TEXT,
    user_audio_url TEXT,
    transcript TEXT,
    accuracy_score INTEGER CHECK (accuracy_score BETWEEN 0 AND 100),
    rhythm_score INTEGER CHECK (rhythm_score BETWEEN 0 AND 100),
    pitch_data TEXT, -- JSON
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Scenarios for AI chat
CREATE TABLE IF NOT EXISTS conversation_scenarios (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 5),
    category TEXT NOT NULL, -- 'daily', 'travel', 'work', 'social'
    icon TEXT,
    system_prompt TEXT NOT NULL,
    starter_messages TEXT, -- JSON array
    vocabulary_hints TEXT, -- JSON
    grammar_points TEXT, -- JSON
    is_active INTEGER DEFAULT 1,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Micro-interactions library
CREATE TABLE IF NOT EXISTS micro_interactions (
    id TEXT PRIMARY KEY,
    prompt TEXT NOT NULL,
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 3),
    category TEXT NOT NULL, -- 'reaction', 'filler', 'small_talk'
    options TEXT NOT NULL, -- JSON array of response options
    best_answer TEXT,
    explanation TEXT,
    context_notes TEXT,
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_conversation_sessions_user_id ON conversation_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_conversation_sessions_status ON conversation_sessions(status);
CREATE INDEX IF NOT EXISTS idx_conversation_messages_session_id ON conversation_messages(session_id);
CREATE INDEX IF NOT EXISTS idx_conversation_messages_created_at ON conversation_messages(created_at);
CREATE INDEX IF NOT EXISTS idx_shadowing_attempts_user_id ON shadowing_attempts(user_id);
CREATE INDEX IF NOT EXISTS idx_conversation_scenarios_level ON conversation_scenarios(level);
CREATE INDEX IF NOT EXISTS idx_conversation_scenarios_category ON conversation_scenarios(category);
CREATE INDEX IF NOT EXISTS idx_micro_interactions_level ON micro_interactions(level);
