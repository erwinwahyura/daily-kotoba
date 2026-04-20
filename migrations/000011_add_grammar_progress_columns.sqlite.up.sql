-- Add grammar progress columns to user_progress table (SQLite version)
-- These columns might already exist if coming from a different migration path
-- Using CREATE TABLE IF NOT EXISTS for the grammar status table

-- Create user_grammar_status table for tracking grammar pattern status
CREATE TABLE IF NOT EXISTS user_grammar_status (
    user_id TEXT REFERENCES users(id) ON DELETE CASCADE,
    grammar_id TEXT REFERENCES grammar_patterns(id) ON DELETE CASCADE,
    status TEXT NOT NULL CHECK (status IN ('studied', 'skipped', 'mastered')),
    marked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, grammar_id)
);

-- Index for faster lookups
CREATE INDEX IF NOT EXISTS idx_grammar_status_user ON user_grammar_status(user_id);
CREATE INDEX IF NOT EXISTS idx_grammar_status_pattern ON user_grammar_status(grammar_id);

-- Note: current_grammar_index, last_grammar_id, grammar_learned_count columns
-- should be added to user_progress only if they don't exist
-- SQLite doesn't support IF NOT EXISTS for ALTER TABLE ADD COLUMN
-- The application should handle missing columns gracefully