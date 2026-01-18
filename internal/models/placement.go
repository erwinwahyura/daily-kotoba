package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// PlacementQuestion represents a placement test question
type PlacementQuestion struct {
	ID              string        `json:"id"`
	QuestionText    string        `json:"question"`
	CorrectAnswer   string        `json:"-"` // Hidden from API response
	WrongAnswers    WrongAnswers  `json:"-"` // Hidden from API response
	Options         []string      `json:"options"` // Shuffled options for client
	DifficultyLevel string        `json:"difficulty"`
	OrderIndex      int           `json:"order_index"`
	CreatedAt       time.Time     `json:"created_at"`
}

// PlacementQuestionResponse is what gets sent to the client (without correct answer)
type PlacementQuestionResponse struct {
	ID              string   `json:"id"`
	QuestionText    string   `json:"question"`
	Options         []string `json:"options"`
	DifficultyLevel string   `json:"difficulty"`
	OrderIndex      int      `json:"order_index"`
}

// WrongAnswers is a custom type for JSONB array
type WrongAnswers []string

// Scan implements sql.Scanner for reading from database
func (w *WrongAnswers) Scan(value interface{}) error {
	if value == nil {
		*w = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan WrongAnswers")
	}

	return json.Unmarshal(bytes, w)
}

// Value implements driver.Valuer for writing to database
func (w WrongAnswers) Value() (driver.Value, error) {
	if len(w) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(w)
}

// PlacementTestResult represents a user's placement test result
type PlacementTestResult struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	TestScore     int       `json:"test_score"`
	AssignedLevel string    `json:"assigned_level"`
	CompletedAt   time.Time `json:"completed_at"`
}

// PlacementTestSubmission represents the request body for submitting test
type PlacementTestSubmission struct {
	Answers map[string]string `json:"answers"` // map[question_id]answer
}

// PlacementTestResponse represents the response after submitting test
type PlacementTestResponse struct {
	Score         int            `json:"score"`
	TotalQuestions int           `json:"total_questions"`
	AssignedLevel string         `json:"assigned_level"`
	Breakdown     map[string]int `json:"breakdown"` // correct answers per level
}
