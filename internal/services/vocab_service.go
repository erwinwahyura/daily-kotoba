package services

import (
	"fmt"
	"math"
	"time"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type VocabService struct {
	vocabRepo    *repository.VocabRepository
	progressRepo *repository.ProgressRepository
	userRepo     *repository.UserRepository
}

func NewVocabService(
	vocabRepo *repository.VocabRepository,
	progressRepo *repository.ProgressRepository,
	userRepo *repository.UserRepository,
) *VocabService {
	return &VocabService{
		vocabRepo:    vocabRepo,
		progressRepo: progressRepo,
		userRepo:     userRepo,
	}
}

func (s *VocabService) GetDailyWord(userID string) (*models.VocabularyWithProgress, error) {
	// Get user to know their current level
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get user's progress
	progress, err := s.progressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Update streak if needed
	if err := s.updateStreak(userID, progress); err != nil {
		return nil, err
	}

	// Get vocabulary at current index for user's level
	vocab, err := s.vocabRepo.GetByLevelAndIndex(user.CurrentLevel, progress.CurrentVocabIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get vocabulary: %w", err)
	}

	// Get total words count for the level
	totalWords, err := s.vocabRepo.GetTotalCountByLevel(user.CurrentLevel)
	if err != nil {
		return nil, err
	}

	// Create response with progress
	response := &models.VocabularyWithProgress{
		Vocabulary: vocab,
		Progress: &models.VocabularyProgress{
			CurrentIndex:      progress.CurrentVocabIndex,
			TotalWordsInLevel: totalWords,
			WordsLearned:      progress.WordsLearnedCount,
			StreakDays:        progress.StreakDays,
		},
	}

	return response, nil
}

func (s *VocabService) GetVocabByID(vocabID string) (*models.Vocabulary, error) {
	return s.vocabRepo.GetByID(vocabID)
}

func (s *VocabService) SkipToNextWord(userID, vocabID, status string) (*models.VocabularyWithProgress, error) {
	// Get user to know their current level
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Mark current word with status
	if err := s.progressRepo.MarkVocabStatus(userID, vocabID, status); err != nil {
		return nil, err
	}

	// Increment appropriate counter
	if status == "known" || status == "skipped" {
		if err := s.progressRepo.IncrementWordsSkipped(userID); err != nil {
			return nil, err
		}
	}

	// Increment vocab index
	progress, err := s.progressRepo.IncrementVocabIndex(userID)
	if err != nil {
		return nil, err
	}

	// Get total words for the level
	totalWords, err := s.vocabRepo.GetTotalCountByLevel(user.CurrentLevel)
	if err != nil {
		return nil, err
	}

	// Check if we've reached the end of the level
	// If yes, cycle back to 0 (Phase 1 logic)
	if progress.CurrentVocabIndex >= totalWords {
		progress.CurrentVocabIndex = 0
		if err := s.progressRepo.Update(progress); err != nil {
			return nil, err
		}
	}

	// Get next vocabulary word
	nextVocab, err := s.vocabRepo.GetByLevelAndIndex(user.CurrentLevel, progress.CurrentVocabIndex)
	if err != nil {
		return nil, fmt.Errorf("failed to get next vocabulary: %w", err)
	}

	// Create response with updated progress
	response := &models.VocabularyWithProgress{
		Vocabulary: nextVocab,
		Progress: &models.VocabularyProgress{
			CurrentIndex:      progress.CurrentVocabIndex,
			TotalWordsInLevel: totalWords,
			WordsLearned:      progress.WordsLearnedCount,
			StreakDays:        progress.StreakDays,
		},
	}

	return response, nil
}

func (s *VocabService) GetVocabularyByLevel(level string, page, limit int) (*models.VocabularyListResponse, error) {
	// Validate page and limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	// Get vocabulary list
	vocabList, total, err := s.vocabRepo.GetByLevel(level, page, limit)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	response := &models.VocabularyListResponse{
		Vocabulary: vocabList,
		Pagination: models.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}

	return response, nil
}

func (s *VocabService) GetProgress(userID string) (*models.ProgressStats, error) {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Get progress
	progress, err := s.progressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get total words for current level
	totalWords, err := s.vocabRepo.GetTotalCountByLevel(user.CurrentLevel)
	if err != nil {
		return nil, err
	}

	// Calculate days active
	var daysActive int
	if progress.LastStudyDate != nil {
		daysActive = int(time.Since(*progress.LastStudyDate).Hours() / 24)
	}

	stats := &models.ProgressStats{
		CurrentVocabIndex: progress.CurrentVocabIndex,
		CurrentLevel:      user.CurrentLevel,
		StreakDays:        progress.StreakDays,
		LastStudyDate:     progress.LastStudyDate,
		WordsLearned:      progress.WordsLearnedCount,
		WordsSkipped:      progress.WordsSkippedCount,
		TotalWordsInLevel: totalWords,
		TotalDaysActive:   daysActive,
		WordsLearnedByLevel: map[string]int{
			user.CurrentLevel: progress.WordsLearnedCount,
		},
	}

	return stats, nil
}

// updateStreak checks if it's a new day and updates the streak accordingly
func (s *VocabService) updateStreak(userID string, progress *models.UserProgress) error {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// If no last study date, this is the first time
	if progress.LastStudyDate == nil {
		return s.progressRepo.UpdateStreak(userID, 1, today)
	}

	lastStudy := time.Date(
		progress.LastStudyDate.Year(),
		progress.LastStudyDate.Month(),
		progress.LastStudyDate.Day(),
		0, 0, 0, 0,
		progress.LastStudyDate.Location(),
	)

	daysDiff := int(today.Sub(lastStudy).Hours() / 24)

	// Same day, no update needed
	if daysDiff == 0 {
		return nil
	}

	// Consecutive day, increment streak
	if daysDiff == 1 {
		return s.progressRepo.UpdateStreak(userID, progress.StreakDays+1, today)
	}

	// Streak broken, reset to 1
	return s.progressRepo.UpdateStreak(userID, 1, today)
}

// Helper functions for level word counts (can be moved to config later)
func (s *VocabService) getTotalWordsForLevel(level string) int {
	counts := map[string]int{
		"N5": 800,
		"N4": 300, // Starting with smaller set for MVP
		"N3": 650,
		"N2": 1000,
		"N1": 2000,
	}
	if count, ok := counts[level]; ok {
		return count
	}
	return 200 // Default fallback
}

// BulkCreateVocabulary is a helper for seeding data
func (s *VocabService) BulkCreateVocabulary(vocabList []models.Vocabulary) error {
	return s.vocabRepo.BulkCreate(vocabList)
}
