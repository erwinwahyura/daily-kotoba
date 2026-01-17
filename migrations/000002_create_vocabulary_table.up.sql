CREATE TABLE vocabulary (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    word VARCHAR(100) NOT NULL,
    reading VARCHAR(100) NOT NULL,
    short_meaning VARCHAR(255) NOT NULL,
    detailed_explanation TEXT NOT NULL,
    example_sentences JSONB,
    usage_notes TEXT,
    jlpt_level VARCHAR(10) NOT NULL,
    index_position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vocab_level ON vocabulary(jlpt_level);
CREATE INDEX idx_vocab_level_index ON vocabulary(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_vocab_level_position_unique ON vocabulary(jlpt_level, index_position);
