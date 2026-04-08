-- JLPT Mock Test System
CREATE TABLE IF NOT EXISTS jlpt_tests (
    id TEXT PRIMARY KEY,
    level TEXT NOT NULL,
    section TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    time_limit_minutes INTEGER NOT NULL,
    total_questions INTEGER NOT NULL,
    passing_score INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_jlpt_tests_level ON jlpt_tests(level);

CREATE TABLE IF NOT EXISTS jlpt_questions (
    id TEXT PRIMARY KEY,
    test_id TEXT REFERENCES jlpt_tests(id) ON DELETE CASCADE,
    question_num INTEGER NOT NULL,
    type TEXT NOT NULL,
    question TEXT NOT NULL,
    question_reading TEXT,
    english_prompt TEXT,
    options TEXT NOT NULL, -- JSON array
    correct_index INTEGER NOT NULL,
    explanation TEXT NOT NULL,
    point_value INTEGER DEFAULT 1,
    skill_tested TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_jlpt_questions_test ON jlpt_questions(test_id);

CREATE TABLE IF NOT EXISTS user_test_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    test_id TEXT REFERENCES jlpt_tests(id) ON DELETE CASCADE,
    level TEXT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    time_spent_sec INTEGER DEFAULT 0,
    answers TEXT, -- JSON map: question_id -> answer_index
    score INTEGER DEFAULT 0,
    correct_count INTEGER DEFAULT 0,
    status TEXT DEFAULT 'in_progress'
);

CREATE INDEX IF NOT EXISTS idx_user_tests_user ON user_test_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_tests_status ON user_test_sessions(status);