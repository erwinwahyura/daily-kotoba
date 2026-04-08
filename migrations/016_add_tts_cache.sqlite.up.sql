-- TTS Audio Cache Table
CREATE TABLE IF NOT EXISTS cached_audio (
    id TEXT PRIMARY KEY,
    text_hash TEXT NOT NULL UNIQUE,
    text_content TEXT NOT NULL,
    voice_id TEXT NOT NULL DEFAULT 'XB0fDUnXU5powFXDhLx',
    audio_data BLOB NOT NULL,
    content_type TEXT NOT NULL DEFAULT 'audio/mpeg',
    duration REAL DEFAULT 0,
    usage_count INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_cached_audio_hash ON cached_audio(text_hash);
CREATE INDEX IF NOT EXISTS idx_cached_audio_used ON cached_audio(last_used_at);