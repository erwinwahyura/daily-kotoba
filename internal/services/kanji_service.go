package services

import (
	"fmt"
	"math"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
	"github.com/google/uuid"
)

// KanjiService handles kanji writing practice business logic
type KanjiService struct {
	kanjiRepo *repository.KanjiRepository
}

// NewKanjiService creates a new service
func NewKanjiService(kanjiRepo *repository.KanjiRepository) *KanjiService {
	return &KanjiService{
		kanjiRepo: kanjiRepo,
	}
}

// GetKanjiByCharacter retrieves kanji details
func (s *KanjiService) GetKanjiByCharacter(char string) (*models.Kanji, error) {
	return s.kanjiRepo.GetKanjiByCharacter(char)
}

// GetKanjiByLevel retrieves kanji for a JLPT level
func (s *KanjiService) GetKanjiByLevel(level string, limit int) (*models.KanjiListResponse, error) {
	kanji, err := s.kanjiRepo.GetKanjiByLevel(level, limit)
	if err != nil {
		return nil, err
	}

	return &models.KanjiListResponse{
		Kanji:      kanji,
		TotalCount: len(kanji),
		Level:      level,
	}, nil
}

// StartPracticeSession creates a new practice session
func (s *KanjiService) StartPracticeSession(userID, kanjiChar string) (*models.KanjiPracticeSession, error) {
	// Get kanji details
	kanji, err := s.kanjiRepo.GetKanjiByCharacter(kanjiChar)
	if err != nil {
		return nil, err
	}

	session := &models.KanjiPracticeSession{
		ID:        uuid.New().String(),
		UserID:    userID,
		KanjiID:   kanji.ID,
		KanjiChar: kanji.Character,
		StartedAt: time.Now(),
		Status:    "in_progress",
		Accuracy:  0,
		UserStrokes: []models.UserStroke{},
	}

	if err := s.kanjiRepo.CreatePracticeSession(session); err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return session, nil
}

// CompareStroke compares user's stroke with reference and returns accuracy
func (s *KanjiService) CompareStroke(sessionID string, strokeNum int, userPath []models.Point) (*models.KanjiCompareResult, error) {
	// Get session
	session, err := s.kanjiRepo.GetPracticeSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Get kanji reference
	kanji, err := s.kanjiRepo.GetKanjiByCharacter(session.KanjiChar)
	if err != nil {
		return nil, fmt.Errorf("kanji not found: %w", err)
	}

	// Find reference stroke
	var refStroke *models.Stroke
	for _, s := range kanji.StrokeOrder {
		if s.StrokeNum == strokeNum {
			refStroke = &s
			break
		}
	}

	if refStroke == nil {
		return nil, fmt.Errorf("stroke %d not found for kanji %s", strokeNum, kanji.Character)
	}

	// Check if correct stroke order
	orderCorrect := strokeNum == len(session.UserStrokes)+1

	// Calculate accuracy using path comparison
	accuracy := s.calculateStrokeAccuracy(refStroke, userPath)

	// Generate feedback
	feedback := s.generateFeedback(accuracy, refStroke.Direction)

	// Record user's stroke
	userStroke := models.UserStroke{
		StrokeNum: strokeNum,
		Path:      userPath,
		Timestamp: time.Now(),
	}
	if len(userPath) >= 2 {
		// Estimate duration based on path length (simplified)
		userStroke.Duration = len(userPath) * 10
	}

	session.UserStrokes = append(session.UserStrokes, userStroke)

	// Update overall accuracy
	if len(session.UserStrokes) > 0 {
		totalAccuracy := 0.0
		for _, us := range session.UserStrokes {
			// Recalculate accuracy for each stroke
			for _, ref := range kanji.StrokeOrder {
				if ref.StrokeNum == us.StrokeNum {
					totalAccuracy += s.calculateStrokeAccuracy(&ref, us.Path)
					break
				}
			}
		}
		session.Accuracy = totalAccuracy / float64(len(session.UserStrokes))
	}

	// Check if completed
	if len(session.UserStrokes) >= kanji.StrokeCount {
		session.Status = "completed"
		now := time.Now()
		session.CompletedAt = &now
	}

	// Save session
	if err := s.kanjiRepo.UpdatePracticeSession(session); err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}

	return &models.KanjiCompareResult{
		Accuracy:     accuracy,
		Feedback:     feedback,
		Direction:    refStroke.Direction,
		OrderCorrect: orderCorrect,
	}, nil
}

