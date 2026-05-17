package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
)

type JLPTService struct {
	jlptRepo *repository.JLPTRepository
}

func NewJLPTService(jlptRepo *repository.JLPTRepository) *JLPTService {
	return &JLPTService{jlptRepo: jlptRepo}
}

// GetAvailableTests returns tests for a specific level
func (s *JLPTService) GetAvailableTests(level string) ([]models.JLPTTest, error) {
	return s.jlptRepo.GetTestsByLevel(level)
}

// GetLevelInfo returns JLPT level metadata
func (s *JLPTService) GetLevelInfo() []map[string]interface{} {
	return models.GetJLPTLevelInfo()
}

// StartTest begins a new test session. count=0 means all questions.
func (s *JLPTService) StartTest(userID string, level, section string, count int) (*models.UserTestSession, []models.JLPTQuestion, error) {
	tests, err := s.jlptRepo.GetTestsByLevel(level)
	if err != nil {
		return nil, nil, err
	}
	if len(tests) == 0 {
		return nil, nil, fmt.Errorf("no tests available for level %s", level)
	}

	// Collect questions from all matching tests (or one section if specified)
	var allQuestions []models.JLPTQuestion
	var primaryTestID string
	for _, t := range tests {
		if section != "" && t.Section != section {
			continue
		}
		if primaryTestID == "" {
			primaryTestID = t.ID
		}
		qs, err := s.jlptRepo.GetQuestionsByTestID(t.ID)
		if err != nil {
			return nil, nil, err
		}
		allQuestions = append(allQuestions, qs...)
	}
	if len(allQuestions) == 0 {
		return nil, nil, fmt.Errorf("no questions available for level %s", level)
	}

	// Shuffle and limit
	rand.Shuffle(len(allQuestions), func(i, j int) {
		allQuestions[i], allQuestions[j] = allQuestions[j], allQuestions[i]
	})
	if count > 0 && count < len(allQuestions) {
		allQuestions = allQuestions[:count]
	}

	// Re-number questions after shuffle/slice
	for i := range allQuestions {
		allQuestions[i].QuestionNum = i + 1
	}

	session := &models.UserTestSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		TestID:    primaryTestID,
		Level:     level,
		StartedAt: time.Now(),
		Status:    "in_progress",
		Answers:   make(map[string]int),
	}

	if err := s.jlptRepo.CreateTestSession(session); err != nil {
		return nil, nil, err
	}

	return session, allQuestions, nil
}

// SubmitAnswer records an answer during a test
func (s *JLPTService) SubmitAnswer(sessionID, questionID string, answerIndex int) error {
	return s.jlptRepo.UpdateSessionAnswer(sessionID, questionID, answerIndex)
}

// CompleteTest finishes a test and calculates results
func (s *JLPTService) CompleteTest(sessionID string, finalAnswers map[string]int) (*models.TestResult, error) {
	// Get session
	session, err := s.jlptRepo.GetTestSession(sessionID)
	if err != nil {
		return nil, err
	}

	// Get test and questions
	test, err := s.jlptRepo.GetTestByID(session.TestID)
	if err != nil {
		return nil, err
	}

	questions, err := s.jlptRepo.GetQuestionsByTestID(session.TestID)
	if err != nil {
		return nil, err
	}

	// Calculate score
	correctCount := 0
	score := 0
	reviewItems := []models.ReviewItem{}

	for _, q := range questions {
		userAnswer, hasAnswer := finalAnswers[q.ID]
		if !hasAnswer {
			userAnswer = -1 // No answer
		}

		isCorrect := userAnswer == q.CorrectIndex
		if isCorrect {
			correctCount++
			score += q.PointValue
		} else {
			// Add to review items
			userAnswerText := "No answer"
			if userAnswer >= 0 && userAnswer < len(q.Options) {
				userAnswerText = q.Options[userAnswer]
			}
			
			// Safely get correct answer with bounds checking
			correctAnswer := "Unknown"
			if q.CorrectIndex >= 0 && q.CorrectIndex < len(q.Options) {
				correctAnswer = q.Options[q.CorrectIndex]
			}
			
			reviewItems = append(reviewItems, models.ReviewItem{
				QuestionNum:   q.QuestionNum,
				Question:      q.Question,
				YourAnswer:    userAnswerText,
				CorrectAnswer: correctAnswer,
				Explanation:   q.Explanation,
			})
		}
	}

	// Calculate time spent
	timeSpent := int(time.Since(session.StartedAt).Seconds())

	// Mark complete
	if err := s.jlptRepo.CompleteTestSession(sessionID, score, correctCount, timeSpent); err != nil {
		return nil, err
	}

	// Build result
	percentage := 0.0
	if len(questions) > 0 {
		percentage = float64(correctCount) / float64(len(questions)) * 100
	}

	passed := score >= test.PassingScore

	// Format time
	timeSpentStr := fmt.Sprintf("%d:%02d", timeSpent/60, timeSpent%60)

	return &models.TestResult{
		SessionID:       sessionID,
		Level:           session.Level,
		Score:           score,
		TotalQuestions:  len(questions),
		CorrectCount:    correctCount,
		IncorrectCount:  len(questions) - correctCount,
		Percentage:      percentage,
		Passed:          passed,
		TimeSpent:       timeSpentStr,
		TimeLimit:       test.TimeLimit,
		ReviewQuestions: reviewItems,
	}, nil
}

// GetUserHistory returns test history for a user
func (s *JLPTService) GetUserHistory(userID string) ([]models.UserTestSession, error) {
	return s.jlptRepo.GetUserTestHistory(userID)
}

// GetTestProgress returns current progress for an active test
func (s *JLPTService) GetTestProgress(sessionID string) (answeredCount, totalCount int, err error) {
	session, err := s.jlptRepo.GetTestSession(sessionID)
	if err != nil {
		return 0, 0, err
	}

	questions, err := s.jlptRepo.GetQuestionsByTestID(session.TestID)
	if err != nil {
		return 0, 0, err
	}

	return len(session.Answers), len(questions), nil
}