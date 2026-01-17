package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Vocabulary struct {
	ID                  string          `json:"id" db:"id"`
	Word                string          `json:"word" db:"word"`
	Reading             string          `json:"reading" db:"reading"`
	ShortMeaning        string          `json:"short_meaning" db:"short_meaning"`
	DetailedExplanation string          `json:"detailed_explanation" db:"detailed_explanation"`
	ExampleSentences    ExampleSentences `json:"example_sentences" db:"example_sentences"`
	UsageNotes          string          `json:"usage_notes" db:"usage_notes"`
	JLPTLevel           string          `json:"jlpt_level" db:"jlpt_level"`
	IndexPosition       int             `json:"index_position" db:"index_position"`
	CreatedAt           time.Time       `json:"created_at" db:"created_at"`
}

// ExampleSentences is a custom type for JSONB support
type ExampleSentences []string

// Scan implements sql.Scanner interface
func (es *ExampleSentences) Scan(value interface{}) error {
	if value == nil {
		*es = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value")
	}

	return json.Unmarshal(bytes, es)
}

// Value implements driver.Valuer interface
func (es ExampleSentences) Value() (driver.Value, error) {
	if es == nil {
		return nil, nil
	}
	return json.Marshal(es)
}

type VocabularyWithProgress struct {
	Vocabulary *Vocabulary         `json:"vocabulary"`
	Progress   *VocabularyProgress `json:"progress"`
}

type VocabularyProgress struct {
	CurrentIndex       int `json:"current_index"`
	TotalWordsInLevel  int `json:"total_words_in_level"`
	WordsLearned       int `json:"words_learned"`
	StreakDays         int `json:"streak_days"`
}

type SkipRequest struct {
	Status string `json:"status" binding:"required,oneof=known skipped"`
}

type VocabularyListResponse struct {
	Vocabulary []Vocabulary       `json:"vocabulary"`
	Pagination PaginationResponse `json:"pagination"`
}

type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
