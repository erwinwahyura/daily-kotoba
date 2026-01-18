package repository

import (
	"database/sql"
	"fmt"

	"github.com/yourusername/kotoba-api/internal/models"
)

type PlacementRepository struct {
	db *sql.DB
}

func NewPlacementRepository(db *sql.DB) *PlacementRepository {
	return &PlacementRepository{db: db}
}

// GetAllQuestions retrieves all placement test questions
func (r *PlacementRepository) GetAllQuestions() ([]models.PlacementQuestion, error) {
	query := `
		SELECT id, question_text, correct_answer, wrong_answers, difficulty_level, order_index, created_at
		FROM placement_questions
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query placement questions: %w", err)
	}
	defer rows.Close()

	var questions []models.PlacementQuestion
	for rows.Next() {
		var q models.PlacementQuestion
		err := rows.Scan(
			&q.ID,
			&q.QuestionText,
			&q.CorrectAnswer,
			&q.WrongAnswers,
			&q.DifficultyLevel,
			&q.OrderIndex,
			&q.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan placement question: %w", err)
		}
		questions = append(questions, q)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating placement questions: %w", err)
	}

	return questions, nil
}

// GetQuestionByID retrieves a single placement question by ID
func (r *PlacementRepository) GetQuestionByID(questionID string) (*models.PlacementQuestion, error) {
	query := `
		SELECT id, question_text, correct_answer, wrong_answers, difficulty_level, order_index, created_at
		FROM placement_questions
		WHERE id = $1
	`

	var q models.PlacementQuestion
	err := r.db.QueryRow(query, questionID).Scan(
		&q.ID,
		&q.QuestionText,
		&q.CorrectAnswer,
		&q.WrongAnswers,
		&q.DifficultyLevel,
		&q.OrderIndex,
		&q.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("placement question not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query placement question: %w", err)
	}

	return &q, nil
}

// SaveTestResult saves a user's placement test result
func (r *PlacementRepository) SaveTestResult(result *models.PlacementTestResult) error {
	query := `
		INSERT INTO placement_test_results (user_id, test_score, assigned_level, completed_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.db.QueryRow(
		query,
		result.UserID,
		result.TestScore,
		result.AssignedLevel,
		result.CompletedAt,
	).Scan(&result.ID)

	if err != nil {
		return fmt.Errorf("failed to save placement test result: %w", err)
	}

	return nil
}

// GetUserTestResult retrieves a user's most recent placement test result
func (r *PlacementRepository) GetUserTestResult(userID string) (*models.PlacementTestResult, error) {
	query := `
		SELECT id, user_id, test_score, assigned_level, completed_at
		FROM placement_test_results
		WHERE user_id = $1
		ORDER BY completed_at DESC
		LIMIT 1
	`

	var result models.PlacementTestResult
	err := r.db.QueryRow(query, userID).Scan(
		&result.ID,
		&result.UserID,
		&result.TestScore,
		&result.AssignedLevel,
		&result.CompletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil // No test result found
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query placement test result: %w", err)
	}

	return &result, nil
}
