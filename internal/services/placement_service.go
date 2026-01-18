package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type PlacementService struct {
	placementRepo *repository.PlacementRepository
	userRepo      *repository.UserRepository
}

func NewPlacementService(placementRepo *repository.PlacementRepository, userRepo *repository.UserRepository) *PlacementService {
	return &PlacementService{
		placementRepo: placementRepo,
		userRepo:      userRepo,
	}
}

// GetPlacementTest retrieves all placement test questions with shuffled options
func (s *PlacementService) GetPlacementTest() ([]models.PlacementQuestionResponse, error) {
	questions, err := s.placementRepo.GetAllQuestions()
	if err != nil {
		return nil, err
	}

	// Convert to response format with shuffled options
	responses := make([]models.PlacementQuestionResponse, len(questions))
	for i, q := range questions {
		// Combine correct answer with wrong answers
		allOptions := append([]string{q.CorrectAnswer}, q.WrongAnswers...)

		// Shuffle options
		shuffledOptions := s.shuffleOptions(allOptions)

		responses[i] = models.PlacementQuestionResponse{
			ID:              q.ID,
			QuestionText:    q.QuestionText,
			Options:         shuffledOptions,
			DifficultyLevel: q.DifficultyLevel,
			OrderIndex:      q.OrderIndex,
		}
	}

	return responses, nil
}

// shuffleOptions randomly shuffles the answer options
func (s *PlacementService) shuffleOptions(options []string) []string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]string, len(options))
	copy(shuffled, options)

	r.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}

// SubmitPlacementTest evaluates the test and assigns a JLPT level
func (s *PlacementService) SubmitPlacementTest(userID string, submission *models.PlacementTestSubmission) (*models.PlacementTestResponse, error) {
	// Get all questions
	questions, err := s.placementRepo.GetAllQuestions()
	if err != nil {
		return nil, err
	}

	// Create map of question ID to question for easy lookup
	questionMap := make(map[string]models.PlacementQuestion)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	// Score the test
	totalScore := 0
	breakdown := map[string]int{
		"N5": 0,
		"N4": 0,
		"N3": 0,
		"N2": 0,
		"N1": 0,
	}

	for questionID, userAnswer := range submission.Answers {
		question, exists := questionMap[questionID]
		if !exists {
			continue // Skip invalid question IDs
		}

		if question.CorrectAnswer == userAnswer {
			totalScore++
			breakdown[question.DifficultyLevel]++
		}
	}

	// Determine assigned level based on score
	assignedLevel := s.calculateLevel(totalScore, breakdown)

	// Save test result
	result := &models.PlacementTestResult{
		UserID:        userID,
		TestScore:     totalScore,
		AssignedLevel: assignedLevel,
		CompletedAt:   time.Now(),
	}

	err = s.placementRepo.SaveTestResult(result)
	if err != nil {
		return nil, fmt.Errorf("failed to save test result: %w", err)
	}

	// Update user's current level
	err = s.userRepo.UpdateUserLevel(userID, assignedLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to update user level: %w", err)
	}

	// Return response
	return &models.PlacementTestResponse{
		Score:          totalScore,
		TotalQuestions: len(questions),
		AssignedLevel:  assignedLevel,
		Breakdown:      breakdown,
	}, nil
}

// calculateLevel determines JLPT level based on test score and breakdown
func (s *PlacementService) calculateLevel(totalScore int, breakdown map[string]int) string {
	// Scoring algorithm:
	// - Focus on highest level where user got most questions correct
	// - If user scored well on higher levels, assign that level
	// - Otherwise, assign based on total score

	// If user got 2+ N1 questions correct (out of 2), assign N1
	if breakdown["N1"] >= 2 {
		return "N1"
	}

	// If user got 2+ N2 questions correct (out of 3), assign N2
	if breakdown["N2"] >= 2 {
		return "N2"
	}

	// If user got 3+ N3 questions correct (out of 5), assign N3
	if breakdown["N3"] >= 3 {
		return "N3"
	}

	// If user got 3+ N4 questions correct (out of 5), assign N4
	if breakdown["N4"] >= 3 {
		return "N4"
	}

	// Default to N5 (beginner level)
	return "N5"
}

// GetUserTestResult retrieves a user's placement test result
func (s *PlacementService) GetUserTestResult(userID string) (*models.PlacementTestResult, error) {
	return s.placementRepo.GetUserTestResult(userID)
}
