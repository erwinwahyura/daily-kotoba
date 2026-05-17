package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type ListeningRepository struct {
	db *db.DB
}

func NewListeningRepository(db *db.DB) *ListeningRepository {
	return &ListeningRepository{db: db}
}

func (r *ListeningRepository) GetExercisesByLevel(level string) ([]models.ListeningExerciseSummary, error) {
	rows, err := r.db.Query(
		"SELECT id, title, difficulty, duration_seconds, topic FROM listening_exercises WHERE jlpt_level = "+r.db.Placeholder(1)+" ORDER BY difficulty",
		level,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ListeningExerciseSummary
	for rows.Next() {
		var ex models.ListeningExerciseSummary
		if err := rows.Scan(&ex.ID, &ex.Title, &ex.Difficulty, &ex.Duration, &ex.Topic); err != nil {
			continue
		}
		list = append(list, ex)
	}
	return list, rows.Err()
}

func (r *ListeningRepository) GetExercise(id string) (*models.ListeningExercise, error) {
	ex := &models.ListeningExercise{}
	err := r.db.QueryRow(
		"SELECT id, title, jlpt_level, difficulty, topic, transcript, translation, duration_seconds, tts_cache_id, created_at FROM listening_exercises WHERE id = "+r.db.Placeholder(1),
		id,
	).Scan(&ex.ID, &ex.Title, &ex.JLPTLevel, &ex.Difficulty, &ex.Topic, &ex.Transcript, &ex.Translation, &ex.DurationSeconds, &ex.TTSCacheID, &ex.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return ex, err
}

func (r *ListeningRepository) GetQuestions(exerciseID string) ([]models.ListeningQuestion, error) {
	rows, err := r.db.Query(
		"SELECT id, exercise_id, question_num, question, options, correct_index FROM listening_questions WHERE exercise_id = "+r.db.Placeholder(1)+" ORDER BY question_num",
		exerciseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ListeningQuestion
	for rows.Next() {
		var q models.ListeningQuestion
		if err := rows.Scan(&q.ID, &q.ExerciseID, &q.QuestionNum, &q.Question, &q.OptionsJSON, &q.CorrectIdx); err != nil {
			continue
		}
		_ = json.Unmarshal([]byte(q.OptionsJSON), &q.Options)
		list = append(list, q)
	}
	return list, rows.Err()
}

func (r *ListeningRepository) GetVocabulary(exerciseID string) ([]models.ListeningVocabulary, error) {
	rows, err := r.db.Query(
		"SELECT id, exercise_id, word, reading, meaning FROM listening_vocabulary WHERE exercise_id = "+r.db.Placeholder(1),
		exerciseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ListeningVocabulary
	for rows.Next() {
		var v models.ListeningVocabulary
		if err := rows.Scan(&v.ID, &v.ExerciseID, &v.Word, &v.Reading, &v.Meaning); err != nil {
			continue
		}
		list = append(list, v)
	}
	return list, rows.Err()
}

func (r *ListeningRepository) SetTTSCacheID(exerciseID, cacheID string) error {
	_, err := r.db.Exec(
		"UPDATE listening_exercises SET tts_cache_id = "+r.db.Placeholder(1)+" WHERE id = "+r.db.Placeholder(2),
		cacheID, exerciseID,
	)
	return err
}

func (r *ListeningRepository) CreateSession(userID, exerciseID string, totalQuestions int) (*models.ListeningSession, error) {
	id := r.db.GenerateUUID()
	now := time.Now()
	_, err := r.db.Exec(
		"INSERT INTO listening_sessions (id, user_id, exercise_id, started_at, score, total_questions, status) VALUES ("+
			r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", 0, "+r.db.Placeholder(5)+", 'in_progress')",
		id, userID, exerciseID, now, totalQuestions,
	)
	if err != nil {
		return nil, err
	}
	return &models.ListeningSession{
		ID: id, UserID: userID, ExerciseID: exerciseID,
		StartedAt: now, TotalQuestions: totalQuestions, Status: "in_progress",
	}, nil
}

func (r *ListeningRepository) GetSession(sessionID string) (*models.ListeningSession, error) {
	s := &models.ListeningSession{}
	err := r.db.QueryRow(
		"SELECT id, user_id, exercise_id, started_at, completed_at, score, total_questions, status FROM listening_sessions WHERE id = "+r.db.Placeholder(1),
		sessionID,
	).Scan(&s.ID, &s.UserID, &s.ExerciseID, &s.StartedAt, &s.CompletedAt, &s.Score, &s.TotalQuestions, &s.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *ListeningRepository) SubmitAnswer(sessionID, questionID string, answerIdx, correctIdx int) (bool, error) {
	isCorrect := answerIdx == correctIdx
	id := r.db.GenerateUUID()
	_, err := r.db.Exec(
		"INSERT OR IGNORE INTO listening_session_answers (id, session_id, question_id, answer_index, is_correct) VALUES ("+
			r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", "+r.db.Placeholder(5)+")",
		id, sessionID, questionID, answerIdx, isCorrect,
	)
	if err != nil {
		return false, err
	}
	if isCorrect {
		_, err = r.db.Exec(
			"UPDATE listening_sessions SET score = score + 1 WHERE id = "+r.db.Placeholder(1),
			sessionID,
		)
	}
	return isCorrect, err
}

func (r *ListeningRepository) GetQuestion(questionID string) (*models.ListeningQuestion, error) {
	q := &models.ListeningQuestion{}
	err := r.db.QueryRow(
		"SELECT id, exercise_id, question_num, question, options, correct_index FROM listening_questions WHERE id = "+r.db.Placeholder(1),
		questionID,
	).Scan(&q.ID, &q.ExerciseID, &q.QuestionNum, &q.Question, &q.OptionsJSON, &q.CorrectIdx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	_ = json.Unmarshal([]byte(q.OptionsJSON), &q.Options)
	return q, err
}
