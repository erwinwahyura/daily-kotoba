package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

// GoalsRepository handles goals and achievements data access
type GoalsRepository struct {
	db *db.DB
}

// NewGoalsRepository creates a new repository
func NewGoalsRepository(db *db.DB) *GoalsRepository {
	return &GoalsRepository{db: db}
}

// GetGoalSettings retrieves user's goal settings
func (r *GoalsRepository) GetGoalSettings(userID string) (*models.GoalSettings, error) {
	settings := &models.GoalSettings{}
	
	query := `
		SELECT id, user_id, vocab_target, grammar_target, kanji_target, 
		       conjugation_target, reading_target, enable_reminders, reminder_time,
		       created_at, updated_at
		FROM goal_settings WHERE user_id = $1
	`
	
	err := r.db.QueryRow(query, userID).Scan(
		&settings.ID, &settings.UserID, &settings.VocabTarget, &settings.GrammarTarget,
		&settings.KanjiTarget, &settings.ConjugationTarget, &settings.ReadingTarget,
		&settings.EnableReminders, &settings.ReminderTime, &settings.CreatedAt, &settings.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		// Create default settings
		return r.CreateDefaultGoalSettings(userID)
	}
	
	return settings, err
}

// CreateDefaultGoalSettings creates default settings for a user
func (r *GoalsRepository) CreateDefaultGoalSettings(userID string) (*models.GoalSettings, error) {
	settings := &models.GoalSettings{
		ID:                "settings_" + userID,
		UserID:            userID,
		VocabTarget:       10,
		GrammarTarget:     5,
		KanjiTarget:       5,
		ConjugationTarget: 20,
		ReadingTarget:     1,
		EnableReminders:   false,
		ReminderTime:      "20:00",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	
	query := `
		INSERT INTO goal_settings (id, user_id, vocab_target, grammar_target, kanji_target,
			conjugation_target, reading_target, enable_reminders, reminder_time, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.Exec(query, settings.ID, settings.UserID, settings.VocabTarget, settings.GrammarTarget,
		settings.KanjiTarget, settings.ConjugationTarget, settings.ReadingTarget,
		settings.EnableReminders, settings.ReminderTime, settings.CreatedAt, settings.UpdatedAt)
	
	return settings, err
}

// UpdateGoalSettings updates user's goal settings
func (r *GoalsRepository) UpdateGoalSettings(settings *models.GoalSettings) error {
	query := `
		UPDATE goal_settings SET
			vocab_target = $1, grammar_target = $2, kanji_target = $3,
			conjugation_target = $4, reading_target = $5,
			enable_reminders = $6, reminder_time = $7, updated_at = $8
		WHERE user_id = $9
	`
	
	_, err := r.db.Exec(query, settings.VocabTarget, settings.GrammarTarget, settings.KanjiTarget,
		settings.ConjugationTarget, settings.ReadingTarget, settings.EnableReminders,
		settings.ReminderTime, time.Now(), settings.UserID)
	
	return err
}

// GetOrCreateDailyGoal gets or creates today's goal for a user
func (r *GoalsRepository) GetOrCreateDailyGoal(userID string, date time.Time) (*models.DailyGoal, error) {
	// Try to get existing
	goal, err := r.GetDailyGoal(userID, date)
	if err == nil {
		return goal, nil
	}
	
	// Create new goal with settings defaults
	settings, err := r.GetGoalSettings(userID)
	if err != nil {
		return nil, err
	}
	
	goal = &models.DailyGoal{
		ID:                "goal_" + userID + "_" + date.Format("2006-01-02"),
		UserID:            userID,
		Date:              date,
		VocabTarget:       settings.VocabTarget,
		GrammarTarget:     settings.GrammarTarget,
		KanjiTarget:       settings.KanjiTarget,
		ConjugationTarget: settings.ConjugationTarget,
		ReadingTarget:     settings.ReadingTarget,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	
	query := `
		INSERT INTO daily_goals (id, user_id, date, vocab_target, grammar_target, kanji_target,
			conjugation_target, reading_target, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err = r.db.Exec(query, goal.ID, goal.UserID, goal.Date, goal.VocabTarget, goal.GrammarTarget,
		goal.KanjiTarget, goal.ConjugationTarget, goal.ReadingTarget, goal.CreatedAt, goal.UpdatedAt)
	
	return goal, err
}

// GetDailyGoal retrieves a specific daily goal
func (r *GoalsRepository) GetDailyGoal(userID string, date time.Time) (*models.DailyGoal, error) {
	goal := &models.DailyGoal{}
	var completedAt sql.NullTime
	
	query := `
		SELECT id, user_id, date, vocab_target, vocab_completed, grammar_target, grammar_completed,
			kanji_target, kanji_completed, conjugation_target, conjugation_completed,
			reading_target, reading_completed, is_completed, completed_at, created_at, updated_at
		FROM daily_goals WHERE user_id = $1 AND date = $2
	`
	
	err := r.db.QueryRow(query, userID, date).Scan(
		&goal.ID, &goal.UserID, &goal.Date, &goal.VocabTarget, &goal.VocabCompleted,
		&goal.GrammarTarget, &goal.GrammarCompleted, &goal.KanjiTarget, &goal.KanjiCompleted,
		&goal.ConjugationTarget, &goal.ConjugationCompleted, &goal.ReadingTarget, &goal.ReadingCompleted,
		&goal.IsCompleted, &completedAt, &goal.CreatedAt, &goal.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("goal not found")
	}
	
	if completedAt.Valid {
		goal.CompletedAt = &completedAt.Time
	}
	
	return goal, err
}

// UpdateDailyGoal updates progress
func (r *GoalsRepository) UpdateDailyGoal(goal *models.DailyGoal) error {
	// Check if completed
	if !goal.IsCompleted && r.isGoalComplete(goal) {
		goal.IsCompleted = true
		now := time.Now()
		goal.CompletedAt = &now
	}
	
	query := `
		UPDATE daily_goals SET
			vocab_completed = $1, grammar_completed = $2, kanji_completed = $3,
			conjugation_completed = $4, reading_completed = $5,
			is_completed = $6, completed_at = $7, updated_at = $8
		WHERE id = $9
	`
	
	_, err := r.db.Exec(query, goal.VocabCompleted, goal.GrammarCompleted, goal.KanjiCompleted,
		goal.ConjugationCompleted, goal.ReadingCompleted, goal.IsCompleted, goal.CompletedAt,
		time.Now(), goal.ID)
	
	return err
}

// isGoalComplete checks if all targets are met
func (r *GoalsRepository) isGoalComplete(goal *models.DailyGoal) bool {
	return goal.VocabCompleted >= goal.VocabTarget &&
		goal.GrammarCompleted >= goal.GrammarTarget &&
		goal.KanjiCompleted >= goal.KanjiTarget &&
		goal.ConjugationCompleted >= goal.ConjugationTarget &&
		goal.ReadingCompleted >= goal.ReadingTarget
}

// GetUserStreak retrieves user's streak info
func (r *GoalsRepository) GetUserStreak(userID string) (*models.UserStreak, error) {
	streak := &models.UserStreak{}
	var lastActivity sql.NullTime
	
	query := `
		SELECT id, user_id, current_streak, longest_streak, last_activity_date, total_active_days, created_at, updated_at
		FROM user_streaks WHERE user_id = $1
	`
	
	err := r.db.QueryRow(query, userID).Scan(
		&streak.ID, &streak.UserID, &streak.CurrentStreak, &streak.LongestStreak,
		&lastActivity, &streak.TotalActiveDays, &streak.CreatedAt, &streak.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		// Create new streak
		return r.CreateUserStreak(userID)
	}
	
	if lastActivity.Valid {
		streak.LastActivityDate = &lastActivity.Time
	}
	
	return streak, err
}

// CreateUserStreak creates new streak record
func (r *GoalsRepository) CreateUserStreak(userID string) (*models.UserStreak, error) {
	streak := &models.UserStreak{
		ID:              "streak_" + userID,
		UserID:          userID,
		CurrentStreak:   0,
		LongestStreak:   0,
		TotalActiveDays: 0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	
	query := `
		INSERT INTO user_streaks (id, user_id, current_streak, longest_streak, total_active_days, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Exec(query, streak.ID, streak.UserID, streak.CurrentStreak,
		streak.LongestStreak, streak.TotalActiveDays, streak.CreatedAt, streak.UpdatedAt)
	
	return streak, err
}

// UpdateStreak updates user's streak
func (r *GoalsRepository) UpdateStreak(streak *models.UserStreak) error {
	query := `
		UPDATE user_streaks SET
			current_streak = $1, longest_streak = $2, last_activity_date = $3,
			total_active_days = $4, updated_at = $5
		WHERE user_id = $6
	`
	
	_, err := r.db.Exec(query, streak.CurrentStreak, streak.LongestStreak, streak.LastActivityDate,
		streak.TotalActiveDays, time.Now(), streak.UserID)
	
	return err
}

// RecordActivity records activity and updates streak
func (r *GoalsRepository) RecordActivity(userID string) error {
	streak, err := r.GetUserStreak(userID)
	if err != nil {
		return err
	}
	
	today := time.Now().Truncate(24 * time.Hour)
	
	// Check if already active today
	if streak.LastActivityDate != nil {
		lastActivity := streak.LastActivityDate.Truncate(24 * time.Hour)
		
		if lastActivity.Equal(today) {
			// Already active today, no change
			return nil
		}
		
		// Check if yesterday
		yesterday := today.Add(-24 * time.Hour)
		if lastActivity.Equal(yesterday) {
			// Continue streak
			streak.CurrentStreak++
			if streak.CurrentStreak > streak.LongestStreak {
				streak.LongestStreak = streak.CurrentStreak
			}
		} else {
			// Streak broken
			streak.CurrentStreak = 1
		}
	} else {
		// First activity
		streak.CurrentStreak = 1
		streak.LongestStreak = 1
	}
	
	streak.LastActivityDate = &today
	streak.TotalActiveDays++
	
	return r.UpdateStreak(streak)
}

// GetUserAchievements retrieves user's earned achievements
func (r *GoalsRepository) GetUserAchievements(userID string) ([]models.Achievement, error) {
	query := `
		SELECT id, user_id, achievement_id, type, name, description, icon, level, earned_at
		FROM achievements WHERE user_id = $1 ORDER BY earned_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var achievements []models.Achievement
	for rows.Next() {
		a := models.Achievement{}
		err := rows.Scan(&a.ID, &a.UserID, &a.Type, &a.Name, &a.Description, &a.Icon, &a.Level, &a.EarnedAt)
		if err != nil {
			continue
		}
		achievements = append(achievements, a)
	}
	
	return achievements, rows.Err()
}

// AddAchievement adds a new achievement
func (r *GoalsRepository) AddAchievement(userID string, def models.AchievementDefinition) error {
	// Check if already earned
	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM achievements WHERE user_id = $1 AND achievement_id = $2)",
		userID, def.ID).Scan(&exists)
	if err != nil || exists {
		return err // Already earned or error
	}
	
	query := `
		INSERT INTO achievements (id, user_id, achievement_id, type, name, description, icon, level, earned_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	
	_, err = r.db.Exec(query,
		"ach_"+userID+"_"+def.ID,
		userID,
		def.ID,
		def.Type,
		def.Name,
		def.Description,
		def.Icon,
		def.Level,
		time.Now(),
	)
	
	return err
}

// GetRecentDailyGoals gets goals for last N days
func (r *GoalsRepository) GetRecentDailyGoals(userID string, days int) ([]models.DailyGoal, error) {
	query := `
		SELECT id, user_id, date, vocab_target, vocab_completed, grammar_target, grammar_completed,
			kanji_target, kanji_completed, conjugation_target, conjugation_completed,
			reading_target, reading_completed, is_completed, completed_at, created_at, updated_at
		FROM daily_goals 
		WHERE user_id = $1 AND date >= date('now', '-$2 days')
		ORDER BY date DESC
	`
	
	rows, err := r.db.Query(query, userID, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var goals []models.DailyGoal
	for rows.Next() {
		g := models.DailyGoal{}
		var completedAt sql.NullTime
		
		err := rows.Scan(
			&g.ID, &g.UserID, &g.Date, &g.VocabTarget, &g.VocabCompleted,
			&g.GrammarTarget, &g.GrammarCompleted, &g.KanjiTarget, &g.KanjiCompleted,
			&g.ConjugationTarget, &g.ConjugationCompleted, &g.ReadingTarget, &g.ReadingCompleted,
			&g.IsCompleted, &completedAt, &g.CreatedAt, &g.UpdatedAt,
		)
		if err != nil {
			continue
		}
		
		if completedAt.Valid {
			g.CompletedAt = &completedAt.Time
		}
		
		goals = append(goals, g)
	}
	
	return goals, rows.Err()
}
