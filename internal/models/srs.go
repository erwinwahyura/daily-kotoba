package models

import (
	"time"
)

// SRSSchedule represents a user's spaced repetition schedule for an item
type SRSSchedule struct {
	ID             string    `json:"id" db:"id"`
	UserID         string    `json:"user_id" db:"user_id"`
	ItemID         string    `json:"item_id" db:"item_id"`
	ItemType       string    `json:"item_type" db:"item_type"` // "vocabulary" or "grammar"
	IntervalDays   int       `json:"interval_days" db:"interval_days"`
	Repetitions    int       `json:"repetitions" db:"repetitions"`
	EaseFactor     float64   `json:"ease_factor" db:"ease_factor"`
	LastReviewedAt *time.Time `json:"last_reviewed_at" db:"last_reviewed_at"`
	NextReviewAt   time.Time `json:"next_review_at" db:"next_review_at"`
	TotalReviews   int       `json:"total_reviews" db:"total_reviews"`
	CorrectReviews int     `json:"correct_reviews" db:"correct_reviews"`
	Streak         int       `json:"streak" db:"streak"`
	Status         string    `json:"status" db:"status"` // "learning", "review", "mastered", "lapsed"
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// SRSReviewHistory tracks individual review attempts
type SRSReviewHistory struct {
	ID                string    `json:"id" db:"id"`
	ScheduleID        string    `json:"schedule_id" db:"schedule_id"`
	UserID            string    `json:"user_id" db:"user_id"`
	Quality           int       `json:"quality" db:"quality"` // 0-5 (SM-2)
	ResponseTimeMs    int       `json:"response_time_ms" db:"response_time_ms"`
	ItemType          string    `json:"item_type" db:"item_type"`
	ItemID            string    `json:"item_id" db:"item_id"`
	IntervalBefore    int       `json:"interval_before" db:"interval_before"`
	IntervalAfter     int       `json:"interval_after" db:"interval_after"`
	EaseFactorBefore  float64   `json:"ease_factor_before" db:"ease_factor_before"`
	EaseFactorAfter   float64   `json:"ease_factor_after" db:"ease_factor_after"`
	ReviewedAt        time.Time `json:"reviewed_at" db:"reviewed_at"`
}

// StudySession tracks daily study activity
type StudySession struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	SessionDate     string    `json:"session_date" db:"session_date"` // YYYY-MM-DD
	NewItems        int       `json:"new_items" db:"new_items"`
	ReviewItems     int       `json:"review_items" db:"review_items"`
	DrillItems      int       `json:"drill_items" db:"drill_items"`
	TotalTimeSeconds int      `json:"total_time_seconds" db:"total_time_seconds"`
	CorrectCount    int       `json:"correct_count" db:"correct_count"`
	TotalAttempts   int       `json:"total_attempts" db:"total_attempts"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

// SRSReviewRequest is sent by client when reviewing an item
type SRSReviewRequest struct {
	ItemID         string `json:"item_id" binding:"required"`
	ItemType       string `json:"item_type" binding:"required,oneof=vocabulary grammar"`
	Quality        int    `json:"quality" binding:"required,min=0,max=5"` // 0-5 SM-2 rating
	ResponseTimeMs int    `json:"response_time_ms"`                      // Optional
}

// SRSReviewResponse contains the result and next review info
type SRSReviewResponse struct {
	Schedule        *SRSSchedule `json:"schedule"`
	NextReview      string       `json:"next_review"` // Human readable: "1 day", "3 days"
	StatusChanged   bool         `json:"status_changed"`
	NewAchievement  *Achievement `json:"new_achievement,omitempty"`
}

// SRSQueueResponse returns items due for review
type SRSQueueResponse struct {
	DueItems        []SRSDueItem `json:"due_items"`
	TotalDue        int          `json:"total_due"`
	NewItems        int          `json:"new_items"`
	LearningItems   int          `json:"learning_items"`
	ReviewItems     int          `json:"review_items"`
	MasteredItems   int          `json:"mastered_items"`
}

// SRSDueItem is a single item ready for review
type SRSDueItem struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"` // "vocabulary" or "grammar"
	Data           interface{} `json:"data"` // Vocabulary or GrammarPattern
	Schedule       *SRSSchedule `json:"schedule"`
	DaysOverdue    int         `json:"days_overdue"`
}

// SRSStats provides overview statistics
type SRSStats struct {
	TotalItems      int     `json:"total_items"`
	Learning        int     `json:"learning"`        // 0-2 repetitions
	Review          int     `json:"review"`          // 3+ repetitions, due soon
	Mastered        int     `json:"mastered"`        // 8+ repetitions
	Lapsed          int     `json:"lapsed"`          // Failed reviews
	DueToday        int     `json:"due_today"`
	DueTomorrow     int     `json:"due_tomorrow"`
	CurrentStreak   int     `json:"current_streak"`
	LongestStreak   int     `json:"longest_streak"`
	TotalReviews    int     `json:"total_reviews"`
	Accuracy        float64 `json:"accuracy"`      // 0.0-1.0
}



// SM2Result holds the calculation result from SM-2 algorithm
type SM2Result struct {
	Interval   int     `json:"interval"`
	Repetitions int    `json:"repetitions"`
	EaseFactor float64 `json:"ease_factor"`
	Status     string  `json:"status"`
}

// CalculateSM2 implements the SuperMemo-2 algorithm
func CalculateSM2(quality int, interval int, repetitions int, easeFactor float64) *SM2Result {
	result := &SM2Result{}
	
	// Quality 0-2: Failed review
	if quality < 3 {
		result.Repetitions = 0
		result.Interval = 1 // Review again tomorrow
		result.EaseFactor = max(1.3, easeFactor-0.2)
		result.Status = "learning"
		return result
	}
	
	// Quality 3-5: Successful review
	result.Repetitions = repetitions + 1
	
	// Calculate interval
	if result.Repetitions == 1 {
		result.Interval = 1
	} else if result.Repetitions == 2 {
		result.Interval = 6
	} else {
		result.Interval = int(float64(interval) * easeFactor)
	}
	
	// Update ease factor
	result.EaseFactor = easeFactor + (0.1 - float64(5-quality)*(0.08+float64(5-quality)*0.02))
	result.EaseFactor = max(1.3, result.EaseFactor)
	
	// Determine status
	if result.Repetitions >= 8 {
		result.Status = "mastered"
	} else {
		result.Status = "review"
	}
	
	return result
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}