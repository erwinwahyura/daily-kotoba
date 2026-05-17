package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

type TTSRepository struct {
	db *db.DB
}

func NewTTSRepository(database *db.DB) *TTSRepository {
	return &TTSRepository{db: database}
}

// GetCachedAudio retrieves cached audio by text hash
func (r *TTSRepository) GetCachedAudio(textHash string) (*models.CachedAudio, error) {
	cached := &models.CachedAudio{}
	query := `
		SELECT id, text_hash, text_content, voice_id, audio_data, 
		       content_type, duration, usage_count, created_at, last_used_at
		FROM cached_audio
		WHERE text_hash = $1
	`
	err := r.db.QueryRow(query, textHash).Scan(
		&cached.ID, &cached.TextHash, &cached.TextContent, &cached.VoiceID,
		&cached.AudioData, &cached.ContentType, &cached.Duration,
		&cached.UsageCount, &cached.CreatedAt, &cached.LastUsedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return cached, err
}

// SaveCachedAudio stores generated audio in cache
func (r *TTSRepository) SaveCachedAudio(cached *models.CachedAudio) error {
	query := `
		INSERT INTO cached_audio 
		(id, text_hash, text_content, voice_id, audio_data, content_type, duration, usage_count, created_at, last_used_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (text_hash) DO UPDATE SET
			last_used_at = EXCLUDED.last_used_at,
			usage_count = cached_audio.usage_count + 1
	`
	_, err := r.db.Exec(query,
		cached.ID, cached.TextHash, cached.TextContent, cached.VoiceID,
		cached.AudioData, cached.ContentType, cached.Duration,
		cached.UsageCount, cached.CreatedAt, cached.LastUsedAt,
	)
	return err
}

// IncrementUsageCount increases the usage counter
func (r *TTSRepository) IncrementUsageCount(textHash string) error {
	query := `UPDATE cached_audio SET usage_count = usage_count + 1, last_used_at = $1 WHERE text_hash = $2`
	_, err := r.db.Exec(query, time.Now(), textHash)
	return err
}

// GetAudioByID retrieves audio data by its unique ID
func (r *TTSRepository) GetAudioByID(audioID string) ([]byte, string, error) {
	var audioData []byte
	var contentType string
	query := `SELECT audio_data, content_type FROM cached_audio WHERE id = $1`
	err := r.db.QueryRow(query, audioID).Scan(&audioData, &contentType)
	if err == sql.ErrNoRows {
		return nil, "", fmt.Errorf("audio not found")
	}
	return audioData, contentType, err
}

// CleanupOldCache removes least recently used audio files beyond max cache size
func (r *TTSRepository) CleanupOldCache(maxSizeMB int64) error {
	// Calculate total size (SQLite doesn't have easy blob size query)
	// For now, just delete old entries based on last_used_at
	cutoff := time.Now().AddDate(0, -1, 0) // 1 month old
	query := `DELETE FROM cached_audio WHERE last_used_at < $1 AND usage_count < 5`
	_, err := r.db.Exec(query, cutoff)
	return err
}

// GetCacheStats returns statistics about the audio cache
func (r *TTSRepository) GetCacheStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Count
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM cached_audio").Scan(&count)
	if err != nil {
		return nil, err
	}
	stats["cached_count"] = count
	
	// Total duration
	var totalDuration float64
	err = r.db.QueryRow("SELECT COALESCE(SUM(duration), 0) FROM cached_audio").Scan(&totalDuration)
	if err != nil {
		return nil, err
	}
	stats["total_duration_seconds"] = totalDuration
	
	// Most used
	var mostUsed string
	var mostUsedCount int
	row := r.db.QueryRow("SELECT text_content, usage_count FROM cached_audio ORDER BY usage_count DESC LIMIT 1")
	if err := row.Scan(&mostUsed, &mostUsedCount); err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	stats["most_used_phrase"] = mostUsed
	stats["most_used_count"] = mostUsedCount
	
	return stats, nil
}