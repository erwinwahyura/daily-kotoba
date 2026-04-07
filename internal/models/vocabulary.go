package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
	// Enhanced fields
	RelatedWords        RelatedWords    `json:"related_words" db:"related_words"`
	WordType            string          `json:"word_type" db:"word_type"`
	Register            string          `json:"register" db:"register"`
	CommonMistakes      string          `json:"common_mistakes" db:"common_mistakes"`
}

// RelatedWords stores synonyms, antonyms, and confusable words
type RelatedWords struct {
	Synonyms     []string `json:"synonyms"`
	Antonyms     []string `json:"antonyms"`
	Confusable   []string `json:"confusable"` // Words that look/sound similar
	SeeAlso      []string `json:"see_also"`   // Related concepts
}

// ExampleSentences is a custom type for JSONB support
type ExampleSentences []string

// Scan implements sql.Scanner interface
func (es *ExampleSentences) Scan(value interface{}) error {
	if value == nil {
		*es = []string{}
		return nil
	}

	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal JSONB value: expected []byte or string, got %T", value)
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

// Scan implements sql.Scanner interface for RelatedWords
func (rw *RelatedWords) Scan(value interface{}) error {
	if value == nil {
		*rw = RelatedWords{}
		return nil
	}
	
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal RelatedWords JSONB value: expected []byte or string, got %T", value)
	}
	
	return json.Unmarshal(bytes, rw)
}

// Value implements driver.Valuer interface for RelatedWords
func (rw RelatedWords) Value() (driver.Value, error) {
	return json.Marshal(rw)
}
