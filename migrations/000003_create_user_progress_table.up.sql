CREATE TABLE user_progress (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    current_vocab_index INTEGER NOT NULL DEFAULT 0,
    last_word_id UUID REFERENCES vocabulary(id),
    streak_days INTEGER NOT NULL DEFAULT 0,
    last_study_date DATE,
    words_learned_count INTEGER NOT NULL DEFAULT 0,
    words_skipped_count INTEGER NOT NULL DEFAULT 0,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER update_user_progress_updated_at BEFORE UPDATE ON user_progress
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
