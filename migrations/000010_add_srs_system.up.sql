-- SRS (Spaced Repetition) System Tables
-- Implements SM-2 algorithm for optimal retention

CREATE TABLE srs_schedules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id UUID NOT NULL,  -- Can reference vocabulary.id or grammar_patterns.id
    item_type VARCHAR(20) NOT NULL CHECK (item_type IN ('vocabulary', 'grammar')),
    
    -- SM-2 Algorithm fields
    interval_days INTEGER DEFAULT 0,      -- Current interval in days
    repetitions INTEGER DEFAULT 0,      -- Number of successful reviews
    ease_factor REAL DEFAULT 2.5,      -- Easiness factor (starts at 2.5)
    
    -- Review tracking
    last_reviewed_at TIMESTAMP,
    next_review_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Performance stats
    total_reviews INTEGER DEFAULT 0,
    correct_reviews INTEGER DEFAULT 0,
    streak INTEGER DEFAULT 0,           -- Current correct streak
    
    -- Status
    status VARCHAR(20) DEFAULT 'learning' CHECK (status IN ('learning', 'review', 'mastered', 'lapsed')),
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    -- Unique constraint: one schedule per user per item
    UNIQUE(user_id, item_id, item_type)
);

-- Indexes for efficient queries
CREATE INDEX idx_srs_user_next_review ON srs_schedules(user_id, next_review_at) 
    WHERE status IN ('learning', 'review');
CREATE INDEX idx_srs_user_status ON srs_schedules(user_id, status);
CREATE INDEX idx_srs_overdue ON srs_schedules(next_review_at) 
    WHERE next_review_at <= CURRENT_TIMESTAMP AND status IN ('learning', 'review');

-- Review history table (audit log)
CREATE TABLE srs_review_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    schedule_id UUID NOT NULL REFERENCES srs_schedules(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Review details
    quality INTEGER NOT NULL CHECK (quality BETWEEN 0 AND 5),  -- SM-2 quality: 0=complete blackout, 5=perfect
    response_time_ms INTEGER,  -- How long they took to answer
    
    -- What was shown
    item_type VARCHAR(20) NOT NULL,
    item_id UUID NOT NULL,
    
    -- Scheduling before/after
    interval_before INTEGER,
    interval_after INTEGER,
    ease_factor_before REAL,
    ease_factor_after REAL,
    
    reviewed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_srs_history_user ON srs_review_history(user_id, reviewed_at DESC);
CREATE INDEX idx_srs_history_schedule ON srs_review_history(schedule_id, reviewed_at DESC);

-- Daily study stats (for streaks and analytics)
CREATE TABLE study_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_date DATE NOT NULL DEFAULT CURRENT_DATE,
    
    -- Activity counts
    new_items INTEGER DEFAULT 0,           -- New words/patterns learned
    review_items INTEGER DEFAULT 0,        -- SRS reviews completed
    drill_items INTEGER DEFAULT 0,         -- Conjugation drills done
    
    -- Time spent
    total_time_seconds INTEGER DEFAULT 0,
    
    -- Accuracy
    correct_count INTEGER DEFAULT 0,
    total_attempts INTEGER DEFAULT 0,
    
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, session_date)
);

CREATE INDEX idx_study_sessions_user_date ON study_sessions(user_id, session_date DESC);

-- Function to update SRS schedule based on SM-2 algorithm
-- Quality: 0-5 (0=complete failure, 5=perfect)
CREATE OR REPLACE FUNCTION calculate_srs_review(
    p_quality INTEGER,
    p_current_interval INTEGER,
    p_repetitions INTEGER,
    p_ease_factor REAL
) RETURNS TABLE (
    new_interval INTEGER,
    new_repetitions INTEGER,
    new_ease_factor REAL,
    new_status VARCHAR(20)
) AS $$
BEGIN
    -- If quality < 3, reset repetitions but don't reset interval completely
    IF p_quality < 3 THEN
        new_repetitions := 0;
        new_interval := 1;  -- Review again tomorrow
        new_ease_factor := GREATEST(1.3, p_ease_factor - 0.2);
        new_status := 'learning';
    ELSE
        -- Successful review
        new_repetitions := p_repetitions + 1;
        
        -- Calculate new interval
        IF new_repetitions = 1 THEN
            new_interval := 1;  -- 1 day
        ELSIF new_repetitions = 2 THEN
            new_interval := 6;  -- 6 days
        ELSE
            -- Interval * ease factor, rounded
            new_interval := ROUND(p_current_interval * p_ease_factor)::INTEGER;
        END IF;
        
        -- Update ease factor based on quality
        new_ease_factor := p_ease_factor + (0.1 - (5 - p_quality) * (0.08 + (5 - p_quality) * 0.02));
        new_ease_factor := GREATEST(1.3, new_ease_factor);  -- Minimum 1.3
        
        -- Status based on repetitions
        IF new_repetitions >= 8 THEN
            new_status := 'mastered';
        ELSE
            new_status := 'review';
        END IF;
    END IF;
    
    RETURN NEXT;
END;
$$ LANGUAGE plpgsql;

-- SQLite version (for our current setup)
-- Note: SQLite doesn't support CREATE OR REPLACE FUNCTION
-- We'll implement the logic in Go code instead
