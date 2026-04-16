package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

type SRSRepository struct {
	db *db.DB
}

func NewSRSRepository(db *db.DB) *SRSRepository {
	return &SRSRepository{db: db}
}

// GetOrCreateSchedule gets existing schedule or creates new one
func (r *SRSRepository) GetOrCreateSchedule(userID, itemID, itemType string) (*models.SRSSchedule, error) {
	// Try to get existing
	schedule, err := r.GetSchedule(userID, itemID, itemType)
	if err == nil {
		return schedule, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create new schedule
	schedule = &models.SRSSchedule{
		ID:           r.db.GenerateUUID(),
		UserID:       userID,
		ItemID:       itemID,
		ItemType:     itemType,
		IntervalDays: 0,
		Repetitions:  0,
		EaseFactor:   2.5,
		Status:       "learning",
		NextReviewAt: time.Now(),
	}

	query := fmt.Sprintf(`
		INSERT INTO srs_schedules 
		(id, user_id, item_id, item_type, interval_days, repetitions, ease_factor, next_review_at, status)
		VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s)
	`, r.db.Placeholder(1), r.db.Placeholder(2), r.db.Placeholder(3), 
	   r.db.Placeholder(4), r.db.Placeholder(5), r.db.Placeholder(6), 
	   r.db.Placeholder(7), r.db.Placeholder(8), r.db.Placeholder(9))

	_, err = r.db.Exec(query, schedule.ID, schedule.UserID, schedule.ItemID, 
		schedule.ItemType, schedule.IntervalDays, schedule.Repetitions, 
		schedule.EaseFactor, schedule.NextReviewAt, schedule.Status)
	
	if err != nil {
		return nil, fmt.Errorf("failed to create schedule: %w", err)
	}

	return schedule, nil
}

// GetSchedule retrieves a specific schedule
func (r *SRSRepository) GetSchedule(userID, itemID, itemType string) (*models.SRSSchedule, error) {
	schedule := &models.SRSSchedule{}
	query := `
		SELECT id, user_id, item_id, item_type, interval_days, repetitions, ease_factor,
		       last_reviewed_at, next_review_at, total_reviews, correct_reviews, streak, status
		FROM srs_schedules
		WHERE user_id = ` + r.db.Placeholder(1) + ` AND item_id = ` + r.db.Placeholder(2) + ` AND item_type = ` + r.db.Placeholder(3)

	err := r.db.QueryRow(query, userID, itemID, itemType).Scan(
		&schedule.ID, &schedule.UserID, &schedule.ItemID, &schedule.ItemType,
		&schedule.IntervalDays, &schedule.Repetitions, &schedule.EaseFactor,
		&schedule.LastReviewedAt, &schedule.NextReviewAt, &schedule.TotalReviews,
		&schedule.CorrectReviews, &schedule.Streak, &schedule.Status,
	)

	if err != nil {
		return nil, err
	}
	return schedule, nil
}

// GetDueItems retrieves items due for review
func (r *SRSRepository) GetDueItems(userID string, limit int) ([]*models.SRSSchedule, error) {
	if limit < 1 {
		limit = 20
	}

	query := `
		SELECT id, user_id, item_id, item_type, interval_days, repetitions, ease_factor,
		       last_reviewed_at, next_review_at, total_reviews, correct_reviews, streak, status
		FROM srs_schedules
		WHERE user_id = ` + r.db.Placeholder(1) + ` 
		  AND status IN ('learning', 'review')
		  AND next_review_at <= ` + r.db.Placeholder(2) + `
		ORDER BY next_review_at ASC
		LIMIT ` + r.db.Placeholder(3)

	rows, err := r.db.Query(query, userID, time.Now(), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*models.SRSSchedule
	for rows.Next() {
		s := &models.SRSSchedule{}
		err := rows.Scan(
			&s.ID, &s.UserID, &s.ItemID, &s.ItemType,
			&s.IntervalDays, &s.Repetitions, &s.EaseFactor,
			&s.LastReviewedAt, &s.NextReviewAt, &s.TotalReviews,
			&s.CorrectReviews, &s.Streak, &s.Status,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	return schedules, rows.Err()
}

// CountDueItems returns counts of items by status
func (r *SRSRepository) CountDueItems(userID string) (dueToday, dueTomorrow, newItems, learning, review, mastered int, err error) {
	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)

	// Due today
	query := `
		SELECT COUNT(*) FROM srs_schedules 
		WHERE user_id = ` + r.db.Placeholder(1) + ` 
		  AND status IN ('learning', 'review')
		  AND next_review_at <= ` + r.db.Placeholder(2)
	err = r.db.QueryRow(query, userID, now).Scan(&dueToday)
	if err != nil {
		return
	}

	// Due tomorrow
	err = r.db.QueryRow(query, userID, tomorrow).Scan(&dueTomorrow)
	if err != nil {
		return
	}

	// By status
	statusQuery := `
		SELECT status, COUNT(*) FROM srs_schedules 
		WHERE user_id = ` + r.db.Placeholder(1) + `
		GROUP BY status`
	
	rows, err := r.db.Query(statusQuery, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err = rows.Scan(&status, &count); err != nil {
			return
		}
		switch status {
		case "learning":
			learning = count
		case "review":
			review = count
		case "mastered":
			mastered = count
		}
	}

	newItems = learning + review + mastered // Already scheduled
	return
}

// UpdateSchedule updates schedule after a review
func (r *SRSRepository) UpdateSchedule(schedule *models.SRSSchedule, result *models.SM2Result) error {
	schedule.IntervalDays = result.Interval
	schedule.Repetitions = result.Repetitions
	schedule.EaseFactor = result.EaseFactor
	schedule.Status = result.Status
	schedule.LastReviewedAt = &[]time.Time{time.Now()}[0]
	schedule.NextReviewAt = time.Now().Add(time.Duration(result.Interval) * 24 * time.Hour)
	schedule.TotalReviews++
	if result.Interval > 1 { // Not a failure (failure resets to 1 day)
		schedule.CorrectReviews++
		schedule.Streak++
	} else {
		schedule.Streak = 0
	}

	query := fmt.Sprintf(`
		UPDATE srs_schedules
		SET interval_days = %s, repetitions = %s, ease_factor = %s,
		    last_reviewed_at = %s, next_review_at = %s, total_reviews = %s,
		    correct_reviews = %s, streak = %s, status = %s
		WHERE id = %s
	`, r.db.Placeholder(1), r.db.Placeholder(2), r.db.Placeholder(3),
	   r.db.Placeholder(4), r.db.Placeholder(5), r.db.Placeholder(6),
	   r.db.Placeholder(7), r.db.Placeholder(8), r.db.Placeholder(9),
	   r.db.Placeholder(10))

	_, err := r.db.Exec(query, schedule.IntervalDays, schedule.Repetitions,
		schedule.EaseFactor, schedule.LastReviewedAt, schedule.NextReviewAt,
		schedule.TotalReviews, schedule.CorrectReviews, schedule.Streak,
		schedule.Status, schedule.ID)

	return err
}

// RecordReviewHistory logs a review attempt
func (r *SRSRepository) RecordReviewHistory(history *models.SRSReviewHistory) error {
	history.ID = r.db.GenerateUUID()
	history.ReviewedAt = time.Now()

	query := fmt.Sprintf(`
		INSERT INTO srs_review_history 
		(id, schedule_id, user_id, quality, response_time_ms, item_type, item_id,
		 interval_before, interval_after, ease_factor_before, ease_factor_after, reviewed_at)
		VALUES (%s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s, %s)
	`, r.db.Placeholder(1), r.db.Placeholder(2), r.db.Placeholder(3),
	   r.db.Placeholder(4), r.db.Placeholder(5), r.db.Placeholder(6),
	   r.db.Placeholder(7), r.db.Placeholder(8), r.db.Placeholder(9),
	   r.db.Placeholder(10), r.db.Placeholder(11), r.db.Placeholder(12))

	_, err := r.db.Exec(query, history.ID, history.ScheduleID, history.UserID,
		history.Quality, history.ResponseTimeMs, history.ItemType, history.ItemID,
		history.IntervalBefore, history.IntervalAfter, history.EaseFactorBefore,
		history.EaseFactorAfter, history.ReviewedAt)

	return err
}

// GetStudySession gets or creates today's study session
func (r *SRSRepository) GetOrCreateStudySession(userID string) (*models.StudySession, error) {
	today := time.Now().Format("2006-01-02")
	
	session := &models.StudySession{}
	query := `
		SELECT id, user_id, session_date, new_items, review_items, drill_items,
		       total_time_seconds, correct_count, total_attempts
		FROM study_sessions
		WHERE user_id = ` + r.db.Placeholder(1) + ` AND session_date = ` + r.db.Placeholder(2)
	
	err := r.db.QueryRow(query, userID, today).Scan(
		&session.ID, &session.UserID, &session.SessionDate,
		&session.NewItems, &session.ReviewItems, &session.DrillItems,
		&session.TotalTimeSeconds, &session.CorrectCount, &session.TotalAttempts,
	)
	
	if err == nil {
		return session, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// Create new session
	session = &models.StudySession{
		ID:          r.db.GenerateUUID(),
		UserID:      userID,
		SessionDate: today,
	}

	insertQuery := fmt.Sprintf(`
		INSERT INTO study_sessions (id, user_id, session_date)
		VALUES (%s, %s, %s)
	`, r.db.Placeholder(1), r.db.Placeholder(2), r.db.Placeholder(3))

	_, err = r.db.Exec(insertQuery, session.ID, session.UserID, session.SessionDate)
	return session, err
}

// UpdateStudySession updates activity counts
func (r *SRSRepository) UpdateStudySession(userID string, newItems, reviewItems, drillItems int, correct, total int) error {
	today := time.Now().Format("2006-01-02")
	
	query := fmt.Sprintf(`
		UPDATE study_sessions
		SET new_items = new_items + %s,
		    review_items = review_items + %s,
		    drill_items = drill_items + %s,
		    correct_count = correct_count + %s,
		    total_attempts = total_attempts + %s
		WHERE user_id = %s AND session_date = %s
	`, r.db.Placeholder(1), r.db.Placeholder(2), r.db.Placeholder(3),
	   r.db.Placeholder(4), r.db.Placeholder(5), r.db.Placeholder(6), r.db.Placeholder(7))

	_, err := r.db.Exec(query, newItems, reviewItems, drillItems, correct, total, userID, today)
	return err
}

// GetSRSStats returns comprehensive SRS statistics
func (r *SRSRepository) GetSRSStats(userID string) (*models.SRSStats, error) {
	stats := &models.SRSStats{}
	
	// Count by status
	query := `
		SELECT status, COUNT(*) FROM srs_schedules
		WHERE user_id = ` + r.db.Placeholder(1) + `
		GROUP BY status`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	for rows.Next() {
		var status string
		var count int
		if err = rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch status {
		case "learning":
			stats.Learning = count
		case "review":
			stats.Review = count
		case "mastered":
			stats.Mastered = count
		case "lapsed":
			stats.Lapsed = count
		}
	}
	
	stats.TotalItems = stats.Learning + stats.Review + stats.Mastered + stats.Lapsed
	
	// Due today
	now := time.Now()
	dueQuery := `
		SELECT COUNT(*) FROM srs_schedules
		WHERE user_id = ` + r.db.Placeholder(1) + `
		  AND status IN ('learning', 'review')
		  AND next_review_at <= ` + r.db.Placeholder(2)
	err = r.db.QueryRow(dueQuery, userID, now).Scan(&stats.DueToday)
	if err != nil {
		return nil, err
	}
	
	// Due tomorrow
	tomorrow := now.Add(24 * time.Hour)
	err = r.db.QueryRow(dueQuery, userID, tomorrow).Scan(&stats.DueTomorrow)
	if err != nil {
		return nil, err
	}
	
	// Overall accuracy from review history
	accuracyQuery := `
		SELECT COUNT(*), COUNT(CASE WHEN quality >= 3 THEN 1 END)
		FROM srs_review_history
		WHERE user_id = ` + r.db.Placeholder(1)
	
	var total, correct int
	err = r.db.QueryRow(accuracyQuery, userID).Scan(&total, &correct)
	if err != nil {
		return nil, err
	}
	
	stats.TotalReviews = total
	if total > 0 {
		stats.Accuracy = float64(correct) / float64(total)
	}
	
	return stats, nil
}