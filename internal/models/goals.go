package models

import "time"

type UserGoalSettings struct {
	UserID              string    `json:"user_id" db:"user_id"`
	VocabTarget         int       `json:"vocab_target" db:"vocab_target"`
	GrammarTarget       int       `json:"grammar_target" db:"grammar_target"`
	KanjiTarget         int       `json:"kanji_target" db:"kanji_target"`
	ConjugationTarget   int       `json:"conjugation_target" db:"conjugation_target"`
	ReadingTarget       int       `json:"reading_target" db:"reading_target"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

type UserDailyActivity struct {
	ID           string    `json:"id" db:"id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ActivityDate string    `json:"activity_date" db:"activity_date"` // YYYY-MM-DD
	ActivityType string    `json:"activity_type" db:"activity_type"`
	Count        int       `json:"count" db:"count"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type AchievementDef struct {
	ID               string    `json:"id" db:"id"`
	Name             string    `json:"name" db:"name"`
	Description      string    `json:"description" db:"description"`
	Icon             string    `json:"icon" db:"icon"`
	Category         string    `json:"category" db:"category"`
	RequirementType  string    `json:"requirement_type" db:"requirement_type"`
	RequirementValue int       `json:"requirement_value" db:"requirement_value"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

type UserAchievement struct {
	ID            string    `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	AchievementID string    `json:"achievement_id" db:"achievement_id"`
	EarnedAt      time.Time `json:"earned_at" db:"earned_at"`
}

// DailyProgressResponse is what GET /goals/daily returns
type DailyProgressResponse struct {
	VocabCompleted       int `json:"vocab_completed"`
	VocabTarget          int `json:"vocab_target"`
	GrammarCompleted     int `json:"grammar_completed"`
	GrammarTarget        int `json:"grammar_target"`
	KanjiCompleted       int `json:"kanji_completed"`
	KanjiTarget          int `json:"kanji_target"`
	ConjugationCompleted int `json:"conjugation_completed"`
	ConjugationTarget    int `json:"conjugation_target"`
	ReadingCompleted     int `json:"reading_completed"`
	ReadingTarget        int `json:"reading_target"`
}

// StreakResponse is what GET /goals/streak returns
type StreakResponse struct {
	CurrentStreak int    `json:"current_streak"`
	LongestStreak int    `json:"longest_streak"`
	LastStudyDate string `json:"last_study_date"`
}

// WeeklyDay represents one day in the weekly heatmap
type WeeklyDay struct {
	Date            string `json:"date"`
	OverallProgress int    `json:"overall_progress"` // 0-100
	IsCompleted     bool   `json:"is_completed"`
}

// RecordActivityRequest body for POST /goals/progress
type RecordActivityRequest struct {
	ActivityType string `json:"activity_type" binding:"required"`
	Count        int    `json:"count"`
}
