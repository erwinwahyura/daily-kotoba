-- Listening practice system

CREATE TABLE IF NOT EXISTS listening_exercises (
    id TEXT PRIMARY KEY,
    title TEXT NOT NULL,
    jlpt_level TEXT NOT NULL,
    difficulty TEXT NOT NULL DEFAULT 'medium', -- easy, medium, hard
    topic TEXT NOT NULL DEFAULT 'general',
    transcript TEXT NOT NULL,
    translation TEXT NOT NULL DEFAULT '',
    duration_seconds INTEGER NOT NULL DEFAULT 60,
    tts_cache_id TEXT, -- FK to cached_audio.id, generated lazily
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_listening_level ON listening_exercises(jlpt_level);

CREATE TABLE IF NOT EXISTS listening_vocabulary (
    id TEXT PRIMARY KEY,
    exercise_id TEXT NOT NULL,
    word TEXT NOT NULL,
    reading TEXT NOT NULL,
    meaning TEXT NOT NULL,
    FOREIGN KEY (exercise_id) REFERENCES listening_exercises(id)
);

CREATE INDEX IF NOT EXISTS idx_listening_vocab_exercise ON listening_vocabulary(exercise_id);

CREATE TABLE IF NOT EXISTS listening_questions (
    id TEXT PRIMARY KEY,
    exercise_id TEXT NOT NULL,
    question_num INTEGER NOT NULL,
    question TEXT NOT NULL,
    options TEXT NOT NULL DEFAULT '[]', -- JSON array of strings
    correct_index INTEGER NOT NULL DEFAULT 0,
    FOREIGN KEY (exercise_id) REFERENCES listening_exercises(id)
);

CREATE INDEX IF NOT EXISTS idx_listening_questions_exercise ON listening_questions(exercise_id);

CREATE TABLE IF NOT EXISTS listening_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    exercise_id TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    score INTEGER NOT NULL DEFAULT 0,
    total_questions INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'in_progress', -- in_progress, completed
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (exercise_id) REFERENCES listening_exercises(id)
);

CREATE INDEX IF NOT EXISTS idx_listening_sessions_user ON listening_sessions(user_id);

CREATE TABLE IF NOT EXISTS listening_session_answers (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL,
    question_id TEXT NOT NULL,
    answer_index INTEGER NOT NULL,
    is_correct INTEGER NOT NULL DEFAULT 0,
    answered_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (session_id) REFERENCES listening_sessions(id),
    FOREIGN KEY (question_id) REFERENCES listening_questions(id),
    UNIQUE(session_id, question_id)
);
