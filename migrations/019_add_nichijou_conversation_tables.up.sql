-- Nichijou Conversation Features - Phase 1
-- AI Chat and Shadowing Tables

-- Conversation sessions
CREATE TABLE IF NOT EXISTS conversation_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    mode VARCHAR(20) NOT NULL CHECK (mode IN ('ai', 'peer', 'shadowing', 'micro')),
    scenario_id VARCHAR(50),
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 5),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'completed', 'abandoned')),
    naturalness_avg INTEGER CHECK (naturalness_avg BETWEEN 0 AND 100),
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages within sessions
CREATE TABLE IF NOT EXISTS conversation_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    session_id UUID NOT NULL REFERENCES conversation_sessions(id) ON DELETE CASCADE,
    sender VARCHAR(10) NOT NULL CHECK (sender IN ('user', 'ai', 'peer', 'system')),
    content TEXT NOT NULL,
    correction TEXT,
    naturalness_score INTEGER CHECK (naturalness_score BETWEEN 0 AND 100),
    alternatives JSONB, -- Array of alternative phrasings
    metadata JSONB, -- Additional data like grammar points, vocabulary used
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Shadowing attempts
CREATE TABLE IF NOT EXISTS shadowing_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    prompt_id VARCHAR(50) NOT NULL,
    scenario_id VARCHAR(50),
    native_audio_url TEXT,
    user_audio_url TEXT,
    transcript TEXT,
    accuracy_score INTEGER CHECK (accuracy_score BETWEEN 0 AND 100),
    rhythm_score INTEGER CHECK (rhythm_score BETWEEN 0 AND 100),
    pitch_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Scenarios for AI chat
CREATE TABLE IF NOT EXISTS conversation_scenarios (
    id VARCHAR(50) PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 5),
    category VARCHAR(30) NOT NULL, -- 'daily', 'travel', 'work', 'social'
    icon TEXT,
    system_prompt TEXT NOT NULL, -- LLM system prompt for this scenario
    starter_messages JSONB, -- Array of AI starter messages
    vocabulary_hints JSONB, -- Suggested vocab for this scenario
    grammar_points JSONB, -- Key grammar to practice
    is_active BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Micro-interactions library
CREATE TABLE IF NOT EXISTS micro_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    prompt TEXT NOT NULL,
    level INTEGER NOT NULL CHECK (level BETWEEN 1 AND 3),
    category VARCHAR(30) NOT NULL, -- 'reaction', 'filler', 'small_talk'
    options JSONB NOT NULL, -- Array of response options
    best_answer TEXT,
    explanation TEXT,
    context_notes TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_conversation_sessions_user_id ON conversation_sessions(user_id);
CREATE INDEX idx_conversation_sessions_status ON conversation_sessions(status);
CREATE INDEX idx_conversation_messages_session_id ON conversation_messages(session_id);
CREATE INDEX idx_conversation_messages_created_at ON conversation_messages(created_at);
CREATE INDEX idx_shadowing_attempts_user_id ON shadowing_attempts(user_id);
CREATE INDEX idx_conversation_scenarios_level ON conversation_scenarios(level);
CREATE INDEX idx_conversation_scenarios_category ON conversation_scenarios(category);
CREATE INDEX idx_micro_interactions_level ON micro_interactions(level);

-- Trigger to update updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_conversation_sessions_updated_at
    BEFORE UPDATE ON conversation_sessions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
