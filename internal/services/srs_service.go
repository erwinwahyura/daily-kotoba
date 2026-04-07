package services

import (
	"fmt"
	"time"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type SRSService struct {
	srsRepo     *repository.SRSRepository
	vocabRepo   *repository.VocabRepository
	grammarRepo *repository.GrammarRepository
	userRepo    *repository.UserRepository
}

func NewSRSService(
	srsRepo *repository.SRSRepository,
	vocabRepo *repository.VocabRepository,
	grammarRepo *repository.GrammarRepository,
	userRepo *repository.UserRepository,
) *SRSService {
	return &SRSService{
		srsRepo:     srsRepo,
		vocabRepo:   vocabRepo,
		grammarRepo: grammarRepo,
		userRepo:    userRepo,
	}
}

// SubmitReview processes a review and updates SRS schedule
func (s *SRSService) SubmitReview(userID string, req *models.SRSReviewRequest) (*models.SRSReviewResponse, error) {
	// Get or create schedule
	schedule, err := s.srsRepo.GetOrCreateSchedule(userID, req.ItemID, req.ItemType)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule: %w", err)
	}

	// Calculate SM-2 result
	result := models.CalculateSM2(
		req.Quality,
		schedule.IntervalDays,
		schedule.Repetitions,
		schedule.EaseFactor,
	)

	// Record history
	history := &models.SRSReviewHistory{
		ScheduleID:       schedule.ID,
		UserID:           userID,
		Quality:          req.Quality,
		ResponseTimeMs:   req.ResponseTimeMs,
		ItemType:         req.ItemType,
		ItemID:           req.ItemID,
		IntervalBefore:   schedule.IntervalDays,
		IntervalAfter:    result.Interval,
		EaseFactorBefore: schedule.EaseFactor,
		EaseFactorAfter:  result.EaseFactor,
	}

	if err := s.srsRepo.RecordReviewHistory(history); err != nil {
		return nil, fmt.Errorf("failed to record history: %w", err)
	}

	// Update schedule with result
	oldStatus := schedule.Status
	if err := s.srsRepo.UpdateSchedule(schedule, result); err != nil {
		return nil, fmt.Errorf("failed to update schedule: %w", err)
	}

	// Update study session
	correct := 0
	if result.Interval > 1 {
		correct = 1
	}
	reviewType := 0
	if schedule.Repetitions == 0 {
		reviewType = 1 // New item
	}
	if err := s.srsRepo.UpdateStudySession(userID, reviewType, 1, 0, correct, 1); err != nil {
		// Non-fatal, just log
		_ = err
	}

	// Build response
	response := &models.SRSReviewResponse{
		Schedule:      schedule,
		StatusChanged: oldStatus != result.Status,
	}

	// Human-readable next review
	switch result.Interval {
	case 1:
		response.NextReview = "Tomorrow"
	case 6:
		response.NextReview = "In 6 days"
	default:
		response.NextReview = fmt.Sprintf("In %d days", result.Interval)
	}

	// Check for achievements
	achievement := s.checkAchievements(userID, schedule, result)
	if achievement != nil {
		response.NewAchievement = achievement
	}

	return response, nil
}

// GetReviewQueue returns items due for review
func (s *SRSService) GetReviewQueue(userID string, limit int) (*models.SRSQueueResponse, error) {
	schedules, err := s.srsRepo.GetDueItems(userID, limit)
	if err != nil {
		return nil, err
	}

	response := &models.SRSQueueResponse{
		DueItems: make([]models.SRSDueItem, 0, len(schedules)),
	}

	for _, sched := range schedules {
		item := models.SRSDueItem{
			ID:          sched.ItemID,
			Type:        sched.ItemType,
			Schedule:    sched,
			DaysOverdue: int(time.Since(sched.NextReviewAt).Hours() / 24),
		}

		// Load actual item data
		if sched.ItemType == "vocabulary" {
			vocab, err := s.vocabRepo.GetByID(sched.ItemID)
			if err != nil {
				continue // Skip if item not found
			}
			item.Data = vocab
		} else if sched.ItemType == "grammar" {
			pattern, err := s.grammarRepo.GetByID(sched.ItemID)
			if err != nil {
				continue
			}
			item.Data = pattern
		}

		response.DueItems = append(response.DueItems, item)
	}

	// Get counts
	dueToday, dueTomorrow, _, learning, review, mastered, err := s.srsRepo.CountDueItems(userID)
	if err != nil {
		return nil, err
	}

	response.TotalDue = len(response.DueItems)
	response.LearningItems = learning
	response.ReviewItems = review
	response.MasteredItems = mastered
	response.NewItems = dueToday + dueTomorrow

	return response, nil
}

// GetStats returns SRS statistics for user
func (s *SRSService) GetStats(userID string) (*models.SRSStats, error) {
	return s.srsRepo.GetSRSStats(userID)
}

// InitializeItem creates an SRS schedule for a newly learned item
func (s *SRSService) InitializeItem(userID, itemID, itemType string) error {
	_, err := s.srsRepo.GetOrCreateSchedule(userID, itemID, itemType)
	return err
}

// checkAchievements checks if user unlocked any achievements
func (s *SRSService) checkAchievements(userID string, schedule *models.SRSSchedule, result *models.SM2Result) *models.Achievement {
	// Check for streak achievement
	if schedule.Streak == 7 {
		return &models.Achievement{
			ID:          "streak_7",
			Name:        "7-Day Streak",
			Description: "Reviewed 7 items correctly in a row",
			Icon:        "🔥",
			UnlockedAt:  time.Now().Format("2006-01-02"),
		}
	}

	// Check for mastered achievement
	if result.Status == "mastered" && schedule.Status != "mastered" {
		return &models.Achievement{
			ID:          "first_mastered",
			Name:        "First Mastery",
			Description: "Mastered your first item",
			Icon:        "⭐",
			UnlockedAt:  time.Now().Format("2006-01-02"),
		}
	}

	return nil
}