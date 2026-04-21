package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSONB type for PostgreSQL JSONB columns
type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &j)
	case string:
		return json.Unmarshal([]byte(v), &j)
	default:
		return json.Unmarshal([]byte(v.(string)), &j)
	}
}

// ConversationSession represents a conversation practice session
type ConversationSession struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	Mode           string    `json:"mode" db:"mode"` // 'ai', 'peer', 'shadowing', 'micro'
	ScenarioID     string    `json:"scenario_id" db:"scenario_id"`
	Level          int       `json:"level" db:"level"`
	Status         string    `json:"status" db:"status"` // 'active', 'completed', 'abandoned'
	NaturalnessAvg int       `json:"naturalness_avg" db:"naturalness_avg"`
	StartedAt      time.Time `json:"started_at" db:"started_at"`
	EndedAt        *time.Time `json:"ended_at" db:"ended_at"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ConversationMessage represents a single message in a conversation
type ConversationMessage struct {
	ID               string    `json:"id" db:"id"`
	SessionID        string    `json:"session_id" db:"session_id"`
	Sender           string    `json:"sender" db:"sender"` // 'user', 'ai', 'peer', 'system'
	Content          string    `json:"content" db:"content"`
	Correction       string    `json:"correction,omitempty" db:"correction"`
	NaturalnessScore int       `json:"naturalness_score,omitempty" db:"naturalness_score"`
	Alternatives     JSONB     `json:"alternatives,omitempty" db:"alternatives"`
	Metadata         JSONB     `json:"metadata,omitempty" db:"metadata"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// ShadowingAttempt represents a voice shadowing practice attempt
type ShadowingAttempt struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	PromptID      string    `json:"prompt_id" db:"prompt_id"`
	ScenarioID    string    `json:"scenario_id,omitempty" db:"scenario_id"`
	NativeAudioURL string   `json:"native_audio_url,omitempty" db:"native_audio_url"`
	UserAudioURL  string    `json:"user_audio_url,omitempty" db:"user_audio_url"`
	Transcript    string    `json:"transcript,omitempty" db:"transcript"`
	AccuracyScore int       `json:"accuracy_score,omitempty" db:"accuracy_score"`
	RhythmScore   int       `json:"rhythm_score,omitempty" db:"rhythm_score"`
	PitchData     JSONB     `json:"pitch_data,omitempty" db:"pitch_data"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ConversationScenario represents a predefined conversation scenario
type ConversationScenario struct {
	ID              string    `json:"id" db:"id"`
	Name            string    `json:"name" db:"name"`
	Description     string    `json:"description" db:"description"`
	Level           int       `json:"level" db:"level"`
	Category        string    `json:"category" db:"category"` // 'daily', 'travel', 'work', 'social'
	Icon            string    `json:"icon,omitempty" db:"icon"`
	SystemPrompt    string    `json:"system_prompt" db:"system_prompt"`
	StarterMessages JSONB     `json:"starter_messages,omitempty" db:"starter_messages"`
	VocabularyHints JSONB     `json:"vocabulary_hints,omitempty" db:"vocabulary_hints"`
	GrammarPoints   JSONB     `json:"grammar_points,omitempty" db:"grammar_points"`
	IsActive        bool      `json:"is_active" db:"is_active"`
	SortOrder       int       `json:"sort_order" db:"sort_order"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// MicroInteraction represents a short casual exchange
type MicroInteraction struct {
	ID            string    `json:"id" db:"id"`
	Prompt        string    `json:"prompt" db:"prompt"`
	Level         int       `json:"level" db:"level"`
	Category      string    `json:"category" db:"category"` // 'reaction', 'filler', 'small_talk'
	Options       JSONB     `json:"options" db:"options"`
	BestAnswer    string    `json:"best_answer,omitempty" db:"best_answer"`
	Explanation   string    `json:"explanation,omitempty" db:"explanation"`
	ContextNotes  string    `json:"context_notes,omitempty" db:"context_notes"`
	IsActive      bool      `json:"is_active" db:"is_active"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

// ChatRequest represents a user message in a conversation
type ChatRequest struct {
	SessionID string `json:"session_id" binding:"required"`
	Message   string `json:"message" binding:"required"`
}

// ChatResponse represents the AI response with corrections
type ChatResponse struct {
	Message          ConversationMessage `json:"message"`
	AIResponse       ConversationMessage `json:"ai_response"`
	NaturalnessScore int                 `json:"naturalness_score"`
	Suggestions      []string            `json:"suggestions,omitempty"`
}

// StartChatRequest represents the request to start a new conversation
type StartChatRequest struct {
	ScenarioID string `json:"scenario_id" binding:"required"`
	Level      int    `json:"level" binding:"required,min=1,max=5"`
}

// StartChatResponse represents the response when starting a conversation
type StartChatResponse struct {
	Session  ConversationSession   `json:"session"`
	Messages []ConversationMessage `json:"messages"`
	Scenario ConversationScenario  `json:"scenario"`
}

// ScenarioListResponse represents the list of available scenarios
type ScenarioListResponse struct {
	Scenarios []ConversationScenario `json:"scenarios"`
}

// MicroInteractionResponse represents a micro-interaction challenge
type MicroInteractionResponse struct {
	Interaction MicroInteraction `json:"interaction"`
}

// MicroInteractionAnswerRequest represents user's answer to micro-interaction
type MicroInteractionAnswerRequest struct {
	InteractionID string `json:"interaction_id" binding:"required"`
	SelectedOption  int    `json:"selected_option" binding:"required,min=0"`
}

// MicroInteractionAnswerResponse represents feedback on micro-interaction answer
type MicroInteractionAnswerResponse struct {
	Correct     bool   `json:"correct"`
	BestAnswer  string `json:"best_answer"`
	Explanation string `json:"explanation"`
}