// calculateStrokeAccuracy compares user path with reference stroke
func (s *KanjiService) calculateStrokeAccuracy(ref *models.Stroke, userPath []models.Point) float64 {
	if len(userPath) < 2 {
		return 0.0
	}

	// Simplified stroke comparison algorithm
	// 1. Check direction similarity
	// 2. Check start/end point proximity
	// 3. Check overall path shape

	// Get user stroke direction vector
	userStart := userPath[0]
	userEnd := userPath[len(userPath)-1]
	userDx := userEnd.X - userStart.X
	userDy := userEnd.Y - userStart.Y
	userLength := math.Sqrt(userDx*userDx + userDy*userDy)

	// Get reference direction vector
	refDx := ref.EndPoint.X - ref.StartPoint.X
	refDy := ref.EndPoint.Y - ref.StartPoint.Y
	refLength := math.Sqrt(refDx*refDx + refDy*refDy)

	if userLength == 0 || refLength == 0 {
		return 50.0 // Neutral if no movement
	}

	// Normalize vectors
	userDx /= userLength
	userDy /= userLength
	refDx /= refLength
	refDy /= refLength

	// Calculate direction similarity (dot product)
	dotProduct := userDx*refDx + userDy*refDy
	directionScore := (dotProduct + 1) / 2 * 100 // Convert to 0-100

	// Calculate start point proximity (normalized to 100x100 canvas)
	startDist := math.Sqrt(
		math.Pow(userStart.X-ref.StartPoint.X, 2) +
		math.Pow(userStart.Y-ref.StartPoint.Y, 2),
	)
	startScore := math.Max(0, 100-startDist)

	// Calculate end point proximity
	endDist := math.Sqrt(
		math.Pow(userEnd.X-ref.EndPoint.X, 2) +
		math.Pow(userEnd.Y-ref.EndPoint.Y, 2),
	)
	endScore := math.Max(0, 100-endDist)

	// Calculate length ratio
	lengthRatio := math.Min(userLength, refLength) / math.Max(userLength, refLength)
	lengthScore := lengthRatio * 100

	// Weighted average
	accuracy := directionScore*0.5 + startScore*0.2 + endScore*0.2 + lengthScore*0.1

	// Clamp to 0-100
	return math.Min(100, math.Max(0, accuracy))
}

// generateFeedback creates helpful feedback based on accuracy
func (s *KanjiService) generateFeedback(accuracy float64, direction string) string {
	switch {
	case accuracy >= 90:
		return "Excellent! Perfect stroke."
	case accuracy >= 80:
		return "Great job! Very close."
	case accuracy >= 70:
		return "Good! Try to make it " + direction + "."
	case accuracy >= 60:
		return "Getting there. Watch the " + direction + " direction."
	case accuracy >= 50:
		return "Keep practicing. Focus on the stroke direction."
	default:
		return "Try again. Follow the guide line."
	}
}

// GetPracticeSession retrieves session details
func (s *KanjiService) GetPracticeSession(sessionID string) (*models.KanjiPracticeSession, error) {
	return s.kanjiRepo.GetPracticeSession(sessionID)
}

// GetUserStats gets user's kanji practice statistics
func (s *KanjiService) GetUserStats(userID string) (map[string]interface{}, error) {
	return s.kanjiRepo.GetUserKanjiStats(userID)
}

// SeedKanjiData seeds initial kanji data
func (s *KanjiService) SeedKanjiData() error {
	return s.kanjiRepo.SeedSampleKanji()
}
