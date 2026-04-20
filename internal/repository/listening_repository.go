package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

// ListeningRepository handles listening exercise data access
type ListeningRepository struct {
	db *db.DB
}

// NewListeningRepository creates a new repository
func NewListeningRepository(db *db.DB) *ListeningRepository {
	return &ListeningRepository{db: db}
}

// GetExercisesByLevel retrieves listening exercises by JLPT level
func (r *ListeningRepository) GetExercisesByLevel(level string, limit int) ([]models.ListeningExercise, error) {
	if limit <= 0 || limit > 50 {
		limit = 20
	}

	query := `
		SELECT id, title, jlpt_level, difficulty, audio_url, duration, transcript, 
		       translation, vocabulary, questions, topic, created_at
		FROM listening_exercises 
		WHERE jlpt_level = $1 
		ORDER BY difficulty, created_at DESC 
		LIMIT $2
	`

	rows, err := r.db.Query(query, level, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanExercises(rows)
}

// GetExerciseByID retrieves a specific exercise
func (r *ListeningRepository) GetExerciseByID(id string) (*models.ListeningExercise, error) {
	exercise := &models.ListeningExercise{}
	var vocabJSON, questionsJSON []byte

	query := `
		SELECT id, title, jlpt_level, difficulty, audio_url, duration, transcript,
		       translation, vocabulary, questions, topic, created_at
		FROM listening_exercises WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&exercise.ID, &exercise.Title, &exercise.JLPTLevel, &exercise.Difficulty,
		&exercise.AudioURL, &exercise.Duration, &exercise.Transcript,
		&exercise.Translation, &vocabJSON, &questionsJSON, &exercise.Topic, &exercise.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("exercise not found")
	}
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err := json.Unmarshal(vocabJSON, &exercise.Vocabulary); err != nil {
		return nil, fmt.Errorf("failed to parse vocabulary: %w", err)
	}
	if err := json.Unmarshal(questionsJSON, &exercise.Questions); err != nil {
		return nil, fmt.Errorf("failed to parse questions: %w", err)
	}

	return exercise, nil
}

// CreateSession creates a new listening session
func (r *ListeningRepository) CreateSession(session *models.ListeningSession) error {
	answersJSON, _ := json.Marshal(session.Answers)

	query := `
		INSERT INTO listening_sessions (id, user_id, exercise_id, started_at, current_position, 
			answers, status, play_count, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query, session.ID, session.UserID, session.ExerciseID,
		session.StartedAt, session.CurrentPosition, answersJSON, session.Status,
		session.PlayCount, time.Now())

	return err
}

// UpdateSession updates a listening session
func (r *ListeningRepository) UpdateSession(session *models.ListeningSession) error {
	answersJSON, _ := json.Marshal(session.Answers)

	var completedAt interface{}
	if session.CompletedAt != nil {
		completedAt = *session.CompletedAt
	} else {
		completedAt = nil
	}

	query := `
		UPDATE listening_sessions 
		SET current_position = $1, answers = $2, score = $3, status = $4, 
		    completed_at = $5, play_count = $6
		WHERE id = $7
	`

	_, err := r.db.Exec(query, session.CurrentPosition, answersJSON, session.Score,
		session.Status, completedAt, session.PlayCount, session.ID)

	return err
}

// GetSession retrieves a listening session
func (r *ListeningRepository) GetSession(sessionID string) (*models.ListeningSession, error) {
	session := &models.ListeningSession{}
	var answersJSON []byte
	var completedAt sql.NullTime

	query := `
		SELECT id, user_id, exercise_id, started_at, completed_at, current_position,
		       answers, score, status, play_count
		FROM listening_sessions WHERE id = $1
	`

	err := r.db.QueryRow(query, sessionID).Scan(
		&session.ID, &session.UserID, &session.ExerciseID, &session.StartedAt,
		&completedAt, &session.CurrentPosition, &answersJSON, &session.Score,
		&session.Status, &session.PlayCount,
	)
	if err != nil {
		return nil, err
	}

	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}

	json.Unmarshal(answersJSON, &session.Answers)

	return session, nil
}

// GetUserStats retrieves user's listening statistics
func (r *ListeningRepository) GetUserStats(userID string) (*models.ListeningProgress, error) {
	stats := &models.ListeningProgress{
		ByLevel: make(map[string]models.LevelStats),
	}

	// Total completed and average score
	query := `
		SELECT COUNT(*), COALESCE(AVG(score), 0)
		FROM listening_sessions 
		WHERE user_id = $1 AND status = 'completed'
	`
	err := r.db.QueryRow(query, userID).Scan(&stats.CompletedCount, &stats.AverageScore)
	if err != nil {
		return nil, err
	}

	// Total exercises available
	err = r.db.QueryRow("SELECT COUNT(*) FROM listening_exercises").Scan(&stats.TotalExercises)
	if err != nil {
		return nil, err
	}

	// Stats by level
	levelQuery := `
		SELECT e.jlpt_level, COUNT(*), AVG(s.score)
		FROM listening_sessions s
		JOIN listening_exercises e ON s.exercise_id = e.id
		WHERE s.user_id = $1 AND s.status = 'completed'
		GROUP BY e.jlpt_level
	`
	rows, err := r.db.Query(levelQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var level string
		var count int
		var avg float64
		if err := rows.Scan(&level, &count, &avg); err != nil {
			continue
		}
		stats.ByLevel[level] = models.LevelStats{
			Completed: count,
			Average:   avg,
		}
	}

	return stats, rows.Err()
}

