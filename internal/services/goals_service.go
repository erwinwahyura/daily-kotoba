package services

import (
	"fmt"
	"math"
	"time"

	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
)

// GoalsService handles daily goals, streaks, and achievements
type GoalsService struct {
	goalsRepo *repository.GoalsRepository
}

// NewGoalsService creates a new service
func NewGoalsService(goalsRepo *repository.GoalsRepository) *GoalsService {
	return &GoalsService{
		goalsRepo: goalsRepo,
	}
}

// GetDailyProgress gets today's progress summary
func (s *GoalsService) GetDailyProgress(userID string) (*models.DailyProgress, error) {
	// Get or create today's goal
	goal, err := s.goalsRepo.GetOrCreateDailyGoal(userID, time.Now())
	if err != nil {
		return nil, err
	}
	
	// Get streak
	streak, err := s.goalsRepo.GetUserStreak(userID)
	if err != nil {
		return nil, err
	}
	
	// Calculate overall progress
	totalTarget := goal.VocabTarget + goal.GrammarTarget + goal.KanjiTarget + 
		goal.ConjugationTarget + goal.ReadingTarget
	totalCompleted := goal.VocabCompleted + goal.GrammarCompleted + goal.KanjiCompleted + 
		goal.ConjugationCompleted + goal.ReadingCompleted
	
	progress := 0
	if totalTarget > 0 {
		progress = int(math.Min(100, float64(totalCompleted)/float64(totalTarget)*100))
	}
	
	return &models.DailyProgress{
		Date:              time.Now().Format("2006-01-02"),
		VocabCompleted:    goal.VocabCompleted,
		VocabTarget:       goal.VocabTarget,
		GrammarCompleted:  goal.GrammarCompleted,
		GrammarTarget:     goal.GrammarTarget,
		KanjiCompleted:    goal.KanjiCompleted,
		KanjiTarget:       goal.KanjiTarget,
		ConjugationCompleted: goal.ConjugationCompleted,
		ConjugationTarget:    goal.ConjugationTarget,
		ReadingCompleted:  goal.ReadingCompleted,
		ReadingTarget:     goal.ReadingTarget,
		OverallProgress:   progress,
		IsCompleted:       goal.IsCompleted,
		CurrentStreak:     streak.CurrentStreak,
	}, nil
}

// UpdateProgress updates progress for a specific activity type
func (s *GoalsService) UpdateProgress(userID string, activityType string, count int) error {
	// Get today's goal
	goal, err := s.goalsRepo.GetOrCreateDailyGoal(userID, time.Now())
	if err != nil {
		return err
	}
	
	// Update appropriate counter
	switch activityType {
	case "vocab":
		goal.VocabCompleted += count
	case "grammar":
		goal.GrammarCompleted += count
	case "kanji":
		goal.KanjiCompleted += count
	case "conjugation":
		goal.ConjugationCompleted += count
	case "reading":
		goal.ReadingCompleted += count
	default:
		return fmt.Errorf("unknown activity type: %s", activityType)
	}
	
	// Save goal
	if err := s.goalsRepo.UpdateDailyGoal(goal); err != nil {
		return err
	}
	
	// Record activity for streak
	if err := s.goalsRepo.RecordActivity(userID); err != nil {
		return err
	}
	
	return nil
}

// GetGoalSettings retrieves user's goal settings
func (s *GoalsService) GetGoalSettings(userID string) (*models.GoalSettings, error) {
	return s.goalsRepo.GetGoalSettings(userID)
}

// UpdateGoalSettings updates user's goal settings
func (s *GoalsService) UpdateGoalSettings(userID string, settings *models.GoalSettings) error {
	settings.UserID = userID
	return s.goalsRepo.UpdateGoalSettings(settings)
}

// GetUserStreak gets user's streak information
func (s *GoalsService) GetUserStreak(userID string) (*models.UserStreak, error) {
	return s.goalsRepo.GetUserStreak(userID)
}

// GetAchievements retrieves user's achievements
func (s *GoalsService) GetAchievements(userID string) ([]models.Achievement, error) {
	return s.goalsRepo.GetUserAchievements(userID)
}

// CheckAndAwardAchievements checks and awards new achievements
func (s *GoalsService) CheckAndAwardAchievements(userID string) ([]models.Achievement, error) {
	// Get current stats
	streak, err := s.goalsRepo.GetUserStreak(userID)
	if err != nil {
		return nil, err
	}
	
	// Get existing achievements
	existingAchievements, err := s.goalsRepo.GetUserAchievements(userID)
	if err != nil {
		return nil, err
	}
	
	// Create map of earned achievement IDs
	earnedMap := make(map[string]bool)
	for _, a := range existingAchievements {
		earnedMap[a.Type] = true // Using Type as key for simplicity
	}
	
	// Get all achievement definitions
	definitions := models.GetAchievementDefinitions()
	
	var newAchievements []models.Achievement
	
	// Check each achievement
	for _, def := range definitions {
		if earnedMap[def.ID] {
			continue // Already earned
		}
		
		earned := false
		
		switch def.Type {
		case "streak":
			if streak.CurrentStreak >= def.Requirement {
				earned = true
			}
		// Add more cases for other achievement types
		// These would require querying other repositories for counts
		}
		
		if earned {
			if err := s.goalsRepo.AddAchievement(userID, def); err != nil {
				continue
			}
			
			newAchievements = append(newAchievements, models.Achievement{
				ID:          def.ID,
				UserID:      userID,
				Type:        def.Type,
				Name:        def.Name,
				Description: def.Description,
				Icon:        def.Icon,
				Level:       def.Level,
				EarnedAt:    time.Now(),
				IsNew:       true,
			})
		}
	}
	
	return newAchievements, nil
}

// GetWeeklyProgress gets progress for the last 7 days
func (s *GoalsService) GetWeeklyProgress(userID string) ([]models.DailyProgress, error) {
	goals, err := s.goalsRepo.GetRecentDailyGoals(userID, 7)
	if err != nil {
		return nil, err
	}
	
	var progress []models.DailyProgress
	for _, g := range goals {
		totalTarget := g.VocabTarget + g.GrammarTarget + g.KanjiTarget + 
			g.ConjugationTarget + g.ReadingTarget
		totalCompleted := g.VocabCompleted + g.GrammarCompleted + g.KanjiCompleted + 
			g.ConjugationCompleted + g.ReadingCompleted
		
		p := 0
		if totalTarget > 0 {
			p = int(math.Min(100, float64(totalCompleted)/float64(totalTarget)*100))
		}
		
		progress = append(progress, models.DailyProgress{
			Date:              g.Date.Format("2006-01-02"),
			VocabCompleted:    g.VocabCompleted,
			VocabTarget:       g.VocabTarget,
			GrammarCompleted:  g.GrammarCompleted,
			GrammarTarget:     g.GrammarTarget,
			KanjiCompleted:    g.KanjiCompleted,
			KanjiTarget:       g.KanjiTarget,
			ConjugationCompleted: g.ConjugationCompleted,
			ConjugationTarget:    g.ConjugationTarget,
			ReadingCompleted:  g.ReadingCompleted,
			ReadingTarget:     g.ReadingTarget,
			OverallProgress:   p,
			IsCompleted:       g.IsCompleted,
		})
	}
	
	return progress, nil
}

// GetAllAchievementDefinitions returns all available achievements
func (s *GoalsService) GetAllAchievementDefinitions() []models.AchievementDefinition {
	return models.GetAchievementDefinitions()
}
