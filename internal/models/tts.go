package models

import (
	"time"
)

// TTSRequest represents a text-to-speech generation request
type TTSRequest struct {
	Text     string `json:"text" binding:"required"`
	VoiceID  string `json:"voice_id,omitempty"`
	Language string `json:"language,omitempty"` // "ja" for Japanese
}

// TTSResponse represents the audio URL and metadata
type TTSResponse struct {
	AudioURL    string    `json:"audio_url"`
	ContentType string    `json:"content_type"`
	Duration    float64   `json:"duration_seconds"`
	Cached      bool      `json:"cached"`
	GeneratedAt time.Time `json:"generated_at"`
}

// CachedAudio represents a stored TTS file in the database
type CachedAudio struct {
	ID          string    `json:"id" db:"id"`
	TextHash    string    `json:"text_hash" db:"text_hash"`     // MD5 hash of text
	TextContent string    `json:"text_content" db:"text_content"` // Original text
	VoiceID     string    `json:"voice_id" db:"voice_id"`
	AudioData   []byte    `json:"-" db:"audio_data"`              // Binary MP3 data
	ContentType string    `json:"content_type" db:"content_type"`
	Duration    float64   `json:"duration" db:"duration"`
	UsageCount  int       `json:"usage_count" db:"usage_count"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	LastUsedAt  time.Time `json:"last_used_at" db:"last_used_at"`
}

// TTSConfig holds configuration for TTS services
type TTSConfig struct {
	ElevenLabsAPIKey string
	DefaultVoiceID   string
	CacheEnabled     bool
	MaxCacheSize     int64 // MB
}

// GetDefaultJapaneseVoice returns the best voice for Japanese
func GetDefaultJapaneseVoice() string {
	// ElevenLabs voice optimized for Japanese
	// "Rachel" or "Adam" work well, but ideally a native Japanese voice
	return "XB0fDUnXU5powFXDhLx" // ElevenLabs Japanese-optimized voice
}

// Available voices for Japanese learning
func GetAvailableVoices() []map[string]string {
	return []map[string]string{
		{"id": "XB0fDUnXU5powFXDhLx", "name": "Takumi (Japanese Male)", "lang": "ja"},
		{"id": "XrExE9yKIg1WjnnlVkGX", "name": "Matteo (Clear)", "lang": "multi"},
		{"id": "21m00Tcm4TlvDq8ikWAM", "name": "Rachel (Natural)", "lang": "en"},
	}
}
