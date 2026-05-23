package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// GrammarPattern represents an N3-N1 grammar form with detailed teaching content
type GrammarPattern struct {
	ID                   string           `json:"id" db:"id"`
	Pattern              string           `json:"pattern" db:"pattern"`
	PlainForm            string           `json:"plain_form" db:"plain_form"`
	Meaning              string           `json:"meaning" db:"meaning"`
	DetailedExplanation  string           `json:"detailed_explanation" db:"detailed_explanation"`
	ConjugationRules     string           `json:"conjugation_rules" db:"conjugation_rules"`
	UsageExamples        UsageExamples    `json:"usage_examples" db:"usage_examples"`
	NuanceNotes          string           `json:"nuance_notes" db:"nuance_notes"`
	JLPTLevel            string           `json:"jlpt_level" db:"jlpt_level"`
	RelatedPatterns      RelatedPatterns  `json:"related_patterns" db:"related_patterns"`
	CommonMistakes       string           `json:"common_mistakes" db:"common_mistakes"`
	IndexPosition        int              `json:"index_position" db:"index_position"`
	QuizQuestions        QuizQuestions    `json:"quiz_questions" db:"quiz_questions"`
	SomatomeWeek         int              `json:"somatome_week" db:"somatome_week"`
	SomatomeDay          int              `json:"somatome_day" db:"somatome_day"`
	Source               string           `json:"source" db:"source"`
	CreatedAt            time.Time        `json:"created_at" db:"created_at"`
}

// QuizQuestion is a multiple-choice question for a grammar pattern
type QuizQuestion struct {
	ID           string   `json:"id"`
	Question     string   `json:"question"`     // Japanese sentence with blank ___
	English      string   `json:"english"`      // English hint
	Options      []string `json:"options"`      // 4 choices
	CorrectIndex int      `json:"correct_index"`
	Explanation  string   `json:"explanation"`
}

// QuizQuestions is a custom JSON type
type QuizQuestions []QuizQuestion

func (qq *QuizQuestions) Scan(value interface{}) error {
	if value == nil {
		*qq = []QuizQuestion{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal QuizQuestions: got %T", value)
	}
	return json.Unmarshal(bytes, qq)
}

func (qq QuizQuestions) Value() (driver.Value, error) {
	if qq == nil {
		return "[]", nil
	}
	return json.Marshal(qq)
}

// UsageExample pairs a sentence with detailed explanation
type UsageExample struct {
	Japanese     string `json:"japanese"`
	Reading      string `json:"reading"`       // Full reading with kanji
	Meaning      string `json:"meaning"`       // English meaning
	Nuance       string `json:"nuance"`        // Why this pattern fits here
	Context      string `json:"context"`       // Situation where used
	Alternative  string `json:"alternative"` // What you might say instead
}

// UsageExamples is a custom type for JSONB
type UsageExamples []UsageExample

// Scan implements sql.Scanner interface
func (ue *UsageExamples) Scan(value interface{}) error {
	if value == nil {
		*ue = []UsageExample{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal UsageExamples JSONB value: expected []byte or string, got %T", value)
	}
	return json.Unmarshal(bytes, ue)
}

// Value implements driver.Valuer interface
func (ue UsageExamples) Value() (driver.Value, error) {
	if ue == nil {
		return nil, nil
	}
	return json.Marshal(ue)
}

// RelatedPattern links to similar/confusable grammar
type RelatedPattern struct {
	Pattern      string `json:"pattern"`
	Relationship string `json:"relationship"` // "similar meaning but different nuance", "often confused with", etc.
	KeyDifference string `json:"key_difference"`
}

// RelatedPatterns is a custom type for JSONB
type RelatedPatterns []RelatedPattern

// Scan implements sql.Scanner interface
func (rp *RelatedPatterns) Scan(value interface{}) error {
	if value == nil {
		*rp = []RelatedPattern{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to unmarshal RelatedPatterns JSONB value: expected []byte or string, got %T", value)
	}
	return json.Unmarshal(bytes, rp)
}

// Value implements driver.Valuer interface
func (rp RelatedPatterns) Value() (driver.Value, error) {
	return json.Marshal(rp)
}

// GrammarProgress tracks user progress through grammar patterns
type GrammarProgress struct {
	CurrentIndex      int `json:"current_index"`
	TotalPatterns     int `json:"total_patterns"`
	PatternsLearned   int `json:"patterns_learned"`
	MasteryByLevel    map[string]int `json:"mastery_by_level"` // N3, N2, N1 counts
}

type GrammarPatternResponse struct {
	Pattern  *GrammarPattern `json:"pattern"`
	Progress *GrammarProgress `json:"progress"`
}

type GrammarListResponse struct {
	Patterns   []GrammarPattern `json:"patterns"`
	Pagination PaginationResponse `json:"pagination"`
}
