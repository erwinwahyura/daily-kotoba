package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type ProgressRepository struct {
	db *db.DB
}

func NewProgressRepository(db *db.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

func (r *ProgressRepository) GetByUserID(userID string) (*models.UserProgress, error) {
	progress := &models.UserProgress{}
	query := `
		SELECT user_id, current_vocab_index, current_grammar_index, last_word_id, last_grammar_id,
		       streak_days, last_study_date, words_learned_count, words_skipped_count,
		       grammar_learned_count, updated_at
		FROM user_progress
		WHERE user_id = $1
	`
	err := r.db.QueryRow(query, userID).Scan(
		&progress.UserID,
		&progress.CurrentVocabIndex,
		&progress.CurrentGrammarIndex,
		&progress.LastWordID,
		&progress.LastGrammarID,
		&progress.StreakDays,
		&progress.LastStudyDate,
		&progress.WordsLearnedCount,
		&progress.WordsSkippedCount,
		&progress.GrammarLearnedCount,
		&progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user progress not found")
	}
	if err != nil {
		return nil, err
	}

	return progress, nil
}

func (r *ProgressRepository) Update(progress *models.UserProgress) error {
	query := `
		UPDATE user_progress
		SET current_vocab_index = $1,
		    current_grammar_index = $2,
		    last_word_id = $3,
		    last_grammar_id = $4,
		    streak_days = $5,
		    last_study_date = $6,
		    words_learned_count = $7,
		    words_skipped_count = $8,
		    grammar_learned_count = $9,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $10
	`
	_, err := r.db.Exec(
		query,
		progress.CurrentVocabIndex,
		progress.CurrentGrammarIndex,
		progress.LastWordID,
		progress.LastGrammarID,
		progress.StreakDays,
		progress.LastStudyDate,
		progress.WordsLearnedCount,
		progress.WordsSkippedCount,
		progress.GrammarLearnedCount,
		progress.UserID,
	)
	return err
}

func (r *ProgressRepository) IncrementVocabIndex(userID string) (*models.UserProgress, error) {
	query := `
		UPDATE user_progress
		SET current_vocab_index = current_vocab_index + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
		RETURNING user_id, current_vocab_index, last_word_id, streak_days,
		          last_study_date, words_learned_count, words_skipped_count, updated_at
	`
	progress := &models.UserProgress{}
	err := r.db.QueryRow(query, userID).Scan(
		&progress.UserID,
		&progress.CurrentVocabIndex,
		&progress.LastWordID,
		&progress.StreakDays,
		&progress.LastStudyDate,
		&progress.WordsLearnedCount,
		&progress.WordsSkippedCount,
		&progress.UpdatedAt,
	)
	return progress, err
}

func (r *ProgressRepository) UpdateStreak(userID string, streakDays int, lastStudyDate time.Time) error {
	query := `
		UPDATE user_progress
		SET streak_days = $1,
		    last_study_date = $2,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $3
	`
	_, err := r.db.Exec(query, streakDays, lastStudyDate, userID)
	return err
}

func (r *ProgressRepository) MarkVocabStatus(userID, vocabID, status string) error {
	query := `
		INSERT INTO user_vocab_status (user_id, vocab_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, vocab_id)
		DO UPDATE SET status = EXCLUDED.status, marked_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, vocabID, status)
	return err
}

func (r *ProgressRepository) GetVocabStatusCount(userID, status string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_vocab_status WHERE user_id = $1 AND status = $2`
	err := r.db.QueryRow(query, userID, status).Scan(&count)
	return count, err
}

func (r *ProgressRepository) IncrementWordsLearned(userID string) error {
	query := `
		UPDATE user_progress
		SET words_learned_count = words_learned_count + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *ProgressRepository) IncrementWordsSkipped(userID string) error {
	query := `
		UPDATE user_progress
		SET words_skipped_count = words_skipped_count + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}

// Grammar progress methods

func (r *ProgressRepository) IncrementGrammarIndex(userID string) (*models.UserProgress, error) {
	query := `
		UPDATE user_progress
		SET current_grammar_index = current_grammar_index + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
		RETURNING user_id, current_vocab_index, current_grammar_index, last_word_id, last_grammar_id,
		          streak_days, last_study_date, words_learned_count, words_skipped_count,
		          grammar_learned_count, updated_at
	`
	progress := &models.UserProgress{}
	err := r.db.QueryRow(query, userID).Scan(
		&progress.UserID,
		&progress.CurrentVocabIndex,
		&progress.CurrentGrammarIndex,
		&progress.LastWordID,
		&progress.LastGrammarID,
		&progress.StreakDays,
		&progress.LastStudyDate,
		&progress.WordsLearnedCount,
		&progress.WordsSkippedCount,
		&progress.GrammarLearnedCount,
		&progress.UpdatedAt,
	)
	return progress, err
}

func (r *ProgressRepository) MarkGrammarStatus(userID, grammarID, status string) error {
	query := `
		INSERT INTO user_grammar_status (user_id, grammar_id, status)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, grammar_id)
		DO UPDATE SET status = EXCLUDED.status, marked_at = CURRENT_TIMESTAMP
	`
	_, err := r.db.Exec(query, userID, grammarID, status)
	return err
}

func (r *ProgressRepository) IncrementGrammarLearned(userID string) error {
	query := `
		UPDATE user_progress
		SET grammar_learned_count = grammar_learned_count + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err := r.db.Exec(query, userID)
	return err
}
