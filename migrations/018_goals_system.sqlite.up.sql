-- Goals system: user settings, daily activity tracking, achievements

CREATE TABLE IF NOT EXISTS user_goal_settings (
    user_id TEXT PRIMARY KEY,
    vocab_target INTEGER NOT NULL DEFAULT 10,
    grammar_target INTEGER NOT NULL DEFAULT 5,
    kanji_target INTEGER NOT NULL DEFAULT 5,
    conjugation_target INTEGER NOT NULL DEFAULT 20,
    reading_target INTEGER NOT NULL DEFAULT 1,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS user_daily_activities (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    activity_date DATE NOT NULL,
    activity_type TEXT NOT NULL, -- vocab, grammar, kanji, conjugation, reading, review
    count INTEGER NOT NULL DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    UNIQUE(user_id, activity_date, activity_type)
);

CREATE INDEX IF NOT EXISTS idx_daily_activities_user_date ON user_daily_activities(user_id, activity_date);

CREATE TABLE IF NOT EXISTS achievements (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT '🏆',
    category TEXT NOT NULL, -- streak, vocab, grammar, conjugation, reading, kanji
    requirement_type TEXT NOT NULL, -- streak_days, total_count, single_session
    requirement_value INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_achievements (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    achievement_id TEXT NOT NULL,
    earned_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (achievement_id) REFERENCES achievements(id),
    UNIQUE(user_id, achievement_id)
);

CREATE INDEX IF NOT EXISTS idx_user_achievements_user ON user_achievements(user_id);
