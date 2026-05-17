package models

import "time"

// DailyGoal represents a user's daily study target
type DailyGoal struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	Date            time.Time `json:"date" db:"date"` // Date for this goal (YYYY-MM-DD)
	VocabTarget     int       `json:"vocab_target" db:"vocab_target"`     // Target words to learn
	VocabCompleted  int       `json:"vocab_completed" db:"vocab_completed"`
	GrammarTarget int       `json:"grammar_target" db:"grammar_target"` // Target patterns to learn
	GrammarCompleted int    `json:"grammar_completed" db:"grammar_completed"`
	KanjiTarget     int       `json:"kanji_target" db:"kanji_target"`     // Target kanji to practice
	KanjiCompleted  int       `json:"kanji_completed" db:"kanji_completed"`
	ConjugationTarget int   `json:"conjugation_target" db:"conjugation_target"` // Target conjugation drills
	ConjugationCompleted int `json:"conjugation_completed" db:"conjugation_completed"`
	ReadingTarget   int       `json:"reading_target" db:"reading_target"`   // Target reading articles
	ReadingCompleted int    `json:"reading_completed" db:"reading_completed"`
	IsCompleted     bool      `json:"is_completed" db:"is_completed"`
	CompletedAt     *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}

// UserStreak tracks consecutive days of activity
type UserStreak struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	CurrentStreak    int     `json:"current_streak" db:"current_streak"`
	LongestStreak     int     `json:"longest_streak" db:"longest_streak"`
	LastActivityDate  *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
	TotalActiveDays   int     `json:"total_active_days" db:"total_active_days"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// GoalSettings represents user's default daily targets
type GoalSettings struct {
	ID               string    `json:"id" db:"id"`
	UserID           string    `json:"user_id" db:"user_id"`
	VocabTarget      int       `json:"vocab_target" db:"vocab_target"`      // Default: 10
	GrammarTarget    int     `json:"grammar_target" db:"grammar_target"`   // Default: 5
	KanjiTarget      int       `json:"kanji_target" db:"kanji_target"`       // Default: 5
	ConjugationTarget int      `json:"conjugation_target" db:"conjugation_target"` // Default: 20
	ReadingTarget    int       `json:"reading_target" db:"reading_target"`   // Default: 1
	EnableReminders  bool      `json:"enable_reminders" db:"enable_reminders"`
	ReminderTime     string    `json:"reminder_time" db:"reminder_time"`     // "20:00"
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// DailyProgress represents today's progress summary
type DailyProgress struct {
	Date              string `json:"date"`
	VocabCompleted    int    `json:"vocab_completed"`
	VocabTarget       int    `json:"vocab_target"`
	GrammarCompleted  int    `json:"grammar_completed"`
	GrammarTarget     int    `json:"grammar_target"`
	KanjiCompleted    int    `json:"kanji_completed"`
	KanjiTarget       int    `json:"kanji_target"`
	ConjugationCompleted int `json:"conjugation_completed"`
	ConjugationTarget    int `json:"conjugation_target"`
	ReadingCompleted  int    `json:"reading_completed"`
	ReadingTarget     int    `json:"reading_target"`
	OverallProgress   int    `json:"overall_progress"` // 0-100
	IsCompleted       bool   `json:"is_completed"`
	CurrentStreak     int    `json:"current_streak"`
}

// Achievement represents a badge/achievement
type Achievement struct {
	ID          string    `json:"id" db:"id"`
	UserID      string    `json:"user_id" db:"user_id"`
	Type        string    `json:"type" db:"type"` // streak, vocab, grammar, kanji, conjugation, reading
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Icon        string    `json:"icon" db:"icon"` // Emoji or icon name
	Level       int       `json:"level" db:"level"` // 1, 2, 3 for tiered achievements
	EarnedAt    time.Time `json:"earned_at" db:"earned_at"`
	IsNew       bool      `json:"is_new,omitempty" db:"-"` // Not stored, for notifications
}

// AchievementDefinition defines available achievements
type AchievementDefinition struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Level       int    `json:"level"`
	Requirement int    `json:"requirement"` // Number to achieve (e.g., 7 for 7-day streak)
}

// GetAchievementDefinitions returns all available achievements
func GetAchievementDefinitions() []AchievementDefinition {
	return []AchievementDefinition{
		// Streak achievements
		{ID: "streak_3", Type: "streak", Name: "Getting Started", Description: "3-day study streak", Icon: "🔥", Level: 1, Requirement: 3},
		{ID: "streak_7", Type: "streak", Name: "Week Warrior", Description: "7-day study streak", Icon: "🔥", Level: 2, Requirement: 7},
		{ID: "streak_14", Type: "streak", Name: "Fortnight Focus", Description: "14-day study streak", Icon: "🔥", Level: 3, Requirement: 14},
		{ID: "streak_30", Type: "streak", Name: "Monthly Master", Description: "30-day study streak", Icon: "🔥", Level: 4, Requirement: 30},
		{ID: "streak_100", Type: "streak", Name: "Century Scholar", Description: "100-day study streak", Icon: "👑", Level: 5, Requirement: 100},
		
		// Vocabulary achievements
		{ID: "vocab_50", Type: "vocab", Name: "Word Collector", Description: "Learn 50 words", Icon: "📚", Level: 1, Requirement: 50},
		{ID: "vocab_100", Type: "vocab", Name: "Vocabulary Builder", Description: "Learn 100 words", Icon: "📚", Level: 2, Requirement: 100},
		{ID: "vocab_500", Type: "vocab", Name: "Word Master", Description: "Learn 500 words", Icon: "📚", Level: 3, Requirement: 500},
		{ID: "vocab_1000", Type: "vocab", Name: "Vocabulary Expert", Description: "Learn 1,000 words", Icon: "📚", Level: 4, Requirement: 1000},
		
		// Grammar achievements
		{ID: "grammar_10", Type: "grammar", Name: "Pattern Beginner", Description: "Learn 10 grammar patterns", Icon: "📝", Level: 1, Requirement: 10},
		{ID: "grammar_50", Type: "grammar", Name: "Pattern Student", Description: "Learn 50 grammar patterns", Icon: "📝", Level: 2, Requirement: 50},
		{ID: "grammar_100", Type: "grammar", Name: "Grammar Expert", Description: "Learn 100 grammar patterns", Icon: "📝", Level: 3, Requirement: 100},
		
		// Kanji achievements
		{ID: "kanji_50", Type: "kanji", Name: "Kanji Novice", Description: "Practice 50 kanji", Icon: "漢", Level: 1, Requirement: 50},
		{ID: "kanji_100", Type: "kanji", Name: "Kanji Learner", Description: "Practice 100 kanji", Icon: "漢", Level: 2, Requirement: 100},
		{ID: "kanji_500", Type: "kanji", Name: "Kanji Master", Description: "Practice 500 kanji", Icon: "漢", Level: 3, Requirement: 500},
		
		// Conjugation achievements
		{ID: "conj_100", Type: "conjugation", Name: "Conjugation Beginner", Description: "Complete 100 conjugation drills", Icon: "活", Level: 1, Requirement: 100},
		{ID: "conj_500", Type: "conjugation", Name: "Conjugation Expert", Description: "Complete 500 conjugation drills", Icon: "活", Level: 2, Requirement: 500},
		{ID: "conj_1000", Type: "conjugation", Name: "Conjugation Master", Description: "Complete 1,000 conjugation drills", Icon: "活", Level: 3, Requirement: 1000},
		
		// Level achievements
		{ID: "level_n5", Type: "level", Name: "N5 Complete", Description: "Complete N5 level", Icon: "🎓", Level: 1, Requirement: 1},
		{ID: "level_n4", Type: "level", Name: "N4 Complete", Description: "Complete N4 level", Icon: "🎓", Level: 2, Requirement: 1},
		{ID: "level_n3", Type: "level", Name: "N3 Complete", Description: "Complete N3 level", Icon: "🎓", Level: 3, Requirement: 1},
		{ID: "level_n2", Type: "level", Name: "N2 Complete", Description: "Complete N2 level", Icon: "🎓", Level: 4, Requirement: 1},
		{ID: "level_n1", Type: "level", Name: "N1 Complete", Description: "Complete N1 level", Icon: "🎓", Level: 5, Requirement: 1},
	}
}
