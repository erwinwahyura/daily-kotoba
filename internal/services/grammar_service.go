package services

import (
	"math"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type GrammarService struct {
	grammarRepo  *repository.GrammarRepository
	progressRepo *repository.ProgressRepository
	userRepo     *repository.UserRepository
}

func NewGrammarService(
	grammarRepo *repository.GrammarRepository,
	progressRepo *repository.ProgressRepository,
	userRepo *repository.UserRepository,
) *GrammarService {
	return &GrammarService{
		grammarRepo:  grammarRepo,
		progressRepo: progressRepo,
		userRepo:     userRepo,
	}
}

func (s *GrammarService) GetDailyPattern(userID string) (*models.GrammarPatternResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	progress, err := s.progressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	pattern, err := s.grammarRepo.GetByLevelAndIndex(user.CurrentLevel, progress.CurrentGrammarIndex)
	if err != nil {
		return nil, err
	}

	totalPatterns, err := s.grammarRepo.GetTotalCountByLevel(user.CurrentLevel)
	if err != nil {
		return nil, err
	}

	return &models.GrammarPatternResponse{
		Pattern: pattern,
		Progress: &models.GrammarProgress{
			CurrentIndex:    progress.CurrentGrammarIndex,
			TotalPatterns:   totalPatterns,
			PatternsLearned: progress.GrammarLearnedCount,
		},
	}, nil
}

func (s *GrammarService) GetPatternByID(patternID string) (*models.GrammarPattern, error) {
	return s.grammarRepo.GetByID(patternID)
}

func (s *GrammarService) GetPatternsByLevel(level string, page, limit int) (*models.GrammarListResponse, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	patterns, total, err := s.grammarRepo.GetByLevel(level, page, limit)
	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &models.GrammarListResponse{
		Patterns: patterns,
		Pagination: models.PaginationResponse{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *GrammarService) BulkCreatePatterns(patterns []models.GrammarPattern) error {
	return s.grammarRepo.BulkCreate(patterns)
}
