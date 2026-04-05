-- Revert vocabulary enhancements
ALTER TABLE vocabulary DROP COLUMN IF EXISTS related_words;
ALTER TABLE vocabulary DROP COLUMN IF EXISTS word_type;
ALTER TABLE vocabulary DROP COLUMN IF EXISTS register;
ALTER TABLE vocabulary DROP COLUMN IF EXISTS common_mistakes;
DROP INDEX IF EXISTS idx_vocab_word;
