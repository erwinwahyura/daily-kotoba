package services

import (
	"fmt"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
	"github.com/google/uuid"
)

// ListeningService handles listening practice business logic
type ListeningService struct {
	repo *repository.ListeningRepository
}

// NewListeningService creates a new service
func NewListeningService(repo *repository.ListeningRepository) *ListeningService {
	return &ListeningService{repo: repo}
}

// GetExercisesByLevel retrieves exercises for a JLPT level
func (s *ListeningService) GetExercisesByLevel(level string, limit int) (*models.ListeningListResponse, error) {
	exercises, err := s.repo.GetExercisesByLevel(level, limit)
	if err != nil {
		return nil, err
	}

	return &models.ListeningListResponse{
		Exercises:  exercises,
		TotalCount: len(exercises),
		Level:      level,
	}, nil
}

// GetExercise retrieves a specific exercise
func (s *ListeningService) GetExercise(id string) (*models.ListeningExercise, error) {
	return s.repo.GetExerciseByID(id)
}

// StartSession creates a new listening session
func (s *ListeningService) StartSession(userID, exerciseID string) (*models.ListeningSession, error) {
	// Verify exercise exists
	_, err := s.repo.GetExerciseByID(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise not found: %w", err)
	}

	session := &models.ListeningSession{
		ID:         uuid.New().String(),
		UserID:     userID,
		ExerciseID: exerciseID,
		StartedAt:  time.Now(),
		Status:     "in_progress",
		Answers:    []models.Answer{},
	}

	if err := s.repo.CreateSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// SubmitAnswer processes a user's answer
func (s *ListeningService) SubmitAnswer(sessionID string, questionID string, answer int, audioPosition int) (*models.Answer, error) {
	// Get session
	session, err := s.repo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Get exercise to check correct answer
	exercise, err := s.repo.GetExerciseByID(session.ExerciseID)
	if err != nil {
		return nil, fmt.Errorf("exercise not found: %w", err)
	}

	// Find question
	var question *models.ListeningQuestion
	for _, q := range exercise.Questions {
		if q.ID == questionID {
			question = &q
			break
		}
	}
	if question == nil {
		return nil, fmt.Errorf("question not found")
	}

	// Check if already answered
	for _, a := range session.Answers {
		if a.QuestionID == questionID {
			return nil, fmt.Errorf("question already answered")
		}
	}

	// Record answer
	isCorrect := answer == question.Correct
	userAnswer := models.Answer{
		QuestionID:    questionID,
		Answer:        answer,
		IsCorrect:     isCorrect,
		Timestamp:     time.Now(),
		AudioPosition: audioPosition,
	}

	session.Answers = append(session.Answers, userAnswer)

	// Update score
	correctCount := 0
	for _, a := range session.Answers {
		if a.IsCorrect {
			correctCount++
		}
	}
	session.Score = int(float64(correctCount) / float64(len(exercise.Questions)) * 100)

	// Check if completed
	if len(session.Answers) >= len(exercise.Questions) {
		session.Status = "completed"
		now := time.Now()
		session.CompletedAt = &now
	}

	// Save session
	if err := s.repo.UpdateSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &userAnswer, nil
}

// GetSession retrieves session progress
func (s *ListeningService) GetSession(sessionID string) (*models.ListeningSession, error) {
	return s.repo.GetSession(sessionID)
}

// GetUserStats retrieves user's listening statistics
func (s *ListeningService) GetUserStats(userID string) (*models.ListeningProgress, error) {
	return s.repo.GetUserStats(userID)
}

// UpdatePlayCount increments the play count
func (s *ListeningService) UpdatePlayCount(sessionID string, position int) error {
	session, err := s.repo.GetSession(sessionID)
	if err != nil {
		return err
	}

	session.PlayCount++
	session.CurrentPosition = position

	return s.repo.UpdateSession(session)
}

// SeedExercises adds sample exercises to the database
func (s *ListeningService) SeedExercises() error {
	return s.repo.SeedSampleExercises()
}
