-- Grammar patterns table for N3-N1 grammar forms
CREATE TABLE grammar_patterns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    pattern VARCHAR(100) NOT NULL, -- e.g., "〜わけにはいかない"
    plain_form VARCHAR(100) NOT NULL, -- e.g., "わけにはいかない"
    meaning TEXT NOT NULL,
    detailed_explanation TEXT NOT NULL,
    conjugation_rules TEXT NOT NULL, -- How to attach to verbs/adjectives
    usage_examples JSONB NOT NULL, -- Example sentences with explanations
    nuance_notes TEXT NOT NULL, -- When to use vs alternatives
    jlpt_level VARCHAR(10) NOT NULL,
    related_patterns JSONB DEFAULT '[]'::jsonb, -- Similar patterns that confuse learners
    common_mistakes TEXT DEFAULT '',
    index_position INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_grammar_level ON grammar_patterns(jlpt_level);
CREATE INDEX idx_grammar_level_index ON grammar_patterns(jlpt_level, index_position);
CREATE UNIQUE INDEX idx_grammar_pattern_unique ON grammar_patterns(jlpt_level, index_position);
