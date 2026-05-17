package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

// ConversationRepository handles conversation session and message operations
type ConversationRepository struct {
	db *db.DB
}

// NewConversationRepository creates a new conversation repository
func NewConversationRepository(db *db.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

// CreateSession creates a new conversation session
func (r *ConversationRepository) CreateSession(session *models.ConversationSession) error {
	query := `
		INSERT INTO conversation_sessions (id, user_id, mode, scenario_id, level, status, started_at, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query,
		session.ID,
		session.UserID,
		session.Mode,
		session.ScenarioID,
		session.Level,
		session.Status,
		session.StartedAt,
		session.CreatedAt,
		session.UpdatedAt,
	)
	return err
}

// GetSession retrieves a conversation session by ID
func (r *ConversationRepository) GetSession(id string) (*models.ConversationSession, error) {
	query := `
		SELECT id, user_id, mode, scenario_id, level, status, naturalness_avg, 
		       started_at, ended_at, created_at, updated_at
		FROM conversation_sessions
		WHERE id = ?
	`
	var session models.ConversationSession
	var endedAt sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&session.ID,
		&session.UserID,
		&session.Mode,
		&session.ScenarioID,
		&session.Level,
		&session.Status,
		&session.NaturalnessAvg,
		&session.StartedAt,
		&endedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	return &session, nil
}

// GetActiveSessionByUser gets the active session for a user
func (r *ConversationRepository) GetActiveSessionByUser(userID string) (*models.ConversationSession, error) {
	query := `
		SELECT id, user_id, mode, scenario_id, level, status, naturalness_avg, 
		       started_at, ended_at, created_at, updated_at
		FROM conversation_sessions
		WHERE user_id = ? AND status = 'active'
		ORDER BY started_at DESC
		LIMIT 1
	`
	var session models.ConversationSession
	var endedAt sql.NullTime
	err := r.db.QueryRow(query, userID).Scan(
		&session.ID,
		&session.UserID,
		&session.Mode,
		&session.ScenarioID,
		&session.Level,
		&session.Status,
		&session.NaturalnessAvg,
		&session.StartedAt,
		&endedAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if endedAt.Valid {
		session.EndedAt = &endedAt.Time
	}
	return &session, nil
}

// UpdateSessionStatus updates the status of a session
func (r *ConversationRepository) UpdateSessionStatus(id string, status string, naturalnessAvg *int) error {
	query := `UPDATE conversation_sessions SET status = ?`
	args := []interface{}{status}
	
	if status == "completed" || status == "abandoned" {
		query += `, ended_at = ?`
		args = append(args, time.Now())
	}
	
	if naturalnessAvg != nil {
		query += `, naturalness_avg = ?`
		args = append(args, *naturalnessAvg)
	}
	
	query += `, updated_at = ? WHERE id = ?`
	args = append(args, time.Now(), id)
	
	_, err := r.db.Exec(query, args...)
	return err
}

// CreateMessage creates a new conversation message
func (r *ConversationRepository) CreateMessage(msg *models.ConversationMessage) error {
	query := `
		INSERT INTO conversation_messages (id, session_id, sender, content, correction, 
		                                   naturalness_score, alternatives, metadata, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	alternativesJSON, _ := json.Marshal(msg.Alternatives)
	metadataJSON, _ := json.Marshal(msg.Metadata)
	
	_, err := r.db.Exec(query,
		msg.ID,
		msg.SessionID,
		msg.Sender,
		msg.Content,
		msg.Correction,
		msg.NaturalnessScore,
		alternativesJSON,
		metadataJSON,
		msg.CreatedAt,
	)
	return err
}

// GetMessagesBySession retrieves all messages for a session
func (r *ConversationRepository) GetMessagesBySession(sessionID string) ([]models.ConversationMessage, error) {
	query := `
		SELECT id, session_id, sender, content, correction, naturalness_score, 
		       alternatives, metadata, created_at
		FROM conversation_messages
		WHERE session_id = ?
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var messages []models.ConversationMessage
	for rows.Next() {
		var msg models.ConversationMessage
		var alternativesJSON, metadataJSON []byte
		err := rows.Scan(
			&msg.ID,
			&msg.SessionID,
			&msg.Sender,
			&msg.Content,
			&msg.Correction,
			&msg.NaturalnessScore,
			&alternativesJSON,
			&metadataJSON,
			&msg.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if len(alternativesJSON) > 0 {
			json.Unmarshal(alternativesJSON, &msg.Alternatives)
		}
		if len(metadataJSON) > 0 {
			json.Unmarshal(metadataJSON, &msg.Metadata)
		}
		
		messages = append(messages, msg)
	}
	
	return messages, rows.Err()
}

// GetUserSessions retrieves all sessions for a user
func (r *ConversationRepository) GetUserSessions(userID string, limit int) ([]models.ConversationSession, error) {
	query := `
		SELECT id, user_id, mode, scenario_id, level, status, naturalness_avg, 
		       started_at, ended_at, created_at, updated_at
		FROM conversation_sessions
		WHERE user_id = ?
		ORDER BY started_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sessions []models.ConversationSession
	for rows.Next() {
		var session models.ConversationSession
		var endedAt sql.NullTime
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Mode,
			&session.ScenarioID,
			&session.Level,
			&session.Status,
			&session.NaturalnessAvg,
			&session.StartedAt,
			&endedAt,
			&session.CreatedAt,
			&session.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if endedAt.Valid {
			session.EndedAt = &endedAt.Time
		}
		sessions = append(sessions, session)
	}
	
	return sessions, rows.Err()
}

// CreateScenario creates a new conversation scenario
func (r *ConversationRepository) CreateScenario(scenario *models.ConversationScenario) error {
	query := `
		INSERT INTO conversation_scenarios (id, name, description, level, category, icon, 
		                                   system_prompt, starter_messages, vocabulary_hints, 
		                                   grammar_points, is_active, sort_order, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	
	starterJSON, _ := json.Marshal(scenario.StarterMessages)
	vocabJSON, _ := json.Marshal(scenario.VocabularyHints)
	grammarJSON, _ := json.Marshal(scenario.GrammarPoints)
	
	isActive := 0
	if scenario.IsActive {
		isActive = 1
	}
	
	_, err := r.db.Exec(query,
		scenario.ID,
		scenario.Name,
		scenario.Description,
		scenario.Level,
		scenario.Category,
		scenario.Icon,
		scenario.SystemPrompt,
		starterJSON,
		vocabJSON,
		grammarJSON,
		isActive,
		scenario.SortOrder,
		scenario.CreatedAt,
	)
	return err
}

// GetScenario retrieves a scenario by ID
func (r *ConversationRepository) GetScenario(id string) (*models.ConversationScenario, error) {
	query := `
		SELECT id, name, description, level, category, icon, system_prompt, 
		       starter_messages, vocabulary_hints, grammar_points, is_active, sort_order, created_at
		FROM conversation_scenarios
		WHERE id = ?
	`
	var scenario models.ConversationScenario
	var starterJSON, vocabJSON, grammarJSON []byte
	var isActive int
	
	err := r.db.QueryRow(query, id).Scan(
		&scenario.ID,
		&scenario.Name,
		&scenario.Description,
		&scenario.Level,
		&scenario.Category,
		&scenario.Icon,
		&scenario.SystemPrompt,
		&starterJSON,
		&vocabJSON,
		&grammarJSON,
		&isActive,
		&scenario.SortOrder,
		&scenario.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	scenario.IsActive = isActive == 1
	
	if len(starterJSON) > 0 {
		json.Unmarshal(starterJSON, &scenario.StarterMessages)
	}
	if len(vocabJSON) > 0 {
		json.Unmarshal(vocabJSON, &scenario.VocabularyHints)
	}
	if len(grammarJSON) > 0 {
		json.Unmarshal(grammarJSON, &scenario.GrammarPoints)
	}
	
	return &scenario, nil
}

// GetScenariosByLevel retrieves active scenarios for a level
func (r *ConversationRepository) GetScenariosByLevel(level int, category string) ([]models.ConversationScenario, error) {
	query := `
		SELECT id, name, description, level, category, icon, system_prompt, 
		       starter_messages, vocabulary_hints, grammar_points, is_active, sort_order, created_at
		FROM conversation_scenarios
		WHERE is_active = 1 AND level <= ?
	`
	args := []interface{}{level}
	
	if category != "" {
		query += ` AND category = ?`
		args = append(args, category)
	}
	
	query += ` ORDER BY sort_order ASC, level ASC`
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var scenarios []models.ConversationScenario
	for rows.Next() {
		var scenario models.ConversationScenario
		var starterJSON, vocabJSON, grammarJSON []byte
		var isActive int
		
		err := rows.Scan(
			&scenario.ID,
			&scenario.Name,
			&scenario.Description,
			&scenario.Level,
			&scenario.Category,
			&scenario.Icon,
			&scenario.SystemPrompt,
			&starterJSON,
			&vocabJSON,
			&grammarJSON,
			&isActive,
			&scenario.SortOrder,
			&scenario.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		scenario.IsActive = isActive == 1
		
		if len(starterJSON) > 0 {
			json.Unmarshal(starterJSON, &scenario.StarterMessages)
		}
		if len(vocabJSON) > 0 {
			json.Unmarshal(vocabJSON, &scenario.VocabularyHints)
		}
		if len(grammarJSON) > 0 {
			json.Unmarshal(grammarJSON, &scenario.GrammarPoints)
		}
		
		scenarios = append(scenarios, scenario)
	}
	
	return scenarios, rows.Err()
}

// GetAllScenarios retrieves all scenarios (for admin/seeding)
func (r *ConversationRepository) GetAllScenarios() ([]models.ConversationScenario, error) {
	query := `
		SELECT id, name, description, level, category, icon, system_prompt, 
		       starter_messages, vocabulary_hints, grammar_points, is_active, sort_order, created_at
		FROM conversation_scenarios
		ORDER BY sort_order ASC
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var scenarios []models.ConversationScenario
	for rows.Next() {
		var scenario models.ConversationScenario
		var starterJSON, vocabJSON, grammarJSON []byte
		var isActive int
		
		err := rows.Scan(
			&scenario.ID,
			&scenario.Name,
			&scenario.Description,
			&scenario.Level,
			&scenario.Category,
			&scenario.Icon,
			&scenario.SystemPrompt,
			&starterJSON,
			&vocabJSON,
			&grammarJSON,
			&isActive,
			&scenario.SortOrder,
			&scenario.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		scenario.IsActive = isActive == 1
		
		if len(starterJSON) > 0 {
			json.Unmarshal(starterJSON, &scenario.StarterMessages)
		}
		if len(vocabJSON) > 0 {
			json.Unmarshal(vocabJSON, &scenario.VocabularyHints)
		}
		if len(grammarJSON) > 0 {
			json.Unmarshal(grammarJSON, &scenario.GrammarPoints)
		}
		
		scenarios = append(scenarios, scenario)
	}
	
	return scenarios, rows.Err()
}

// SeedScenarios inserts default scenarios if none exist
func (r *ConversationRepository) SeedScenarios() error {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM conversation_scenarios").Scan(&count)
	if err != nil {
		return err
	}
	
	if count > 0 {
		return nil // Already seeded
	}
	
	scenarios := []models.ConversationScenario{
		{
			ID:          "konbini",
			Name:        "コンビニ (Konbini)",
			Description: "Practice small talk at a convenience store",
			Level:       5,
			Category:    "daily",
			Icon:        "🏪",
			SystemPrompt: `You are a friendly convenience store clerk in Japan. Respond naturally to customer questions and comments. Keep responses short and casual. Use appropriate level of politeness (です/ます for N5, more casual for higher levels). Provide gentle corrections if the user makes mistakes.`,
			StarterMessages: models.JSONB{
				"messages": []string{
					"いらっしゃいませ！",
					"お弁当、温めますか？",
					"ポイントカードはお持ちですか？",
				},
			},
			VocabularyHints: models.JSONB{
				"words": []string{"温める", "ポイントカード", "袋", "レジ"},
			},
			GrammarPoints: models.JSONB{
				"points": []string{"〜ますか (polite question)", "お〜 (honorific prefix)"},
			},
			IsActive:    true,
			SortOrder:   1,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "restaurant",
			Name:        "レストラン (Restaurant)",
			Description: "Order food and chat at a restaurant",
			Level:       5,
			Category:    "daily",
			Icon:        "🍜",
			SystemPrompt: `You are a waiter/waitress at a casual Japanese restaurant. Help the customer order, answer questions about the menu, and make small talk. Be friendly and helpful.`,
			StarterMessages: models.JSONB{
				"messages": []string{
					"いらっしゃいませ。何名様ですか？",
					"メニューはこちらになります。",
					"おすすめはラーメンですよ。",
				},
			},
			VocabularyHints: models.JSONB{
				"words": []string{"おすすめ", "メニュー", "お会計", "お持ち帰り"},
			},
			GrammarPoints: models.JSONB{
				"points": []string{"〜名様 (counter for people)", "お〜になる (honorific)"},
			},
			IsActive:    true,
			SortOrder:   2,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "directions",
			Name:        "道を聞く (Asking Directions)",
			Description: "Ask for and give directions",
			Level:       4,
			Category:    "travel",
			Icon:        "🗺️",
			SystemPrompt: `You are a local passerby in Tokyo. Someone asks you for directions. Give clear, helpful directions using landmarks and simple Japanese. Be patient and friendly.`,
			StarterMessages: models.JSONB{
				"messages": []string{
					"すみません、駅はどこですか？",
					"この辺り、詳しいですか？",
					"あの建物の隣ですよ。",
				},
			},
			VocabularyHints: models.JSONB{
				"words": []string{"交差点", "信号", "曲がる", "まっすぐ"},
			},
			GrammarPoints: models.JSONB{
				"points": []string{"〜はどこですか", "〜を曲がる", "〜の隣"},
			},
			IsActive:    true,
			SortOrder:   3,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "weather_chat",
			Name:        "天気の話 (Weather Chat)",
			Description: "Casual small talk about the weather",
			Level:       5,
			Category:    "social",
			Icon:        "☀️",
			SystemPrompt: `You are a friendly acquaintance who starts casual conversations about the weather. Keep it light and natural. Use filler words and casual expressions appropriate for the user's level.`,
			StarterMessages: models.JSONB{
				"messages": []string{
					"あ、雨降ってきたね。",
					"今日、暑いですね。",
					"えーと、傘持ってきた？",
				},
			},
			VocabularyHints: models.JSONB{
				"words": []string{"暑い", "寒い", "雨", "傘", "天気"},
			},
			GrammarPoints: models.JSONB{
				"points": []string{"〜ね (sentence ending)", "〜ですね", "えーと (filler)"},
			},
			IsActive:    true,
			SortOrder:   4,
			CreatedAt:   time.Now(),
		},
		{
			ID:          "weekend_plans",
			Name:        "週末の予定 (Weekend Plans)",
			Description: "Chat about what you're doing this weekend",
			Level:       4,
			Category:    "social",
			Icon:        "🎉",
			SystemPrompt: `You are a coworker or friend asking about weekend plans. Be casual and interested. Share your own plans too. Use natural conversational Japanese.`,
			StarterMessages: models.JSONB{
				"messages": []string{
					"週末、何か予定ある？",
					"今週末、暇？",
					"なんか、いいことあった？",
				},
			},
			VocabularyHints: models.JSONB{
				"words": []string{"週末", "予定", "暇", "遊ぶ", "行く"},
			},
			GrammarPoints: models.JSONB{
				"points": []string{"〜ある？ (casual question)", "なんか (casual 'something')"},
			},
			IsActive:    true,
			SortOrder:   5,
			CreatedAt:   time.Now(),
		},
	}
	
	for _, scenario := range scenarios {
		if err := r.CreateScenario(&scenario); err != nil {
			return fmt.Errorf("failed to seed scenario %s: %w", scenario.ID, err)
		}
	}
	
	return nil
}

// GetSessionStats gets conversation statistics for a user
func (r *ConversationRepository) GetSessionStats(userID string) (map[string]interface{}, error) {
	stats := make(map[string]interface{})
	
	// Total sessions
	var totalSessions int
	err := r.db.QueryRow(
		"SELECT COUNT(*) FROM conversation_sessions WHERE user_id = ?",
		userID,
	).Scan(&totalSessions)
	if err != nil {
		return nil, err
	}
	stats["total_sessions"] = totalSessions
	
	// Completed sessions
	var completedSessions int
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM conversation_sessions WHERE user_id = ? AND status = 'completed'",
		userID,
	).Scan(&completedSessions)
	if err != nil {
		return nil, err
	}
	stats["completed_sessions"] = completedSessions
	
	// Average naturalness score
	var avgNaturalness sql.NullInt64
	err = r.db.QueryRow(
		"SELECT AVG(naturalness_avg) FROM conversation_sessions WHERE user_id = ? AND naturalness_avg IS NOT NULL",
		userID,
	).Scan(&avgNaturalness)
	if err != nil {
		return nil, err
	}
	if avgNaturalness.Valid {
		stats["average_naturalness"] = avgNaturalness.Int64
	} else {
		stats["average_naturalness"] = 0
	}
	
	// Sessions by mode
	rows, err := r.db.Query(
		"SELECT mode, COUNT(*) FROM conversation_sessions WHERE user_id = ? GROUP BY mode",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	modeCounts := make(map[string]int)
	for rows.Next() {
		var mode string
		var count int
		if err := rows.Scan(&mode, &count); err != nil {
			return nil, err
		}
		modeCounts[mode] = count
	}
	stats["sessions_by_mode"] = modeCounts
	
	return stats, nil
}
