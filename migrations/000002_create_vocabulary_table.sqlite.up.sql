CREATE TABLE vocabulary (
    id TEXT PRIMARY KEY,
    word TEXT NOT NULL,
    reading TEXT NOT NULL,
    short_meaning TEXT NOT NULL,
    detailed_explanation TEXT NOT NULL,
    example_sentences TEXT, -- JSON as TEXT
    usage_notes TEXT,
    jlpt_level TEXT NOT NULL,
    index_position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vocab_level ON vocabulary(jlpt_level);
CREATE INDEX idx_vocab_level_index ON vocabulary(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_vocab_level_position_unique ON vocabulary(jlpt_level, index_position);
