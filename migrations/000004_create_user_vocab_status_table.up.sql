CREATE TABLE user_vocab_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocab_id UUID NOT NULL REFERENCES vocabulary(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL,
    marked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, vocab_id)
);

CREATE INDEX idx_user_vocab_status ON user_vocab_status(user_id, status);
CREATE INDEX idx_user_vocab_user_id ON user_vocab_status(user_id);
