package models

import "time"

// ConjugationChallenge represents a conjugation drill exercise
type ConjugationChallenge struct {
	ID            string                `json:"id" db:"id"`
	BaseForm      string                `json:"base_form" db:"base_form"`         // 食べる, 行く, etc.
	Reading       string                `json:"reading" db:"reading"`             // たべる, いく
	Group         string                `json:"group" db:"group"`                   // godan (group 1), ichidan (group 2), irregular
	TargetForm    string                `json:"target_form" db:"target_form"`     // te-form, nai-form, ta-form, etc.
	TargetEnding  string                `json:"target_ending" db:"target_ending"` // て, ない, た (expected answer ending)
	FullAnswer    string                `json:"full_answer" db:"full_answer"`     // 食べて, 行かない
	Hint          string                `json:"hint" db:"hint"`                     // Conjugation hint/rule
	Difficulty    string                `json:"difficulty" db:"difficulty"`       // N5, N4, N3, N2, N1
	JLPTLevel     string                `json:"jlpt_level" db:"jlpt_level"`
	Category      string                `json:"category" db:"category"`           // verb, adjective, copula
	CreatedAt     time.Time             `json:"created_at" db:"created_at"`
}

// ConjugationFormType represents different conjugation forms
type ConjugationFormType struct {
	Name        string `json:"name"`        // te-form, nai-form, ta-form, etc.
	DisplayName string `json:"display_name"` // て形, ない形, た形
	Description string `json:"description"`
	Level       string `json:"level"`       // N5, N4, etc.
	Order       int    `json:"order"`       // Learning sequence order
}

// Available conjugation forms in learning order
type ConjugationForms struct {
	Forms []ConjugationFormType
}

func GetConjugationForms() []ConjugationFormType {
	return []ConjugationFormType{
		{Name: "polite", DisplayName: "丁寧形 (ます)", Description: "Polite/formal ending", Level: "N5", Order: 1},
		{Name: "te", DisplayName: "て形", Description: "Connective form", Level: "N4", Order: 2},
		{Name: "ta", DisplayName: "た形", Description: "Past tense", Level: "N4", Order: 3},
		{Name: "nai", DisplayName: "ない形", Description: "Negative plain", Level: "N4", Order: 4},
		{Name: "nakatta", DisplayName: "なかった形", Description: "Negative past", Level: "N4", Order: 5},
		{Name: "potential", DisplayName: "可能形", Description: "Potential (can do)", Level: "N3", Order: 6},
		{Name: "passive", DisplayName: "受身形", Description: "Passive (is done)", Level: "N3", Order: 7},
		{Name: "causative", DisplayName: "使役形", Description: "Causative (make do)", Level: "N3", Order: 8},
		{Name: "imperative", DisplayName: "命令形", Description: "Imperative (command)", Level: "N3", Order: 9},
		{Name: "conditional", DisplayName: "条件形", Description: "Conditional (if)", Level: "N2", Order: 10},
		{Name: "volitional", DisplayName: "意向形", Description: "Volitional (let's)", Level: "N4", Order: 11},
	}
}

// ConjugationSession tracks a user's drill session
type ConjugationSession struct {
	ID              string    `json:"id" db:"id"`
	UserID          string    `json:"user_id" db:"user_id"`
	CurrentForm     string    `json:"current_form" db:"current_form"`
	CurrentIndex    int       `json:"current_index" db:"current_index"`
	TotalQuestions  int       `json:"total_questions" db:"total_questions"`
	CorrectCount    int       `json:"correct_count" db:"correct_count"`
	WrongCount      int       `json:"wrong_count" db:"wrong_count"`
	Streak          int       `json:"streak" db:"streak"`
	MaxStreak       int       `json:"max_streak" db:"max_streak"`
	StartTime       time.Time `json:"start_time" db:"start_time"`
	LastActive      time.Time `json:"last_active" db:"last_active"`
	CompletedForms  []string  `json:"completed_forms" db:"completed_forms"` // JSON array
}

// ConjugationAttempt records a single attempt
type ConjugationAttempt struct {
	ID           string    `json:"id" db:"id"`
	SessionID    string    `json:"session_id" db:"session_id"`
	UserID       string    `json:"user_id" db:"user_id"`
	ChallengeID  string    `json:"challenge_id" db:"challenge_id"`
	FormType     string    `json:"form_type" db:"form_type"`
	BaseForm     string    `json:"base_form" db:"base_form"`
	UserAnswer   string    `json:"user_answer" db:"user_answer"`
	IsCorrect    bool      `json:"is_correct" db:"is_correct"`
	TimeSpentSec int       `json:"time_spent_sec" db:"time_spent_sec"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ConjugationChallengeResponse for API
type ConjugationChallengeResponse struct {
	Challenge   *ConjugationChallenge `json:"challenge"`
	Progress    *ConjugationProgress   `json:"progress"`
	FormInfo    *ConjugationFormType    `json:"form_info"`
	SessionID   string                 `json:"session_id"`
}

// ConjugationProgress tracks overall progress
type ConjugationProgress struct {
	CurrentForm      string             `json:"current_form"`
	FormsUnlocked    []string           `json:"forms_unlocked"`
	FormsCompleted   []string           `json:"forms_completed"`
	FormMastery      map[string]float64 `json:"form_mastery"` // percentage
	TotalAttempts    int                `json:"total_attempts"`
	CorrectAttempts  int                `json:"correct_attempts"`
	AccuracyRate     float64            `json:"accuracy_rate"`
	CurrentStreak    int                `json:"current_streak"`
	BestStreak       int                `json:"best_streak"`
	DailyGoal        int                `json:"daily_goal"`
	DailyCompleted   int                `json:"daily_completed"`
}

// ConjugationSubmitRequest for answer submission
type ConjugationSubmitRequest struct {
	SessionID   string `json:"session_id" binding:"required"`
	ChallengeID string `json:"challenge_id" binding:"required"`
	Answer      string `json:"answer" binding:"required"`
	TimeSpentMs int    `json:"time_spent_ms"`
}

// ConjugationSubmitResponse after submission
type ConjugationSubmitResponse struct {
	IsCorrect       bool                   `json:"is_correct"`
	CorrectAnswer   string                 `json:"correct_answer"`
	Explanation     string                 `json:"explanation"`
	NextChallenge   *ConjugationChallenge   `json:"next_challenge,omitempty"`
	NextFormInfo    *ConjugationFormType    `json:"next_form_info,omitempty"`
	FormCompleted   bool                   `json:"form_completed"`
	AllFormsCompleted bool                `json:"all_forms_completed"`
	Progress        *ConjugationProgress   `json:"progress"`
	SessionID       string                 `json:"session_id"`
}
