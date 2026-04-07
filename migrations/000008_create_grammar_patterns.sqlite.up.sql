-- Grammar patterns table for N3-N1 grammar forms
CREATE TABLE grammar_patterns (
    id TEXT PRIMARY KEY,
    pattern TEXT NOT NULL,
    plain_form TEXT NOT NULL,
    meaning TEXT NOT NULL,
    detailed_explanation TEXT NOT NULL,
    conjugation_rules TEXT NOT NULL,
    usage_examples TEXT NOT NULL, -- JSON stored as TEXT
    nuance_notes TEXT NOT NULL,
    jlpt_level TEXT NOT NULL,
    related_patterns TEXT DEFAULT '[]', -- JSON stored as TEXT
    common_mistakes TEXT DEFAULT '',
    index_position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_grammar_level ON grammar_patterns(jlpt_level);
CREATE INDEX idx_grammar_level_index ON grammar_patterns(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_grammar_pattern_unique ON grammar_patterns(jlpt_level, index_position);