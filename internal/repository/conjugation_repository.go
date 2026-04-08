package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type ConjugationRepository struct {
	db *db.DB
}

func NewConjugationRepository(database *db.DB) *ConjugationRepository {
	return &ConjugationRepository{db: database}
}

// GetChallengeByID retrieves a specific conjugation challenge
func (r *ConjugationRepository) GetChallengeByID(id string) (*models.ConjugationChallenge, error) {
	challenge := &models.ConjugationChallenge{}
	query := `
		SELECT id, base_form, reading, "group", target_form, target_ending, 
		       full_answer, hint, difficulty, jlpt_level, category, created_at
		FROM conjugation_challenges
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&challenge.ID, &challenge.BaseForm, &challenge.Reading, &challenge.Group,
		&challenge.TargetForm, &challenge.TargetEnding, &challenge.FullAnswer,
		&challenge.Hint, &challenge.Difficulty, &challenge.JLPTLevel,
		&challenge.Category, &challenge.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("challenge not found")
	}
	return challenge, err
}

// GetChallengesByForm retrieves challenges for a specific form type and level
func (r *ConjugationRepository) GetChallengesByForm(formType, jlptLevel string, limit int) ([]*models.ConjugationChallenge, error) {
	if limit < 1 || limit > 50 {
		limit = 20
	}

	query := `
		SELECT id, base_form, reading, "group", target_form, target_ending,
		       full_answer, hint, difficulty, jlpt_level, category, created_at
		FROM conjugation_challenges
		WHERE target_form = $1 AND difficulty = $2
		ORDER BY RANDOM()
		LIMIT $3
	`

	rows, err := r.db.Query(query, formType, jlptLevel, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*models.ConjugationChallenge
	for rows.Next() {
		c := &models.ConjugationChallenge{}
		err := rows.Scan(
			&c.ID, &c.BaseForm, &c.Reading, &c.Group,
			&c.TargetForm, &c.TargetEnding, &c.FullAnswer,
			&c.Hint, &c.Difficulty, &c.JLPTLevel,
			&c.Category, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	return challenges, rows.Err()
}

// GetChallengesByLevel retrieves all challenges for a JLPT level
func (r *ConjugationRepository) GetChallengesByLevel(jlptLevel string) ([]*models.ConjugationChallenge, error) {
	query := `
		SELECT id, base_form, reading, "group", target_form, target_ending,
		       full_answer, hint, difficulty, jlpt_level, category, created_at
		FROM conjugation_challenges
		WHERE difficulty = $1
		ORDER BY target_form, base_form
	`

	rows, err := r.db.Query(query, jlptLevel)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var challenges []*models.ConjugationChallenge
	for rows.Next() {
		c := &models.ConjugationChallenge{}
		err := rows.Scan(
			&c.ID, &c.BaseForm, &c.Reading, &c.Group,
			&c.TargetForm, &c.TargetEnding, &c.FullAnswer,
			&c.Hint, &c.Difficulty, &c.JLPTLevel,
			&c.Category, &c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}

	return challenges, rows.Err()
}

// CreateOrUpdateChallenge inserts or updates a conjugation challenge
func (r *ConjugationRepository) CreateOrUpdateChallenge(challenge *models.ConjugationChallenge) error {
	query := `
		INSERT INTO conjugation_challenges 
		(id, base_form, reading, "group", target_form, target_ending, full_answer, hint, difficulty, jlpt_level, category, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE SET
		base_form = EXCLUDED.base_form,
		reading = EXCLUDED.reading,
		"group" = EXCLUDED.group,
		target_form = EXCLUDED.target_form,
		target_ending = EXCLUDED.target_ending,
		full_answer = EXCLUDED.full_answer,
		hint = EXCLUDED.hint,
		difficulty = EXCLUDED.difficulty,
		jlpt_level = EXCLUDED.jlpt_level,
		category = EXCLUDED.category
	`
	_, err := r.db.Exec(query,
		challenge.ID, challenge.BaseForm, challenge.Reading, challenge.Group,
		challenge.TargetForm, challenge.TargetEnding, challenge.FullAnswer,
		challenge.Hint, challenge.Difficulty, challenge.JLPTLevel,
		challenge.Category, challenge.CreatedAt,
	)
	return err
}

// Session methods

// GetActiveSession retrieves user's active conjugation session
func (r *ConjugationRepository) GetActiveSession(userID string) (*models.ConjugationSession, error) {
	session := &models.ConjugationSession{}
	var completedFormsJSON []byte

	query := `
		SELECT id, user_id, current_form, current_index, total_questions,
		       correct_count, wrong_count, streak, max_streak, start_time, last_active, completed_forms
		FROM conjugation_sessions
		WHERE user_id = $1 AND last_active > $2
		ORDER BY last_active DESC
		LIMIT 1
	`
	// Sessions expire after 24 hours of inactivity
	cutoff := time.Now().Add(-24 * time.Hour)

	err := r.db.QueryRow(query, userID, cutoff).Scan(
		&session.ID, &session.UserID, &session.CurrentForm, &session.CurrentIndex,
		&session.TotalQuestions, &session.CorrectCount, &session.WrongCount,
		&session.Streak, &session.MaxStreak, &session.StartTime,
		&session.LastActive, &completedFormsJSON,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Parse completed forms JSON
	if len(completedFormsJSON) > 0 {
		json.Unmarshal(completedFormsJSON, &session.CompletedForms)
	}

	return session, nil
}

// CreateSession creates a new conjugation session
func (r *ConjugationRepository) CreateSession(session *models.ConjugationSession) error {
	completedFormsJSON, _ := json.Marshal(session.CompletedForms)

	query := `
		INSERT INTO conjugation_sessions
		(id, user_id, current_form, current_index, total_questions, correct_count,
		 wrong_count, streak, max_streak, start_time, last_active, completed_forms)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(query,
		session.ID, session.UserID, session.CurrentForm, session.CurrentIndex,
		session.TotalQuestions, session.CorrectCount, session.WrongCount,
		session.Streak, session.MaxStreak, session.StartTime,
		session.LastActive, completedFormsJSON,
	)
	return err
}

// UpdateSession updates an existing session
func (r *ConjugationRepository) UpdateSession(session *models.ConjugationSession) error {
	completedFormsJSON, _ := json.Marshal(session.CompletedForms)

	query := `
		UPDATE conjugation_sessions
		SET current_form = $1, current_index = $2, total_questions = $3,
		    correct_count = $4, wrong_count = $5, streak = $6, max_streak = $7,
		    last_active = $8, completed_forms = $9
		WHERE id = $10
	`
	_, err := r.db.Exec(query,
		session.CurrentForm, session.CurrentIndex, session.TotalQuestions,
		session.CorrectCount, session.WrongCount, session.Streak, session.MaxStreak,
		session.LastActive, completedFormsJSON, session.ID,
	)
	return err
}

// RecordAttempt records a conjugation attempt
func (r *ConjugationRepository) RecordAttempt(attempt *models.ConjugationAttempt) error {
	query := `
		INSERT INTO conjugation_attempts
		(id, session_id, user_id, challenge_id, form_type, base_form, user_answer,
		 is_correct, time_spent_sec, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query,
		attempt.ID, attempt.SessionID, attempt.UserID, attempt.ChallengeID,
		attempt.FormType, attempt.BaseForm, attempt.UserAnswer,
		attempt.IsCorrect, attempt.TimeSpentSec, attempt.CreatedAt,
	)
	return err
}

// GetProgressStats retrieves user's conjugation progress
func (r *ConjugationRepository) GetProgressStats(userID string) (*models.ConjugationProgress, error) {
	// Calculate accuracy by form type
	query := `
		SELECT form_type, 
		       COUNT(*) as total,
		       SUM(CASE WHEN is_correct THEN 1 ELSE 0 END) as correct
		FROM conjugation_attempts
		WHERE user_id = $1
		GROUP BY form_type
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	formMastery := make(map[string]float64)
	totalAttempts := 0
	correctAttempts := 0

	for rows.Next() {
		var formType string
		var total, correct int
		if err := rows.Scan(&formType, &total, &correct); err != nil {
			continue
		}
		formMastery[formType] = float64(correct) / float64(total) * 100
		totalAttempts += total
		correctAttempts += correct
	}

	accuracy := 0.0
	if totalAttempts > 0 {
		accuracy = float64(correctAttempts) / float64(totalAttempts) * 100
	}

	// Get current streak
	var bestStreak int
	err = r.db.QueryRow(
		"SELECT COALESCE(MAX(max_streak), 0) FROM conjugation_sessions WHERE user_id = $1",
		userID,
	).Scan(&bestStreak)

	return &models.ConjugationProgress{
		FormMastery:     formMastery,
		TotalAttempts:   totalAttempts,
		CorrectAttempts: correctAttempts,
		AccuracyRate:    accuracy,
		BestStreak:      bestStreak,
		DailyGoal:       20,
	}, nil
}