// SeedSampleExercises adds sample listening exercises
func (r *ListeningRepository) SeedSampleExercises() error {
	exercises := []models.ListeningExercise{
		{
			ID:          "listen_n5_001",
			Title:       "At the Restaurant",
			JLPTLevel:   "N5",
			Difficulty:  "easy",
			AudioURL:    "https://example.com/audio/n5_restaurant.mp3",
			Duration:    45,
			Transcript:  "いらっしゃいませ。何名様ですか。\n二人です。\nこちらへどうぞ。メニューをどうぞ。\nありがとうございます。",
			Translation: "Welcome. How many people?\nTwo people.\nThis way please. Here's the menu.\nThank you.",
			Vocabulary: []models.VocabItem{
				{Word: "いらっしゃいませ", Reading: "いらっしゃいませ", Meaning: "Welcome", Timestamp: 0},
				{Word: "何名様", Reading: "なんめいさま", Meaning: "How many people", Timestamp: 2},
				{Word: "二人", Reading: "ふたり", Meaning: "Two people", Timestamp: 5},
				{Word: "メニュー", Reading: "めにゅー", Meaning: "Menu", Timestamp: 12},
			},
			Questions: []models.ListeningQuestion{
				{
					ID:          "q1",
					Question:    "How many people are there?",
					Options:     []string{"One", "Two", "Three", "Four"},
					Correct:     1,
					Timestamp:   45,
					Explanation: "The customer says 二人 (ふたり) which means two people.",
				},
			},
			Topic: "conversation",
		},
		{
			ID:          "listen_n4_001",
			Title:       "Train Station Announcement",
			JLPTLevel:   "N4",
			Difficulty:  "medium",
			AudioURL:    "https://example.com/audio/n4_station.mp3",
			Duration:    60,
			Transcript:  "まもなく、東京駅に到着いたします。\nお出口は左側です。\nお忘れ物のないようにご注意ください。\nありがとうございました。",
			Translation: "We will soon arrive at Tokyo Station.\nThe exit is on the left side.\nPlease be careful not to forget your belongings.\nThank you.",
			Vocabulary: []models.VocabItem{
				{Word: "まもなく", Reading: "まもなく", Meaning: "Soon", Timestamp: 0},
				{Word: "到着", Reading: "とうちゃく", Meaning: "Arrival", Timestamp: 3},
				{Word: "お出口", Reading: "おでぐち", Meaning: "Exit", Timestamp: 10},
				{Word: "お忘れ物", Reading: "おわすれもの", Meaning: "Forgotten items", Timestamp: 20},
			},
			Questions: []models.ListeningQuestion{
				{
					ID:          "q1",
					Question:    "Which side is the exit?",
					Options:     []string{"Right side", "Left side", "Both sides", "Not mentioned"},
					Correct:     1,
					Timestamp:   60,
					Explanation: "The announcement says 左側 (ひだりがわ) which means left side.",
				},
			},
			Topic: "announcement",
		},
	}

	for _, ex := range exercises {
		// Check if exists
		var exists bool
		r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM listening_exercises WHERE id = $1)", ex.ID).Scan(&exists)
		if exists {
			continue
		}

		vocabJSON, _ := json.Marshal(ex.Vocabulary)
		questionsJSON, _ := json.Marshal(ex.Questions)

		query := `
			INSERT INTO listening_exercises (id, title, jlpt_level, difficulty, audio_url, duration,
				transcript, translation, vocabulary, questions, topic, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		`
		_, err := r.db.Exec(query, ex.ID, ex.Title, ex.JLPTLevel, ex.Difficulty,
			ex.AudioURL, ex.Duration, ex.Transcript, ex.Translation, vocabJSON, questionsJSON,
			ex.Topic, time.Now())
		if err != nil {
			return fmt.Errorf("failed to seed exercise %s: %w", ex.ID, err)
		}
	}

	return nil
}

// scanExercises helper to scan exercise rows
func (r *ListeningRepository) scanExercises(rows *sql.Rows) ([]models.ListeningExercise, error) {
	var exercises []models.ListeningExercise

	for rows.Next() {
		ex := models.ListeningExercise{}
		var vocabJSON, questionsJSON []byte

		err := rows.Scan(
			&ex.ID, &ex.Title, &ex.JLPTLevel, &ex.Difficulty, &ex.AudioURL,
			&ex.Duration, &ex.Transcript, &ex.Translation, &vocabJSON, &questionsJSON,
			&ex.Topic, &ex.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal(vocabJSON, &ex.Vocabulary)
		json.Unmarshal(questionsJSON, &ex.Questions)

		exercises = append(exercises, ex)
	}

	return exercises, rows.Err()
}
