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

func (s *GrammarService) SkipToNextPattern(userID, patternID, status string) (*models.GrammarPatternResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Mark current pattern status (placeholder for future tracking)
	_ = patternID
	_ = status

	// Increment grammar index
	progress, err := s.progressRepo.IncrementGrammarIndex(userID)
	if err != nil {
		return nil, err
	}

	totalPatterns, err := s.grammarRepo.GetTotalCountByLevel(user.CurrentLevel)
	if err != nil {
		return nil, err
	}

	// Check if we've reached the end
	if progress.CurrentGrammarIndex >= totalPatterns {
		progress.CurrentGrammarIndex = 0
		if err := s.progressRepo.Update(progress); err != nil {
			return nil, err
		}
	}

	nextPattern, err := s.grammarRepo.GetByLevelAndIndex(user.CurrentLevel, progress.CurrentGrammarIndex)
	if err != nil {
		return nil, err
	}

	return &models.GrammarPatternResponse{
		Pattern: nextPattern,
		Progress: &models.GrammarProgress{
			CurrentIndex:    progress.CurrentGrammarIndex,
			TotalPatterns:   totalPatterns,
			PatternsLearned: progress.GrammarLearnedCount,
		},
	}, nil
}

// ComparisonPair represents two related patterns for side-by-side study
type ComparisonPair struct {
	ID          string                    `json:"id"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	PatternA    *models.GrammarPattern    `json:"pattern_a"`
	PatternB    *models.GrammarPattern    `json:"pattern_b"`
	KeyDifferences []string               `json:"key_differences"`
	UsageBoundaries string                 `json:"usage_boundaries"`
}

// DetailedComparison for two specific patterns
type DetailedComparison struct {
	PatternA        *models.GrammarPattern    `json:"pattern_a"`
	PatternB        *models.GrammarPattern    `json:"pattern_b"`
	KeyDifferences  []DifferencePoint        `json:"key_differences"`
	UsageBoundaries []BoundaryRule           `json:"usage_boundaries"`
	CommonErrors    []CommonError            `json:"common_errors"`
	DecisionTree    string                   `json:"decision_tree"`
}

// DifferencePoint highlights one specific contrast
type DifferencePoint struct {
	Aspect      string `json:"aspect"`       // e.g., "politeness", "restriction", "meaning"
	PatternA    string `json:"pattern_a"`    // How pattern A behaves
	PatternB    string `json:"pattern_b"`    // How pattern B behaves
	ExampleA    string `json:"example_a"`    // Example for A
	ExampleB    string `json:"example_b"`    // Example for B
}

// BoundaryRule defines when to use which pattern
type BoundaryRule struct {
	Situation   string `json:"situation"`    // When this applies
	UsePattern  string `json:"use_pattern"`  // "A" or "B"
	Explanation string `json:"explanation"`  // Why
}

// CommonError shows mistakes learners make
type CommonError struct {
	Error       string `json:"error"`        // The mistake
	Correction  string `json:"correction"`   // How to fix it
	Explanation string `json:"explanation"`  // Why it's wrong
}

// GetComparisonPairs returns related patterns grouped for comparison study
func (s *GrammarService) GetComparisonPairs(userID, level string) ([]ComparisonPair, error) {
	// Get all patterns for the level
	patterns, _, err := s.grammarRepo.GetByLevel(level, 1, 100)
	if err != nil {
		return nil, err
	}

	// Build comparison pairs based on related_patterns in data
	pairs := []ComparisonPair{}
	seen := make(map[string]bool)

	for _, p := range patterns {
		// Skip if already in a pair
		if seen[p.ID] {
			continue
		}

		// Look for related patterns marked for comparison
		for _, related := range p.RelatedPatterns {
			if related.Relationship == "contrast" || related.Relationship == "opposite" {
				// Find the related pattern
				relatedPattern, err := s.grammarRepo.GetByID(related.ID)
				if err != nil {
					continue
				}

				pairID := p.ID + "_vs_" + relatedPattern.ID
				if seen[relatedPattern.ID] || seen[p.ID] {
					continue
				}

				pair := ComparisonPair{
					ID:              pairID,
					Name:            p.Pattern + " vs " + relatedPattern.Pattern,
					Description:     related.KeyDifference,
					PatternA:        p,
					PatternB:        relatedPattern,
					KeyDifferences:  []string{related.KeyDifference},
					UsageBoundaries: "See detailed comparison for usage rules",
				}
				pairs = append(pairs, pair)
				seen[p.ID] = true
				seen[relatedPattern.ID] = true
				break // One pair per pattern for now
			}
		}
	}

	// If no pairs found from relationships, create some defaults
	if len(pairs) == 0 && len(patterns) >= 2 {
		// Create pairs from consecutive patterns (simple fallback)
		for i := 0; i < len(patterns)-1 && i < 6; i += 2 {
			pair := ComparisonPair{
				ID:          patterns[i].ID + "_vs_" + patterns[i+1].ID,
				Name:        patterns[i].Pattern + " vs " + patterns[i+1].Pattern,
				Description: "Compare these related patterns",
				PatternA:    patterns[i],
				PatternB:    patterns[i+1],
				KeyDifferences: []string{
					patterns[i].NuanceNotes,
					patterns[i+1].NuanceNotes,
				},
				UsageBoundaries: "Review pattern examples for usage differences",
			}
			pairs = append(pairs, pair)
		}
	}

	return pairs, nil
}

// ComparePatterns returns detailed analysis between two patterns
func (s *GrammarService) ComparePatterns(patternAID, patternBID string) (*DetailedComparison, error) {
	patternA, err := s.grammarRepo.GetByID(patternAID)
	if err != nil {
		return nil, err
	}

	patternB, err := s.grammarRepo.GetByID(patternBID)
	if err != nil {
		return nil, err
	}

	// Build key differences from usage examples
	differences := []DifferencePoint{}
	
	// Compare meanings
	if patternA.Meaning != patternB.Meaning {
		differences = append(differences, DifferencePoint{
			Aspect:   "core_meaning",
			PatternA: patternA.Meaning,
			PatternB: patternB.Meaning,
		})
	}

	// Compare nuance
	if patternA.NuanceNotes != "" && patternB.NuanceNotes != "" {
		differences = append(differences, DifferencePoint{
			Aspect:   "nuance",
			PatternA: patternA.NuanceNotes,
			PatternB: patternB.NuanceNotes,
		})
	}

	// Extract examples for comparison
	if len(patternA.UsageExamples) > 0 && len(patternB.UsageExamples) > 0 {
		differences = append(differences, DifferencePoint{
			Aspect:   "example_contrast",
			PatternA: patternA.UsageExamples[0].Japanese + " (" + patternA.UsageExamples[0].Meaning + ")",
			PatternB: patternB.UsageExamples[0].Japanese + " (" + patternB.UsageExamples[0].Meaning + ")",
		})
	}

	// Usage boundaries
	boundaries := []BoundaryRule{
		{
			Situation:   "When expressing " + patternA.Meaning,
			UsePattern:  "A",
			Explanation: patternA.NuanceNotes,
		},
		{
			Situation:   "When expressing " + patternB.Meaning,
			UsePattern:  "B",
			Explanation: patternB.NuanceNotes,
		},
	}

	// Common errors
	errors := []CommonError{}
	if patternA.CommonMistakes != "" {
		errors = append(errors, CommonError{
			Error:       "Using " + patternA.Pattern + " when " + patternB.Pattern + " is needed",
			Correction:  "Use " + patternB.Pattern + " instead for " + patternB.Meaning,
			Explanation: patternA.CommonMistakes,
		})
	}

	return &DetailedComparison{
		PatternA:        patternA,
		PatternB:        patternB,
		KeyDifferences:  differences,
		UsageBoundaries: boundaries,
		CommonErrors:    errors,
		DecisionTree:    s.buildDecisionTree(patternA, patternB),
	}, nil
}

// buildDecisionTree creates a decision guide for choosing between patterns
func (s *GrammarService) buildDecisionTree(a, b *models.GrammarPattern) string {
	return "Choose " + a.Pattern + " when: " + a.Meaning + ". Choose " + b.Pattern + " when: " + b.Meaning + "."
}
