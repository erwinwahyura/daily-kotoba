-- Reading comprehension system

CREATE TABLE IF NOT EXISTS reading_articles (
    id TEXT PRIMARY KEY,
    jlpt_level TEXT NOT NULL,
    title TEXT NOT NULL,
    text TEXT NOT NULL,
    topic TEXT NOT NULL DEFAULT 'general',
    difficulty TEXT NOT NULL DEFAULT 'medium',
    word_count INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reading_articles_level ON reading_articles(jlpt_level);

CREATE TABLE IF NOT EXISTS reading_questions (
    id TEXT PRIMARY KEY,
    article_id TEXT NOT NULL,
    question_num INTEGER NOT NULL,
    question TEXT NOT NULL,
    options TEXT NOT NULL DEFAULT '[]',
    correct_index INTEGER NOT NULL DEFAULT 0,
    explanation TEXT NOT NULL DEFAULT '',
    FOREIGN KEY (article_id) REFERENCES reading_articles(id)
);

CREATE INDEX IF NOT EXISTS idx_reading_questions_article ON reading_questions(article_id);
