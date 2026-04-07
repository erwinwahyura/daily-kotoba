CREATE TABLE user_vocab_status (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocab_id TEXT NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    marked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, vocab_id)
);

CREATE INDEX idx_user_vocab_status ON user_vocab_status(user_id, status);
CREATE INDEX idx_user_vocab_user_id ON user_vocab_status(user_id);