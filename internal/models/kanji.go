package models

import (
	"time"
)

// Kanji represents a kanji character with stroke order
type Kanji struct {
	ID          string   `json:"id" db:"id"`
	Character   string   `json:"character" db:"character"`     // 漢
	JLPTLevel   string   `json:"jlpt_level" db:"jlpt_level"`   // N5, N4, etc.
	Meaning     string   `json:"meaning" db:"meaning"`
	Readings    []string `json:"readings" db:"readings"`     // JSON array
	StrokeCount int      `json:"stroke_count" db:"stroke_count"`
	StrokeOrder []Stroke `json:"stroke_order" db:"stroke_order"` // JSON array of stroke data
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// Stroke represents a single stroke path
type Stroke struct {
	StrokeNum   int       `json:"stroke_num"`
	Path        []Point   `json:"path"`         // Array of x,y coordinates
	Direction   string    `json:"direction"`    // horizontal, vertical, diagonal, curve
	StartPoint  Point     `json:"start_point"`
	EndPoint    Point     `json:"end_point"`
}

// Point represents a coordinate
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// KanjiPracticeSession tracks writing practice
type KanjiPracticeSession struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	KanjiID       string    `json:"kanji_id" db:"kanji_id"`
	KanjiChar     string    `json:"kanji_char" db:"kanji_char"`
	StartedAt     time.Time `json:"started_at" db:"started_at"`
	CompletedAt   *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	UserStrokes   []UserStroke `json:"user_strokes" db:"user_strokes"` // JSON
	Accuracy      float64   `json:"accuracy" db:"accuracy"`       // 0-100
	Status        string    `json:"status" db:"status"`             // in_progress, completed
}

// UserStroke captures user's drawn stroke
type UserStroke struct {
	StrokeNum   int       `json:"stroke_num"`
	Path        []Point   `json:"path"`
	Duration    int       `json:"duration_ms"`
	Timestamp   time.Time `json:"timestamp"`
}

// KanjiCompareRequest for stroke comparison
type KanjiCompareRequest struct {
	SessionID   string      `json:"session_id" binding:"required"`
	StrokeNum   int         `json:"stroke_num" binding:"required"`
	UserPath    []Point     `json:"user_path" binding:"required"`
}

// KanjiCompareResult for stroke feedback
type KanjiCompareResult struct {
	Accuracy      float64 `json:"accuracy"`       // 0-100
	Feedback      string  `json:"feedback"`       // "Good!", "Try again", etc.
	Direction     string  `json:"direction"`      // correct, wrong_direction
	OrderCorrect  bool    `json:"order_correct"`  // Is this the expected stroke?
}

// KanjiListResponse for listing kanji
type KanjiListResponse struct {
	Kanji       []Kanji `json:"kanji"`
	TotalCount  int     `json:"total_count"`
	Level       string  `json:"level"`
}