-- Kanji writing practice system

CREATE TABLE IF NOT EXISTS kanji_characters (
    id TEXT PRIMARY KEY,
    character TEXT NOT NULL UNIQUE,
    jlpt_level TEXT NOT NULL,
    meaning TEXT NOT NULL,
    readings TEXT NOT NULL DEFAULT '[]', -- JSON array of reading strings
    stroke_count INTEGER NOT NULL DEFAULT 1,
    stroke_order TEXT NOT NULL DEFAULT '[]', -- JSON array of Stroke objects
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_kanji_jlpt ON kanji_characters(jlpt_level);
CREATE INDEX IF NOT EXISTS idx_kanji_character ON kanji_characters(character);

CREATE TABLE IF NOT EXISTS kanji_practice_sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    kanji_id TEXT NOT NULL,
    kanji_char TEXT NOT NULL,
    started_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    completed_at DATETIME,
    current_stroke INTEGER NOT NULL DEFAULT 0,
    accuracy REAL NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'in_progress', -- in_progress, completed
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (kanji_id) REFERENCES kanji_characters(id)
);

CREATE INDEX IF NOT EXISTS idx_kanji_sessions_user ON kanji_practice_sessions(user_id);
