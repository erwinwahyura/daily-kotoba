-- Daily goals and achievements system

-- User goal settings (default targets)
CREATE TABLE IF NOT EXISTS goal_settings (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    vocab_target INTEGER DEFAULT 10,
    grammar_target INTEGER DEFAULT 5,
    kanji_target INTEGER DEFAULT 5,
    conjugation_target INTEGER DEFAULT 20,
    reading_target INTEGER DEFAULT 1,
    enable_reminders BOOLEAN DEFAULT false,
    reminder_time TEXT DEFAULT '20:00',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Daily goals tracking
CREATE TABLE IF NOT EXISTS daily_goals (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    vocab_target INTEGER DEFAULT 10,
    vocab_completed INTEGER DEFAULT 0,
    grammar_target INTEGER DEFAULT 5,
    grammar_completed INTEGER DEFAULT 0,
    kanji_target INTEGER DEFAULT 5,
    kanji_completed INTEGER DEFAULT 0,
    conjugation_target INTEGER DEFAULT 20,
    conjugation_completed INTEGER DEFAULT 0,
    reading_target INTEGER DEFAULT 1,
    reading_completed INTEGER DEFAULT 0,
    is_completed BOOLEAN DEFAULT false,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, date)
);

-- User streaks tracking
CREATE TABLE IF NOT EXISTS user_streaks (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_activity_date DATE,
    total_active_days INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id)
);

-- Achievements earned
CREATE TABLE IF NOT EXISTS achievements (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    achievement_id TEXT NOT NULL, -- references AchievementDefinition.ID
    type TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    icon TEXT,
    level INTEGER DEFAULT 1,
    earned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, achievement_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_daily_goals_user_date ON daily_goals(user_id, date);
CREATE INDEX IF NOT EXISTS idx_daily_goals_date ON daily_goals(date);
CREATE INDEX IF NOT EXISTS idx_achievements_user ON achievements(user_id);
CREATE INDEX IF NOT EXISTS idx_achievements_type ON achievements(type);
