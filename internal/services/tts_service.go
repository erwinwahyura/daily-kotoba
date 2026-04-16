package services

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
)

type TTSService struct {
	ttsRepo    *repository.TTSRepository
	apiKey     string
	defaultVoice string
	client     *http.Client
}

func NewTTSService(ttsRepo *repository.TTSRepository) *TTSService {
	apiKey := os.Getenv("ELEVENLABS_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("ELEVENLABS_API_KEY")
	}
	
	return &TTSService{
		ttsRepo:      ttsRepo,
		apiKey:       apiKey,
		defaultVoice: models.GetDefaultJapaneseVoice(),
		client:       &http.Client{Timeout: 30 * time.Second},
	}
}

// GenerateAudio generates or retrieves cached TTS audio
func (s *TTSService) GenerateAudio(text, voiceID string) (*models.TTSResponse, error) {
	if text == "" {
		return nil, fmt.Errorf("text is required")
	}
	
	if voiceID == "" {
		voiceID = s.defaultVoice
	}
	
	// Create hash for cache lookup
	textHash := hashText(text + voiceID)
	
	// Check cache first
	cached, err := s.ttsRepo.GetCachedAudio(textHash)
	if err != nil {
		return nil, err
	}
	
	if cached != nil {
		// Update usage stats
		s.ttsRepo.IncrementUsageCount(textHash)
		
		return &models.TTSResponse{
			AudioURL:    fmt.Sprintf("/api/tts/audio/%s", cached.ID),
			ContentType: cached.ContentType,
			Duration:    cached.Duration,
			Cached:      true,
			GeneratedAt: cached.CreatedAt,
		}, nil
	}
	
	// No cache hit - generate via ElevenLabs
	if s.apiKey == "" {
		return nil, fmt.Errorf("TTS service not configured - no API key")
	}
	
	audioData, duration, err := s.generateElevenLabs(text, voiceID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate audio: %w", err)
	}
	
	// Store in cache
	cached = &models.CachedAudio{
		ID:          uuid.New().String(),
		TextHash:    textHash,
		TextContent: text,
		VoiceID:     voiceID,
		AudioData:   audioData,
		ContentType: "audio/mpeg",
		Duration:    duration,
		UsageCount:  1,
		CreatedAt:   time.Now(),
		LastUsedAt:  time.Now(),
	}
	
	if err := s.ttsRepo.SaveCachedAudio(cached); err != nil {
		// Log but don't fail - we can still return the audio
		fmt.Printf("Failed to cache audio: %v\n", err)
	}
	
	return &models.TTSResponse{
		AudioURL:    fmt.Sprintf("/api/tts/audio/%s", cached.ID),
		ContentType: cached.ContentType,
		Duration:    duration,
		Cached:      false,
		GeneratedAt: time.Now(),
	}, nil
}

// GetAudioByID retrieves raw audio data by ID
func (s *TTSService) GetAudioByID(audioID string) ([]byte, string, error) {
	return s.ttsRepo.GetAudioByID(audioID)
}

// generateElevenLabs calls the ElevenLabs API
func (s *TTSService) generateElevenLabs(text, voiceID string) ([]byte, float64, error) {
	url := fmt.Sprintf("https://api.elevenlabs.io/v1/text-to-speech/%s", voiceID)
	
	payload := map[string]interface{}{
		"text":     text,
		"model_id": "eleven_multilingual_v2",
		"voice_settings": map[string]interface{}{
			"stability":        0.5,
			"similarity_boost": 0.75,
		},
	}
	
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, 0, err
	}
	
	req.Header.Set("Accept", "audio/mpeg")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("xi-api-key", s.apiKey)
	
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, 0, fmt.Errorf("ElevenLabs API error %d: %s", resp.StatusCode, string(body))
	}
	
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	
	// Estimate duration (rough calculation: ~3 chars per second for Japanese)
	duration := float64(len(text)) / 3.0
	
	return audioData, duration, nil
}

// hashText creates MD5 hash of text
func hashText(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetAvailableVoices returns list of available TTS voices
func (s *TTSService) GetAvailableVoices() []map[string]string {
	return models.GetAvailableVoices()
}

// GetCacheStats returns cache statistics
func (s *TTSService) GetCacheStats() (map[string]interface{}, error) {
	return s.ttsRepo.GetCacheStats()
}

// CleanupCache removes old cache entries
func (s *TTSService) CleanupCache() error {
	return s.ttsRepo.CleanupOldCache(100) // 100MB max
}