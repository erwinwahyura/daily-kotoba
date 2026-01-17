package models

import "time"

type UserProgress struct {
	UserID            string    `json:"user_id" db:"user_id"`
	CurrentVocabIndex int       `json:"current_vocab_index" db:"current_vocab_index"`
	LastWordID        *string   `json:"last_word_id" db:"last_word_id"`
	StreakDays        int       `json:"streak_days" db:"streak_days"`
	LastStudyDate     *time.Time `json:"last_study_date" db:"last_study_date"`
	WordsLearnedCount int       `json:"words_learned_count" db:"words_learned_count"`
	WordsSkippedCount int       `json:"words_skipped_count" db:"words_skipped_count"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type UserVocabStatus struct {
	ID       string    `json:"id" db:"id"`
	UserID   string    `json:"user_id" db:"user_id"`
	VocabID  string    `json:"vocab_id" db:"vocab_id"`
	Status   string    `json:"status" db:"status"` // 'learning', 'known', 'skipped'
	MarkedAt time.Time `json:"marked_at" db:"marked_at"`
}

type ProgressStats struct {
	CurrentVocabIndex   int               `json:"current_vocab_index"`
	CurrentLevel        string            `json:"current_level"`
	StreakDays          int               `json:"streak_days"`
	LastStudyDate       *time.Time        `json:"last_study_date"`
	WordsLearned        int               `json:"words_learned"`
	WordsSkipped        int               `json:"words_skipped"`
	TotalWordsInLevel   int               `json:"total_words_in_level"`
	TotalDaysActive     int               `json:"total_days_active"`
	WordsLearnedByLevel map[string]int    `json:"words_learned_by_level"`
}
