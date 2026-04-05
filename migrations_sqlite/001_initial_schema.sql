-- SQLite schema migration from PostgreSQL
-- Run this once to initialize SQLite database

-- Enable foreign keys
PRAGMA foreign_keys = ON;

-- Users table
CREATE TABLE users (
    id TEXT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    current_jlpt_level TEXT NOT NULL DEFAULT 'N5',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);

-- Vocabulary table
CREATE TABLE vocabulary (
    id TEXT PRIMARY KEY,
    word TEXT NOT NULL,
    reading TEXT NOT NULL,
    short_meaning TEXT NOT NULL,
    detailed_explanation TEXT NOT NULL,
    example_sentences TEXT, -- JSON array
    usage_notes TEXT,
    jlpt_level TEXT NOT NULL,
    index_position INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    -- Enhanced fields
    related_words TEXT, -- JSON object
    word_type TEXT DEFAULT 'unknown',
    register TEXT DEFAULT 'neutral',
    common_mistakes TEXT DEFAULT ''
);

CREATE INDEX idx_vocab_level ON vocabulary(jlpt_level);
CREATE INDEX idx_vocab_level_index ON vocabulary(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_vocab_level_position_unique ON vocabulary(jlpt_level, index_position);
CREATE INDEX idx_vocab_word ON vocabulary(word);

-- Grammar patterns table
CREATE TABLE grammar_patterns (
    id TEXT PRIMARY KEY,
    pattern TEXT NOT NULL,
    plain_form TEXT NOT NULL,
    meaning TEXT NOT NULL,
    detailed_explanation TEXT NOT NULL,
    conjugation_rules TEXT NOT NULL,
    usage_examples TEXT, -- JSON array
    nuance_notes TEXT NOT NULL,
    jlpt_level TEXT NOT NULL,
    related_patterns TEXT, -- JSON array
    common_mistakes TEXT DEFAULT '',
    index_position INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_grammar_level ON grammar_patterns(jlpt_level);
CREATE INDEX idx_grammar_level_index ON grammar_patterns(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_grammar_pattern_unique ON grammar_patterns(jlpt_level, index_position);

-- User progress table
CREATE TABLE user_progress (
    user_id TEXT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_vocab_index INTEGER NOT NULL DEFAULT 0,
    current_grammar_index INTEGER NOT NULL DEFAULT 0,
    last_word_id TEXT REFERENCES vocabulary(id),
    last_grammar_id TEXT REFERENCES grammar_patterns(id),
    streak_days INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE,
    words_learned_count INTEGER NOT NULL DEFAULT 0,
    words_skipped_count INTEGER NOT NULL DEFAULT 0,
    grammar_learned_count INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- User vocabulary status
CREATE TABLE user_vocab_status (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocab_id TEXT NOT NULL REFERENCES vocabulary(id),
    status TEXT NOT NULL CHECK (status IN ('learning', 'known', 'skipped')),
    marked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, vocab_id)
);

CREATE INDEX idx_user_vocab_status_user ON user_vocab_status(user_id);
CREATE INDEX idx_user_vocab_status_vocab ON user_vocab_status(vocab_id);

-- User grammar status
CREATE TABLE user_grammar_status (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grammar_id TEXT NOT NULL REFERENCES grammar_patterns(id),
    status TEXT NOT NULL CHECK (status IN ('learning', 'mastered', 'skipped')),
    marked_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, grammar_id)
);

CREATE INDEX idx_user_grammar_status_user ON user_grammar_status(user_id);
CREATE INDEX idx_user_grammar_status_grammar ON user_grammar_status(grammar_id);

-- Placement test results
CREATE TABLE placement_test_results (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    test_score INTEGER NOT NULL,
    assigned_level TEXT NOT NULL,
    completed_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Placement questions
CREATE TABLE placement_questions (
    id TEXT PRIMARY KEY,
    question_text TEXT NOT NULL,
    correct_answer TEXT NOT NULL,
    wrong_answers TEXT, -- JSON array
    difficulty_level TEXT NOT NULL,
    order_index INTEGER NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_difficulty ON placement_questions(difficulty_level);
CREATE INDEX idx_placement_order ON placement_questions(order_index);

-- Trigger to auto-update updated_at
CREATE TRIGGER update_user_progress_updated_at 
AFTER UPDATE ON user_progress
BEGIN
    UPDATE user_progress SET updated_at = CURRENT_TIMESTAMP 
    WHERE user_id = NEW.user_id;
END;
