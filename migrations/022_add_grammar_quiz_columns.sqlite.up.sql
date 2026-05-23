ALTER TABLE grammar_patterns ADD COLUMN quiz_questions TEXT DEFAULT '[]';
ALTER TABLE grammar_patterns ADD COLUMN somatome_week INTEGER DEFAULT 0;
ALTER TABLE grammar_patterns ADD COLUMN somatome_day INTEGER DEFAULT 0;
ALTER TABLE grammar_patterns ADD COLUMN source TEXT DEFAULT '';
