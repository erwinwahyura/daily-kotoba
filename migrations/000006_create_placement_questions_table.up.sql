CREATE TABLE placement_questions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_text TEXT NOT NULL,
    correct_answer VARCHAR(255) NOT NULL,
    wrong_answers JSONB NOT NULL,
    difficulty_level VARCHAR(10) NOT NULL,
    order_index INTEGER NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_placement_difficulty ON placement_questions(difficulty_level);
CREATE INDEX idx_placement_order ON placement_questions(order_index);
