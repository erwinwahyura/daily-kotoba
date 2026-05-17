package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type GoalsRepository struct {
	db *db.DB
}

func NewGoalsRepository(db *db.DB) *GoalsRepository {
	return &GoalsRepository{db: db}
}

func (r *GoalsRepository) GetOrCreateSettings(userID string) (*models.UserGoalSettings, error) {
	s := &models.UserGoalSettings{}
	err := r.db.QueryRow(
		"SELECT user_id, vocab_target, grammar_target, kanji_target, conjugation_target, reading_target, updated_at FROM user_goal_settings WHERE user_id = "+r.db.Placeholder(1),
		userID,
	).Scan(&s.UserID, &s.VocabTarget, &s.GrammarTarget, &s.KanjiTarget, &s.ConjugationTarget, &s.ReadingTarget, &s.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		// create defaults
		s = &models.UserGoalSettings{
			UserID: userID, VocabTarget: 10, GrammarTarget: 5,
			KanjiTarget: 5, ConjugationTarget: 20, ReadingTarget: 1,
			UpdatedAt: time.Now(),
		}
		_, insertErr := r.db.Exec(
			"INSERT INTO user_goal_settings (user_id, vocab_target, grammar_target, kanji_target, conjugation_target, reading_target) VALUES ("+r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", "+r.db.Placeholder(5)+", "+r.db.Placeholder(6)+")",
			s.UserID, s.VocabTarget, s.GrammarTarget, s.KanjiTarget, s.ConjugationTarget, s.ReadingTarget,
		)
		return s, insertErr
	}
	return s, err
}

func (r *GoalsRepository) SaveSettings(s *models.UserGoalSettings) error {
	_, err := r.db.Exec(
		"INSERT OR REPLACE INTO user_goal_settings (user_id, vocab_target, grammar_target, kanji_target, conjugation_target, reading_target, updated_at) VALUES ("+r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", "+r.db.Placeholder(5)+", "+r.db.Placeholder(6)+", CURRENT_TIMESTAMP)",
		s.UserID, s.VocabTarget, s.GrammarTarget, s.KanjiTarget, s.ConjugationTarget, s.ReadingTarget,
	)
	return err
}

func (r *GoalsRepository) GetDailyCount(userID, date, activityType string) (int, error) {
	var count int
	err := r.db.QueryRow(
		"SELECT COALESCE(count, 0) FROM user_daily_activities WHERE user_id = "+r.db.Placeholder(1)+" AND activity_date = "+r.db.Placeholder(2)+" AND activity_type = "+r.db.Placeholder(3),
		userID, date, activityType,
	).Scan(&count)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	return count, err
}

func (r *GoalsRepository) IncrementActivity(userID, date, activityType string, count int) error {
	id := r.db.GenerateUUID()
	_, err := r.db.Exec(
		"INSERT INTO user_daily_activities (id, user_id, activity_date, activity_type, count) VALUES ("+
			r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+", "+r.db.Placeholder(4)+", "+r.db.Placeholder(5)+") "+
			"ON CONFLICT(user_id, activity_date, activity_type) DO UPDATE SET count = count + "+r.db.Placeholder(6)+", updated_at = CURRENT_TIMESTAMP",
		id, userID, date, activityType, count, count,
	)
	return err
}

func (r *GoalsRepository) GetWeeklyActivity(userID string, days []string) (map[string]map[string]int, error) {
	result := make(map[string]map[string]int)
	for _, d := range days {
		result[d] = make(map[string]int)
	}

	rows, err := r.db.Query(
		"SELECT activity_date, activity_type, count FROM user_daily_activities WHERE user_id = "+r.db.Placeholder(1)+" AND activity_date >= "+r.db.Placeholder(2)+" ORDER BY activity_date",
		userID, days[0],
	)
	if err != nil {
		return result, err
	}
	defer rows.Close()

	for rows.Next() {
		var date, actType string
		var count int
		if err := rows.Scan(&date, &actType, &count); err != nil {
			continue
		}
		if _, ok := result[date]; ok {
			result[date][actType] = count
		}
	}
	return result, rows.Err()
}

func (r *GoalsRepository) GetAllAchievements() ([]models.AchievementDef, error) {
	rows, err := r.db.Query("SELECT id, name, description, icon, category, requirement_type, requirement_value, created_at FROM achievements ORDER BY requirement_value ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.AchievementDef
	for rows.Next() {
		var a models.AchievementDef
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.Icon, &a.Category, &a.RequirementType, &a.RequirementValue, &a.CreatedAt); err != nil {
			continue
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *GoalsRepository) GetUserAchievements(userID string) ([]models.UserAchievement, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, achievement_id, earned_at FROM user_achievements WHERE user_id = "+r.db.Placeholder(1)+" ORDER BY earned_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.UserAchievement
	for rows.Next() {
		var a models.UserAchievement
		if err := rows.Scan(&a.ID, &a.UserID, &a.AchievementID, &a.EarnedAt); err != nil {
			continue
		}
		list = append(list, a)
	}
	return list, rows.Err()
}

func (r *GoalsRepository) GrantAchievement(userID, achievementID string) error {
	id := r.db.GenerateUUID()
	_, err := r.db.Exec(
		"INSERT OR IGNORE INTO user_achievements (id, user_id, achievement_id) VALUES ("+r.db.Placeholder(1)+", "+r.db.Placeholder(2)+", "+r.db.Placeholder(3)+")",
		id, userID, achievementID,
	)
	return err
}
