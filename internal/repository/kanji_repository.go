package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

// KanjiRepository handles kanji data access
type KanjiRepository struct {
	db *db.DB
}

// NewKanjiRepository creates a new repository
func NewKanjiRepository(db *db.DB) *KanjiRepository {
	return &KanjiRepository{db: db}
}

// GetKanjiByCharacter retrieves kanji by character
func (r *KanjiRepository) GetKanjiByCharacter(char string) (*models.Kanji, error) {
	kanji := &models.Kanji{}
	var readingsJSON, strokeOrderJSON []byte

	query := `
		SELECT id, character, jlpt_level, meaning, readings, stroke_count, stroke_order, created_at
		FROM kanji WHERE character = $1
	`

	err := r.db.QueryRow(query, char).Scan(
		&kanji.ID,
		&kanji.Character,
		&kanji.JLPTLevel,
		&kanji.Meaning,
		&readingsJSON,
		&kanji.StrokeCount,
		&strokeOrderJSON,
		&kanji.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("kanji not found: %s", char)
		}
		return nil, err
	}

	// Parse JSON arrays
	if err := json.Unmarshal(readingsJSON, &kanji.Readings); err != nil {
		return nil, fmt.Errorf("failed to parse readings: %w", err)
	}
	if err := json.Unmarshal(strokeOrderJSON, &kanji.StrokeOrder); err != nil {
		return nil, fmt.Errorf("failed to parse stroke order: %w", err)
	}

	return kanji, nil
}

