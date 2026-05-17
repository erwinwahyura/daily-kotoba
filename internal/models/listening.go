package models

import "time"

type ListeningExercise struct {
	ID              string    `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	JLPTLevel       string    `json:"jlpt_level" db:"jlpt_level"`
	Difficulty      string    `json:"difficulty" db:"difficulty"`
	Topic           string    `json:"topic" db:"topic"`
	Transcript      string    `json:"transcript" db:"transcript"`
	Translation     string    `json:"translation" db:"translation"`
	DurationSeconds int       `json:"duration" db:"duration_seconds"`
	TTSCacheID      *string   `json:"tts_cache_id,omitempty" db:"tts_cache_id"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type ListeningVocabulary struct {
	ID         string `json:"id" db:"id"`
	ExerciseID string `json:"exercise_id" db:"exercise_id"`
	Word       string `json:"word" db:"word"`
	Reading    string `json:"reading" db:"reading"`
	Meaning    string `json:"meaning" db:"meaning"`
}

type ListeningQuestion struct {
	ID          string   `json:"id" db:"id"`
	ExerciseID  string   `json:"exercise_id" db:"exercise_id"`
	QuestionNum int      `json:"question_num" db:"question_num"`
	Question    string   `json:"question" db:"question"`
	Options     []string `json:"options" db:"-"`
	OptionsJSON string   `json:"-" db:"options"`
	CorrectIdx  int      `json:"correct_index" db:"correct_index"`
}

type ListeningSession struct {
	ID             string     `json:"session_id" db:"id"`
	UserID         string     `json:"user_id" db:"user_id"`
	ExerciseID     string     `json:"exercise_id" db:"exercise_id"`
	StartedAt      time.Time  `json:"started_at" db:"started_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	Score          int        `json:"score" db:"score"`
	TotalQuestions int        `json:"total_questions" db:"total_questions"`
	Status         string     `json:"status" db:"status"`
}

type ListeningSessionAnswer struct {
	ID         string    `json:"id" db:"id"`
	SessionID  string    `json:"session_id" db:"session_id"`
	QuestionID string    `json:"question_id" db:"question_id"`
	AnswerIdx  int       `json:"answer_index" db:"answer_index"`
	IsCorrect  bool      `json:"correct" db:"is_correct"`
	AnsweredAt time.Time `json:"answered_at" db:"answered_at"`
}

// ListeningExerciseDetail is the full payload for a single exercise
type ListeningExerciseDetail struct {
	ListeningExercise
	Questions  []ListeningQuestion  `json:"questions"`
	Vocabulary []ListeningVocabulary `json:"vocabulary"`
	AudioURL   string               `json:"audio_url"`
}

// ListeningExerciseSummary is used in the exercises list
type ListeningExerciseSummary struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
	Duration   int    `json:"duration"`
	Topic      string `json:"topic"`
}

type SubmitListeningAnswerRequest struct {
	SessionID  string `json:"session_id" binding:"required"`
	QuestionID string `json:"question_id" binding:"required"`
	Answer     int    `json:"answer"`
}

type SubmitListeningAnswerResponse struct {
	Correct       bool `json:"correct"`
	CorrectAnswer int  `json:"correct_answer"`
}
