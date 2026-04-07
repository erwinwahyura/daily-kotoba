CREATE TABLE placement_questions (
    id TEXT PRIMARY KEY,
    question_text TEXT NOT NULL,
    options TEXT NOT NULL, -- JSON array stored as TEXT
    correct_answer TEXT NOT NULL,
    difficulty_level INTEGER NOT NULL,
    target_jlpt_level TEXT,
    hints TEXT, -- JSON stored as TEXT
    explanation TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_difficulty ON placement_questions(difficulty_level);
CREATE INDEX idx_placement_order ON placement_questions(difficulty_level);