// GetKanjiByLevel retrieves kanji by JLPT level
func (r *KanjiRepository) GetKanjiByLevel(level string, limit int) ([]models.Kanji, error) {
	if limit <= 0 {
		limit = 50
	}

	query := `
		SELECT id, character, jlpt_level, meaning, readings, stroke_count, stroke_order, created_at
		FROM kanji WHERE jlpt_level = $1 ORDER BY stroke_count ASC, character ASC LIMIT $2
	`

	rows, err := r.db.Query(query, level, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kanjiList []models.Kanji
	for rows.Next() {
		kanji := models.Kanji{}
		var readingsJSON, strokeOrderJSON []byte

		err := rows.Scan(
			&kanji.ID,
			&kanji.Character,
			&kanji.JLPTLevel,
			&kanji.Meaning,
			&readingsJSON,
			&kanji.StrokeCount,
			&strokeOrderJSON,
			&kanji.CreatedAt,
		)
		if err != nil {
			continue
		}

		// Parse JSON
		json.Unmarshal(readingsJSON, &kanji.Readings)
		json.Unmarshal(strokeOrderJSON, &kanji.StrokeOrder)

		kanjiList = append(kanjiList, kanji)
	}

	return kanjiList, rows.Err()
}

// CreatePracticeSession creates a new practice session
func (r *KanjiRepository) CreatePracticeSession(session *models.KanjiPracticeSession) error {
	userStrokesJSON, _ := json.Marshal(session.UserStrokes)

	query := `
		INSERT INTO kanji_practice_sessions (id, user_id, kanji_id, kanji_char, started_at, status, user_strokes, accuracy)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		session.ID,
		session.UserID,
		session.KanjiID,
		session.KanjiChar,
		session.StartedAt,
		session.Status,
		userStrokesJSON,
		session.Accuracy,
	)

	return err
}

// UpdatePracticeSession updates session with new strokes
func (r *KanjiRepository) UpdatePracticeSession(session *models.KanjiPracticeSession) error {
	userStrokesJSON, _ := json.Marshal(session.UserStrokes)

	var completedAt interface{}
	if session.CompletedAt != nil {
		completedAt = *session.CompletedAt
	} else {
		completedAt = nil
	}

	query := `
		UPDATE kanji_practice_sessions 
		SET user_strokes = $1, accuracy = $2, status = $3, completed_at = $4
		WHERE id = $5
	`

	_, err := r.db.Exec(query,
		userStrokesJSON,
		session.Accuracy,
		session.Status,
		completedAt,
		session.ID,
	)

	return err
}

// GetPracticeSession retrieves a practice session
func (r *KanjiRepository) GetPracticeSession(sessionID string) (*models.KanjiPracticeSession, error) {
	session := &models.KanjiPracticeSession{}
	var userStrokesJSON []byte
	var completedAt sql.NullTime

	query := `
		SELECT id, user_id, kanji_id, kanji_char, started_at, completed_at, user_strokes, accuracy, status
		FROM kanji_practice_sessions WHERE id = $1
	`

	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID,
		&session.UserID,
		&session.KanjiID,
		&session.KanjiChar,
		&session.StartedAt,
		&completedAt,
		&userStrokesJSON,
		&session.Accuracy,
		&session.Status,
	)

	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}

	json.Unmarshal(userStrokesJSON, &session.UserStrokes)

	return session, nil
}

// GetUserKanjiStats gets user's kanji practice stats
func (r *KanjiRepository) GetUserKanjiStats(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Total practiced
	var total int
	r.db.QueryRow("SELECT COUNT(*) FROM kanji_practice_sessions WHERE user_id = $1", userID).Scan(&total)
	stats["total_practiced"] = total

	// Completed sessions
	var completed int
	r.db.QueryRow("SELECT COUNT(*) FROM kanji_practice_sessions WHERE user_id = $1 AND status = 'completed'", userID).Scan(&completed)
	stats["completed"] = completed

	// Average accuracy
	var avgAccuracy float64
	r.db.QueryRow("SELECT COALESCE(AVG(accuracy), 0) FROM kanji_practice_sessions WHERE user_id = $1 AND status = 'completed'", userID).Scan(&avgAccuracy)
	stats["average_accuracy"] = avgAccuracy

	return stats, nil
}

// SeedSampleKanji seeds initial kanji data
func (r *KanjiRepository) SeedSampleKanji() error {
	sampleKanji := []models.Kanji{
		{
			ID:          "kanji_001",
			Character:   "日",
			JLPTLevel:   "N5",
			Meaning:     "Sun, day",
			Readings:    []string{"にち", "ひ", "か"},
			StrokeCount: 4,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 50, Y: 80}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 35}, EndPoint: models.Point{X: 80, Y: 35}},
				{StrokeNum: 3, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 65}, EndPoint: models.Point{X: 80, Y: 65}},
				{StrokeNum: 4, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 50}, EndPoint: models.Point{X: 80, Y: 50}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_002",
			Character:   "月",
			JLPTLevel:   "N5",
			Meaning:     "Moon, month",
			Readings:    []string{"げつ", "つき"},
			StrokeCount: 4,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "vertical", StartPoint: models.Point{X: 30, Y: 20}, EndPoint: models.Point{X: 30, Y: 80}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 30, Y: 35}, EndPoint: models.Point{X: 70, Y: 35}},
				{StrokeNum: 3, Direction: "horizontal", StartPoint: models.Point{X: 30, Y: 65}, EndPoint: models.Point{X: 70, Y: 65}},
				{StrokeNum: 4, Direction: "curve", StartPoint: models.Point{X: 70, Y: 35}, EndPoint: models.Point{X: 70, Y: 65}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_003",
			Character:   "火",
			JLPTLevel:   "N5",
			Meaning:     "Fire",
			Readings:    []string{"か", "ひ"},
			StrokeCount: 4,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "diagonal", StartPoint: models.Point{X: 35, Y: 25}, EndPoint: models.Point{X: 25, Y: 45}},
				{StrokeNum: 2, Direction: "diagonal", StartPoint: models.Point{X: 65, Y: 25}, EndPoint: models.Point{X: 75, Y: 45}},
				{StrokeNum: 3, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 35}, EndPoint: models.Point{X: 50, Y: 75}},
				{StrokeNum: 4, Direction: "curve", StartPoint: models.Point{X: 30, Y: 60}, EndPoint: models.Point{X: 70, Y: 60}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_004",
			Character:   "水",
			JLPTLevel:   "N5",
			Meaning:     "Water",
			Readings:    []string{"すい", "みず"},
			StrokeCount: 4,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 35, Y: 45}},
				{StrokeNum: 2, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 65, Y: 45}},
				{StrokeNum: 3, Direction: "horizontal", StartPoint: models.Point{X: 30, Y: 50}, EndPoint: models.Point{X: 70, Y: 50}},
				{StrokeNum: 4, Direction: "curve", StartPoint: models.Point{X: 50, Y: 50}, EndPoint: models.Point{X: 50, Y: 80}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_005",
			Character:   "木",
			JLPTLevel:   "N5",
			Meaning:     "Tree, wood",
			Readings:    []string{"もく", "き"},
			StrokeCount: 4,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 50, Y: 80}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 50}, EndPoint: models.Point{X: 80, Y: 50}},
				{StrokeNum: 3, Direction: "diagonal", StartPoint: models.Point{X: 30, Y: 35}, EndPoint: models.Point{X: 15, Y: 55}},
				{StrokeNum: 4, Direction: "diagonal", StartPoint: models.Point{X: 70, Y: 35}, EndPoint: models.Point{X: 85, Y: 55}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_006",
			Character:   "金",
			JLPTLevel:   "N5",
			Meaning:     "Gold, metal, money",
			Readings:    []string{"きん", "かね"},
			StrokeCount: 8,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "horizontal", StartPoint: models.Point{X: 25, Y: 20}, EndPoint: models.Point{X: 40, Y: 20}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 60, Y: 20}, EndPoint: models.Point{X: 75, Y: 20}},
				{StrokeNum: 3, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 35}, EndPoint: models.Point{X: 80, Y: 35}},
				{StrokeNum: 4, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 50, Y: 80}},
				{StrokeNum: 5, Direction: "horizontal", StartPoint: models.Point{X: 30, Y: 50}, EndPoint: models.Point{X: 70, Y: 50}},
				{StrokeNum: 6, Direction: "diagonal", StartPoint: models.Point{X: 35, Y: 50}, EndPoint: models.Point{X: 25, Y: 70}},
				{StrokeNum: 7, Direction: "diagonal", StartPoint: models.Point{X: 65, Y: 50}, EndPoint: models.Point{X: 75, Y: 70}},
				{StrokeNum: 8, Direction: "horizontal", StartPoint: models.Point{X: 25, Y: 80}, EndPoint: models.Point{X: 75, Y: 80}},
			},
			CreatedAt: time.Now(),
		},
		// Additional kanji for practice
		{
			ID:          "kanji_007",
			Character:   "人",
			JLPTLevel:   "N5",
			Meaning:     "Person, people",
			Readings:    []string{"じん", "にん", "ひと"},
			StrokeCount: 2,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 25, Y: 70}},
				{StrokeNum: 2, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 75, Y: 70}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_008",
			Character:   "大",
			JLPTLevel:   "N5",
			Meaning:     "Big, large",
			Readings:    []string{"だい", "おお"},
			StrokeCount: 3,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 35}, EndPoint: models.Point{X: 80, Y: 35}},
				{StrokeNum: 2, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 35}, EndPoint: models.Point{X: 25, Y: 80}},
				{StrokeNum: 3, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 35}, EndPoint: models.Point{X: 75, Y: 80}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_009",
			Character:   "小",
			JLPTLevel:   "N5",
			Meaning:     "Small",
			Readings:    []string{"しょう", "ちい", "こ", "お"},
			StrokeCount: 3,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 20}, EndPoint: models.Point{X: 50, Y: 55}},
				{StrokeNum: 2, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 45}, EndPoint: models.Point{X: 25, Y: 75}},
				{StrokeNum: 3, Direction: "diagonal", StartPoint: models.Point{X: 50, Y: 45}, EndPoint: models.Point{X: 75, Y: 75}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_010",
			Character:   "上",
			JLPTLevel:   "N5",
			Meaning:     "Up, above",
			Readings:    []string{"じょう", "うえ", "あ", "のぼ", "かみ"},
			StrokeCount: 3,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 30}, EndPoint: models.Point{X: 80, Y: 30}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 35, Y: 55}, EndPoint: models.Point{X: 65, Y: 55}},
				{StrokeNum: 3, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 55}, EndPoint: models.Point{X: 50, Y: 85}},
			},
			CreatedAt: time.Now(),
		},
		{
			ID:          "kanji_011",
			Character:   "下",
			JLPTLevel:   "N5",
			Meaning:     "Down, below",
			Readings:    []string{"か", "げ", "くだ", "お", "しも", "さ"},
			StrokeCount: 3,
			StrokeOrder: []models.Stroke{
				{StrokeNum: 1, Direction: "horizontal", StartPoint: models.Point{X: 20, Y: 30}, EndPoint: models.Point{X: 80, Y: 30}},
				{StrokeNum: 2, Direction: "horizontal", StartPoint: models.Point{X: 35, Y: 55}, EndPoint: models.Point{X: 65, Y: 55}},
				{StrokeNum: 3, Direction: "vertical", StartPoint: models.Point{X: 50, Y: 30}, EndPoint: models.Point{X: 50, Y: 85}},
			},
			CreatedAt: time.Now(),
		},
	}

	for _, k := range sampleKanji {
		readingsJSON, _ := json.Marshal(k.Readings)
		strokeJSON, _ := json.Marshal(k.StrokeOrder)

		// Check if exists
		var exists bool
		r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM kanji WHERE character = $1)", k.Character).Scan(&exists)
		if exists {
			continue
		}

		query := `
			INSERT INTO kanji (id, character, jlpt_level, meaning, readings, stroke_count, stroke_order, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err := r.db.Exec(query, k.ID, k.Character, k.JLPTLevel, k.Meaning, readingsJSON, k.StrokeCount, strokeJSON, k.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to seed kanji %s: %w", k.Character, err)
		}
	}

	return nil
}
