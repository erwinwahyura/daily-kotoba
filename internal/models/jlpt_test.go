package models

import (
	"time"
)

// JLPTTest represents a full mock exam
type JLPTTest struct {
	ID          string        `json:"id" db:"id"`
	Level       string        `json:"level" db:"level"`         // N5, N4, N3, N2, N1
	Section     string        `json:"section" db:"section"`     // vocab, grammar, reading, listening
	Title       string        `json:"title" db:"title"`
	Description string        `json:"description" db:"description"`
	TimeLimit   int           `json:"time_limit_minutes" db:"time_limit_minutes"`
	TotalQuestions int        `json:"total_questions" db:"total_questions"`
	PassingScore   int        `json:"passing_score" db:"passing_score"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
}

// JLPTQuestion represents a single test question
type JLPTQuestion struct {
	ID           string   `json:"id" db:"id"`
	TestID       string   `json:"test_id" db:"test_id"`
	QuestionNum  int      `json:"question_num" db:"question_num"`
	Type         string   `json:"type" db:"type"`           // multiple_choice, fill_blank, reading_comp
	Question     string   `json:"question" db:"question"`   // Japanese text
	QuestionReading string `json:"question_reading,omitempty" db:"question_reading"`
	EnglishPrompt   string `json:"english_prompt,omitempty" db:"english_prompt"` // Translation hint
	Options      []string `json:"options" db:"options"`       // JSON array
	CorrectIndex int      `json:"correct_index" db:"correct_index"`
	Explanation  string   `json:"explanation" db:"explanation"` // Why this is correct
	PointValue   int      `json:"point_value" db:"point_value"` // Usually 1
	SkillTested  string   `json:"skill_tested" db:"skill_tested"` // vocab, grammar, kanji, etc.
}

// UserTestSession tracks a user's test attempt
type UserTestSession struct {
	ID              string                 `json:"id" db:"id"`
	UserID          string                 `json:"user_id" db:"user_id"`
	TestID          string                 `json:"test_id" db:"test_id"`
	Level           string                 `json:"level" db:"level"`
	StartedAt       time.Time              `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time             `json:"completed_at,omitempty" db:"completed_at"`
	TimeSpentSec    int                    `json:"time_spent_sec" db:"time_spent_sec"`
	Answers         map[string]int         `json:"answers" db:"answers"` // question_id -> selected_option_index
	Score           int                    `json:"score" db:"score"`
	CorrectCount    int                    `json:"correct_count" db:"correct_count"`
	Status          string                 `json:"status" db:"status"` // in_progress, completed, abandoned
}

// TestResult represents the final results
type TestResult struct {
	SessionID       string                 `json:"session_id"`
	Level           string                 `json:"level"`
	Score           int                    `json:"score"`
	TotalQuestions  int                    `json:"total_questions"`
	CorrectCount    int                    `json:"correct_count"`
	IncorrectCount  int                    `json:"incorrect_count"`
	Percentage      float64                `json:"percentage"`
	Passed          bool                   `json:"passed"`
	TimeSpent       string                 `json:"time_spent"`
	TimeLimit       int                    `json:"time_limit_minutes"`
	SectionBreakdown map[string]SectionScore `json:"section_breakdown,omitempty"`
	ReviewQuestions []ReviewItem           `json:"review_questions,omitempty"`
}

// SectionScore tracks per-section performance
type SectionScore struct {
	Total       int     `json:"total"`
	Correct     int     `json:"correct"`
	Percentage  float64 `json:"percentage"`
}

// ReviewItem for missed questions
type ReviewItem struct {
	QuestionNum   int      `json:"question_num"`
	Question      string   `json:"question"`
	YourAnswer    string   `json:"your_answer"`
	CorrectAnswer string   `json:"correct_answer"`
	Explanation   string   `json:"explanation"`
}

// SubmitAnswerRequest for answering a question
type SubmitAnswerRequest struct {
	SessionID    string `json:"session_id" binding:"required"`
	QuestionID   string `json:"question_id" binding:"required"`
	AnswerIndex  int    `json:"answer_index" binding:"min=0"`
}

// StartTestRequest to begin a new test
type StartTestRequest struct {
	Level   string `json:"level" binding:"required,oneof=N5 N4 N3 N2 N1"`
	Section string `json:"section,omitempty"` // Optional: specific section only
}

// Available JLPT levels and their info
func GetJLPTLevelInfo() []map[string]interface{} {
	return []map[string]interface{}{
		{"level": "N5", "name": "Beginner", "time_limit": 90, "vocab_count": 20, "grammar_count": 10},
		{"level": "N4", "name": "Basic", "time_limit": 105, "vocab_count": 25, "grammar_count": 15},
		{"level": "N3", "name": "Intermediate", "time_limit": 140, "vocab_count": 30, "grammar_count": 20},
		{"level": "N2", "name": "Upper-Intermediate", "time_limit": 155, "vocab_count": 35, "grammar_count": 25},
		{"level": "N1", "name": "Advanced", "time_limit": 170, "vocab_count": 40, "grammar_count": 30},
	}
}
