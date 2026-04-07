-- Add grammar progress tracking to user_progress table
ALTER TABLE user_progress ADD COLUMN current_grammar_index INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user_progress ADD COLUMN grammar_learned_count INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user_progress ADD COLUMN last_grammar_id TEXT REFERENCES grammar_patterns(id);

-- Create grammar status tracking table
CREATE TABLE user_grammar_status (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    grammar_id TEXT NOT NULL REFERENCES grammar_patterns(id),
    status TEXT NOT NULL CHECK (status IN ('learning', 'mastered', 'skipped')),
    marked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, grammar_id)
);

CREATE INDEX idx_user_grammar_status_user ON user_grammar_status(user_id);
CREATE INDEX idx_user_grammar_status_grammar ON user_grammar_status(grammar_id);