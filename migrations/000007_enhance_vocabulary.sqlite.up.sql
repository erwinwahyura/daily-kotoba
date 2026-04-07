-- SQLite has limited ALTER TABLE support
-- Adding columns that don't exist yet
ALTER TABLE vocabulary ADD COLUMN related_words TEXT DEFAULT '[]';
ALTER TABLE vocabulary ADD COLUMN word_type TEXT DEFAULT 'unknown';
ALTER TABLE vocabulary ADD COLUMN register TEXT DEFAULT 'neutral';
ALTER TABLE vocabulary ADD COLUMN common_mistakes TEXT DEFAULT '';

CREATE INDEX idx_vocab_word ON vocabulary(word);