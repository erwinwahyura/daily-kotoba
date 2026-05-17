package repository

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type KanjiRepository struct {
	db *db.DB
}

func NewKanjiRepository(db *db.DB) *KanjiRepository {
	return &KanjiRepository{db: db}
}

func (r *KanjiRepository) GetByCharacter(character string) (*models.Kanji, error) {
	var k models.Kanji
	var readingsJSON, strokeOrderJSON string

	err := r.db.QueryRow(
		"SELECT id, character, jlpt_level, meaning, readings, stroke_count, stroke_order, created_at FROM kanji_characters WHERE character = "+r.db.Placeholder(1),
		character,
	).Scan(&k.ID, &k.Character, &k.JLPTLevel, &k.Meaning, &readingsJSON, &k.StrokeCount, &strokeOrderJSON, &k.CreatedAt)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(readingsJSON), &k.Readings); err != nil {
		k.Readings = []string{}
	}
	if err := json.Unmarshal([]byte(strokeOrderJSON), &k.StrokeOrder); err != nil {
		k.StrokeOrder = []models.Stroke{}
	}
	return &k, nil
}

func (r *KanjiRepository) GetByLevel(level string) ([]models.Kanji, error) {
	rows, err := r.db.Query(
		"SELECT id, character, jlpt_level, meaning, readings, stroke_count, stroke_order, created_at FROM kanji_characters WHERE jlpt_level = "+r.db.Placeholder(1)+" ORDER BY character",
		level,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Kanji
	for rows.Next() {
		var k models.Kanji
		var readingsJSON, strokeOrderJSON string
		if err := rows.Scan(&k.ID, &k.Character, &k.JLPTLevel, &k.Meaning, &readingsJSON, &k.StrokeCount, &strokeOrderJSON, &k.CreatedAt); err != nil {
			continue
		}
		_ = json.Unmarshal([]byte(readingsJSON), &k.Readings)
		_ = json.Unmarshal([]byte(strokeOrderJSON), &k.StrokeOrder)
		list = append(list, k)
	}
	return list, rows.Err()
}

func (r *KanjiRepository) CreateSession(userID, kanjiID, kanjiChar string) (*models.KanjiPracticeSession, error) {
	id := r.db.GenerateUUID()
	now := time.Now()
	_, err := r.db.Exec(
		"INSERT INTO kanji_practice_sessions (id, user_id, kanji_id, kanji_char, started_at, current_stroke, accuracy, status) VALUES ("+
			r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", "+r.db.Placeholder(5)+", 0, 0.0, 'in_progress')",
		id, userID, kanjiID, kanjiChar, now,
	)
	if err != nil {
		return nil, err
	}
	return &models.KanjiPracticeSession{
		ID: id, UserID: userID, KanjiID: kanjiID, KanjiChar: kanjiChar,
		StartedAt: now, Accuracy: 0, Status: "in_progress",
	}, nil
}

func (r *KanjiRepository) GetSession(sessionID string) (*models.KanjiPracticeSession, error) {
	s := &models.KanjiPracticeSession{}
	err := r.db.QueryRow(
		"SELECT id, user_id, kanji_id, kanji_char, started_at, completed_at, current_stroke, accuracy, status FROM kanji_practice_sessions WHERE id = "+r.db.Placeholder(1),
		sessionID,
	).Scan(&s.ID, &s.UserID, &s.KanjiID, &s.KanjiChar, &s.StartedAt, &s.CompletedAt, &s.CurrentStroke, &s.Accuracy, &s.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return s, err
}

func (r *KanjiRepository) UpdateSessionAccuracy(sessionID string, accuracy float64, strokeNum int) error {
	_, err := r.db.Exec(
		"UPDATE kanji_practice_sessions SET accuracy = "+r.db.Placeholder(1)+", current_stroke = "+r.db.Placeholder(2)+" WHERE id = "+r.db.Placeholder(3),
		accuracy, strokeNum, sessionID,
	)
	return err
}
