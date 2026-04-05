-- Add related_words field for synonyms, antonyms, confusable words
ALTER TABLE vocabulary ADD COLUMN related_words JSONB DEFAULT '[]'::jsonb;
ALTER TABLE vocabulary ADD COLUMN word_type VARCHAR(50) DEFAULT 'unknown'; -- verb, noun, adjective, etc.
ALTER TABLE vocabulary ADD COLUMN register VARCHAR(50) DEFAULT 'neutral'; -- formal, casual, neutral, slang
ALTER TABLE vocabulary ADD COLUMN common_mistakes TEXT DEFAULT '';

-- Create index for word lookups
CREATE INDEX idx_vocab_word ON vocabulary(word);
