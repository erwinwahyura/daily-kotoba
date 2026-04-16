package models

import "time"

// ListeningExercise represents a listening comprehension exercise
type ListeningExercise struct {
	ID              string                `json:"id" db:"id"`
	Title           string                `json:"title" db:"title"`
	JLPTLevel       string                `json:"jlpt_level" db:"jlpt_level"`
	Difficulty      string                `json:"difficulty" db:"difficulty"` // easy, medium, hard
	AudioURL        string                `json:"audio_url" db:"audio_url"`
	Duration        int                   `json:"duration" db:"duration"` // seconds
	Transcript      string                `json:"transcript" db:"transcript"`
	Translation     string                `json:"translation" db:"translation"`
	Vocabulary      []VocabItem           `json:"vocabulary" db:"vocabulary"` // Key vocab from audio
	Questions       []ListeningQuestion   `json:"questions" db:"questions"`
	Topic           string                `json:"topic" db:"topic"` // conversation, news, announcement, etc.
	CreatedAt       time.Time             `json:"created_at" db:"created_at"`
}

// VocabItem represents vocabulary from the audio
type VocabItem struct {
	Word        string `json:"word"`
	Reading     string `json:"reading"`
	Meaning     string `json:"meaning"`
	Timestamp   int    `json:"timestamp"` // When it appears in audio (seconds)
}

// ListeningQuestion represents a comprehension question
type ListeningQuestion struct {
	ID          string   `json:"id"`
	Question    string   `json:"question"`
	Options     []string `json:"options"`
	Correct     int      `json:"correct"` // Index of correct answer
	Timestamp   int      `json:"timestamp"` // When to pause (seconds from start)
	Explanation string   `json:"explanation"`
}

// ListeningSession tracks a user's listening practice
type ListeningSession struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	ExerciseID      string    `json:"exercise_id" db:"exercise_id"`
	StartedAt       time.Time `json:"started_at" db:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CurrentPosition int       `json:"current_position" db:"current_position"` // Current playback position
	Answers         []Answer  `json:"answers" db:"answers"` // User's answers
	Score           int       `json:"score" db:"score"`
	Status          string    `json:"status" db:"status"` // in_progress, completed, abandoned
	PlayCount       int       `json:"play_count" db:"play_count"` // How many times audio played
}

// Answer represents a user's answer to a question
type Answer struct {
	QuestionID    string    `json:"question_id"`
	Answer        int       `json:"answer"` // Selected option index
	IsCorrect     bool      `json:"is_correct"`
	Timestamp     time.Time `json:"timestamp"`
	AudioPosition int       `json:"audio_position"` // Where in audio when answered
}

// ListeningProgress tracks user's listening stats
type ListeningProgress struct {
	TotalExercises    int     `json:"total_exercises"`
	CompletedCount    int     `json:"completed_count"`
	AverageScore      float64 `json:"average_score"`
	TotalListeningTime int    `json:"total_listening_time"` // seconds
	ByLevel           map[string]LevelStats `json:"by_level"`
}

// LevelStats represents stats for a JLPT level
type LevelStats struct {
	Completed int     `json:"completed"`
	Average   float64 `json:"average"`
}

// ListeningListResponse for listing exercises
type ListeningListResponse struct {
	Exercises   []ListeningExercise `json:"exercises"`
	TotalCount  int                 `json:"total_count"`
	Level       string              `json:"level,omitempty"`
	Difficulty  string              `json:"difficulty,omitempty"`
}
