-- Conjugation Drill System
CREATE TABLE IF NOT EXISTS conjugation_challenges (
    id TEXT PRIMARY KEY,
    base_form TEXT NOT NULL,
    reading TEXT NOT NULL,
    "group" TEXT NOT NULL,
    target_form TEXT NOT NULL,
    target_ending TEXT NOT NULL,
    full_answer TEXT NOT NULL,
    hint TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    jlpt_level TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_conj_form ON conjugation_challenges(target_form);
CREATE INDEX IF NOT EXISTS idx_conj_level ON conjugation_challenges(jlpt_level);

CREATE TABLE IF NOT EXISTS conjugation_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    current_form TEXT NOT NULL,
    current_index INTEGER NOT NULL DEFAULT 0,
    total_questions INTEGER NOT NULL DEFAULT 0,
    correct_count INTEGER NOT NULL DEFAULT 0,
    wrong_count INTEGER NOT NULL DEFAULT 0,
    streak INTEGER NOT NULL DEFAULT 0,
    max_streak INTEGER NOT NULL DEFAULT 0,
    start_time TIMESTAMP NOT NULL,
    last_active TIMESTAMP NOT NULL,
    completed_forms TEXT
);

CREATE INDEX IF NOT EXISTS idx_conj_session_user ON conjugation_sessions(user_id);

CREATE TABLE IF NOT EXISTS conjugation_attempts (
    id TEXT PRIMARY KEY,
    session_id TEXT REFERENCES conjugation_sessions(id) ON DELETE CASCADE,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    challenge_id TEXT NOT NULL,
    form_type TEXT NOT NULL,
    base_form TEXT NOT NULL,
    user_answer TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL,
    time_spent_sec INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_conj_attempt_user ON conjugation_attempts(user_id);
