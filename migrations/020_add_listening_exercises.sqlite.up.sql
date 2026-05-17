-- Listening practice tables

CREATE TABLE IF NOT EXISTS listening_exercises (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    jlpt_level TEXT NOT NULL,
    difficulty TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    audio_url TEXT NOT NULL,
    duration INTEGER NOT NULL,
    transcript TEXT NOT NULL,
    translation TEXT,
    vocabulary JSONB NOT NULL DEFAULT '[]',
    questions JSONB NOT NULL DEFAULT '[]',
    topic TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS listening_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id TEXT NOT NULL REFERENCES listening_exercises(id) ON DELETE CASCADE,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    current_position INTEGER DEFAULT 0,
    answers JSONB NOT NULL DEFAULT '[]',
    score INTEGER DEFAULT 0,
    status TEXT DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed', 'abandoned')),
    play_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_listening_level ON listening_exercises(jlpt_level);
CREATE INDEX IF NOT EXISTS idx_listening_difficulty ON listening_exercises(difficulty);
CREATE INDEX IF NOT EXISTS idx_listening_topic ON listening_exercises(topic);
CREATE INDEX IF NOT EXISTS idx_listening_sessions_user ON listening_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_listening_sessions_exercise ON listening_sessions(exercise_id);
