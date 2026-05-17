package services

import (
	"fmt"
	"math"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type KanjiService struct {
	kanjiRepo *repository.KanjiRepository
}

func NewKanjiService(kanjiRepo *repository.KanjiRepository) *KanjiService {
	return &KanjiService{kanjiRepo: kanjiRepo}
}

func (s *KanjiService) GetCharacter(char string) (*models.Kanji, error) {
	k, err := s.kanjiRepo.GetByCharacter(char)
	if err != nil {
		return nil, fmt.Errorf("kanji not found: %w", err)
	}
	return k, nil
}

func (s *KanjiService) StartSession(userID, kanjiChar string) (*models.KanjiPracticeSession, error) {
	k, err := s.kanjiRepo.GetByCharacter(kanjiChar)
	if err != nil {
		return nil, fmt.Errorf("kanji not found")
	}
	return s.kanjiRepo.CreateSession(userID, k.ID, kanjiChar)
}

// CompareStroke compares user's drawn stroke against the reference stroke.
// Returns accuracy (0-100) and feedback message.
func (s *KanjiService) CompareStroke(sessionID string, strokeNum int, userPath []models.Point) (*models.KanjiCompareResult, error) {
	session, err := s.kanjiRepo.GetSession(sessionID)
	if err != nil || session == nil {
		return nil, fmt.Errorf("session not found")
	}

	kanji, err := s.kanjiRepo.GetByCharacter(session.KanjiChar)
	if err != nil || kanji == nil {
		return nil, fmt.Errorf("kanji not found")
	}

	// Find reference stroke
	var refStroke *models.Stroke
	for i := range kanji.StrokeOrder {
		if kanji.StrokeOrder[i].StrokeNum == strokeNum {
			refStroke = &kanji.StrokeOrder[i]
			break
		}
	}
	if refStroke == nil {
		// No reference data — give benefit of the doubt
		return &models.KanjiCompareResult{
			Accuracy: 75, Feedback: "Keep practicing!", Direction: "unknown", OrderCorrect: true,
		}, nil
	}

	accuracy := compareStrokePaths(userPath, *refStroke)

	var feedback string
	var directionOK string
	switch {
	case accuracy >= 85:
		feedback = "Excellent! Perfect stroke."
		directionOK = "correct"
	case accuracy >= 70:
		feedback = "Good! Keep practicing."
		directionOK = "correct"
	case accuracy >= 50:
		feedback = "Almost there. Check stroke direction."
		directionOK = "wrong_direction"
	default:
		feedback = "Try again. Follow the guide carefully."
		directionOK = "wrong_direction"
	}

	// Update session accuracy (running average)
	newAccuracy := (session.Accuracy*float64(strokeNum-1) + accuracy) / float64(strokeNum)
	_ = s.kanjiRepo.UpdateSessionAccuracy(sessionID, newAccuracy, strokeNum)

	return &models.KanjiCompareResult{
		Accuracy:     accuracy,
		Feedback:     feedback,
		Direction:    directionOK,
		OrderCorrect: strokeNum == session.CurrentStroke+1,
	}, nil
}

// compareStrokePaths computes a 0-100 similarity between user path and reference stroke.
// Uses direction match + start/end proximity as a simple heuristic.
func compareStrokePaths(userPath []models.Point, ref models.Stroke) float64 {
	if len(userPath) < 2 {
		return 0
	}

	// Normalize user path to 0-1 coordinate space (assuming canvas 300x300)
	norm := normalizePoints(userPath)

	// Score based on start proximity
	startDist := dist(norm[0], ref.StartPoint)
	endDist := dist(norm[len(norm)-1], ref.EndPoint)

	// Direction match: compute overall direction vector of user stroke
	userDir := dirVector(norm[0], norm[len(norm)-1])
	refDir := dirVector(ref.StartPoint, ref.EndPoint)
	dirSimilarity := dotProduct(userDir, refDir) // -1 to 1

	// Combine: start proximity (0-1), end proximity (0-1), direction (-1 to 1)
	startScore := math.Max(0, 1-startDist*3)
	endScore := math.Max(0, 1-endDist*3)
	dirScore := (dirSimilarity + 1) / 2 // 0 to 1

	accuracy := (startScore*0.3 + endScore*0.3 + dirScore*0.4) * 100
	return math.Round(accuracy)
}

func normalizePoints(pts []models.Point) []models.Point {
	if len(pts) == 0 {
		return pts
	}
	minX, minY := pts[0].X, pts[0].Y
	maxX, maxY := pts[0].X, pts[0].Y
	for _, p := range pts {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	rangeX := maxX - minX
	rangeY := maxY - minY
	if rangeX == 0 {
		rangeX = 1
	}
	if rangeY == 0 {
		rangeY = 1
	}
	norm := make([]models.Point, len(pts))
	for i, p := range pts {
		norm[i] = models.Point{X: (p.X - minX) / rangeX, Y: (p.Y - minY) / rangeY}
	}
	return norm
}

func dist(a, b models.Point) float64 {
	dx := a.X - b.X
	dy := a.Y - b.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func dirVector(start, end models.Point) models.Point {
	dx := end.X - start.X
	dy := end.Y - start.Y
	mag := math.Sqrt(dx*dx + dy*dy)
	if mag == 0 {
		return models.Point{X: 0, Y: 0}
	}
	return models.Point{X: dx / mag, Y: dy / mag}
}

func dotProduct(a, b models.Point) float64 {
	return a.X*b.X + a.Y*b.Y
}
