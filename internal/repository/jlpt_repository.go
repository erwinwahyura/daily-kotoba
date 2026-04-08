package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type JLPTRepository struct {
	db *db.DB
}

func NewJLPTRepository(database *db.DB) *JLPTRepository {
	return &JLPTRepository{db: database}
}

// GetTestsByLevel retrieves available tests for a JLPT level
func (r *JLPTRepository) GetTestsByLevel(level string) ([]models.JLPTTest, error) {
	query := `SELECT id, level, section, title, description, time_limit_minutes, total_questions, passing_score, created_at FROM jlpt_tests WHERE level = $1`
	rows, err := r.db.Query(query, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []models.JLPTTest
	for rows.Next() {
		var t models.JLPTTest
		err := rows.Scan(&t.ID, &t.Level, &t.Section, &t.Title, &t.Description, &t.TimeLimit, &t.TotalQuestions, &t.PassingScore, &t.CreatedAt)
		if err != nil {
			continue
		}
		tests = append(tests, t)
	}
	return tests, rows.Err()
}

// GetTestByID retrieves a specific test
func (r *JLPTRepository) GetTestByID(testID string) (*models.JLPTTest, error) {
	var t models.JLPTTest
	query := `SELECT id, level, section, title, description, time_limit_minutes, total_questions, passing_score, created_at FROM jlpt_tests WHERE id = $1`
	err := r.db.QueryRow(query, testID).Scan(&t.ID, &t.Level, &t.Section, &t.Title, &t.Description, &t.TimeLimit, &t.TotalQuestions, &t.PassingScore, &t.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("test not found")
	}
	return &t, err
}

// GetQuestionsByTestID retrieves all questions for a test
func (r *JLPTRepository) GetQuestionsByTestID(testID string) ([]models.JLPTQuestion, error) {
	query := `SELECT id, test_id, question_num, type, question, question_reading, english_prompt, options, correct_index, explanation, point_value, skill_tested FROM jlpt_questions WHERE test_id = $1 ORDER BY question_num`
	rows, err := r.db.Query(query, testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.JLPTQuestion
	for rows.Next() {
		var q models.JLPTQuestion
		var optionsJSON []byte
		err := rows.Scan(&q.ID, &q.TestID, &q.QuestionNum, &q.Type, &q.Question, &q.QuestionReading, &q.EnglishPrompt, &optionsJSON, &q.CorrectIndex, &q.Explanation, &q.PointValue, &q.SkillTested)
		if err != nil {
			continue
		}
		json.Unmarshal(optionsJSON, &q.Options)
		questions = append(questions, q)
	}
	return questions, rows.Err()
}

// CreateTestSession starts a new test session
func (r *JLPTRepository) CreateTestSession(session *models.UserTestSession) error {
	query := `INSERT INTO user_test_sessions (id, user_id, test_id, level, started_at, status, answers) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	answersJSON, _ := json.Marshal(session.Answers)
	_, err := r.db.Exec(query, session.ID, session.UserID, session.TestID, session.Level, session.StartedAt, session.Status, string(answersJSON))
	return err
}

// GetTestSession retrieves an active session
func (r *JLPTRepository) GetTestSession(sessionID string) (*models.UserTestSession, error) {
	var s models.UserTestSession
	var answersJSON string
	query := `SELECT id, user_id, test_id, level, started_at, completed_at, time_spent_sec, answers, score, correct_count, status FROM user_test_sessions WHERE id = $1`
	err := r.db.QueryRow(query, sessionID).Scan(&s.ID, &s.UserID, &s.TestID, &s.Level, &s.StartedAt, &s.CompletedAt, &s.TimeSpentSec, &answersJSON, &s.Score, &s.CorrectCount, &s.Status)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("session not found")
	}
	if err != nil {
		return nil, err
	}
	json.Unmarshal([]byte(answersJSON), &s.Answers)
	return &s, nil
}

// UpdateSessionAnswer saves an answer and updates progress
func (r *JLPTRepository) UpdateSessionAnswer(sessionID string, questionID string, answerIndex int) error {
	session, err := r.GetTestSession(sessionID)
	if err != nil {
		return err
	}
	if session.Answers == nil {
		session.Answers = make(map[string]int)
	}
	session.Answers[questionID] = answerIndex
	
	answersJSON, _ := json.Marshal(session.Answers)
	query := `UPDATE user_test_sessions SET answers = $1 WHERE id = $2`
	_, err = r.db.Exec(query, string(answersJSON), sessionID)
	return err
}

// CompleteTestSession marks session complete with score
func (r *JLPTRepository) CompleteTestSession(sessionID string, score, correctCount, timeSpent int) error {
	now := time.Now()
	query := `UPDATE user_test_sessions SET status = 'completed', completed_at = $1, score = $2, correct_count = $3, time_spent_sec = $4 WHERE id = $5`
	_, err := r.db.Exec(query, now, score, correctCount, timeSpent, sessionID)
	return err
}

// GetUserTestHistory gets all tests taken by user
func (r *JLPTRepository) GetUserTestHistory(userID string) ([]models.UserTestSession, error) {
	query := `SELECT id, test_id, level, started_at, completed_at, score, correct_count, status FROM user_test_sessions WHERE user_id = $1 ORDER BY started_at DESC`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []models.UserTestSession
	for rows.Next() {
		var s models.UserTestSession
		err := rows.Scan(&s.ID, &s.TestID, &s.Level, &s.StartedAt, &s.CompletedAt, &s.Score, &s.CorrectCount, &s.Status)
		if err != nil {
			continue
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}