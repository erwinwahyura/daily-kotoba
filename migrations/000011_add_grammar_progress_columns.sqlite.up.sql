-- Add grammar progress columns to user_progress table
ALTER TABLE user_progress ADD COLUMN current_grammar_index INTEGER NOT NULL DEFAULT 0;
ALTER TABLE user_progress ADD COLUMN last_grammar_id TEXT REFERENCES grammar_patterns(id);
ALTER TABLE user_progress ADD COLUMN grammar_learned_count INTEGER NOT NULL DEFAULT 0;

-- Create user_grammar_status table for tracking grammar pattern status
CREATE TABLE IF NOT EXISTS user_grammar_status (
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    grammar_id TEXT REFERENCES grammar_patterns(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('studied', 'skipped', 'mastered')),
    marked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, grammar_id)
);

-- Index for faster lookups
CREATE INDEX idx_grammar_status_user ON user_grammar_status(user_id);
CREATE INDEX idx_grammar_status_pattern ON user_grammar_status(grammar_id);
