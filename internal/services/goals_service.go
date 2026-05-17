package services

import (
	"fmt"
	"time"

	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

type GoalsService struct {
	goalsRepo    *repository.GoalsRepository
	progressRepo *repository.ProgressRepository
}

func NewGoalsService(goalsRepo *repository.GoalsRepository, progressRepo *repository.ProgressRepository) *GoalsService {
	return &GoalsService{goalsRepo: goalsRepo, progressRepo: progressRepo}
}

func (s *GoalsService) GetDailyProgress(userID string) (*models.DailyProgressResponse, error) {
	settings, err := s.goalsRepo.GetOrCreateSettings(userID)
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	types := []string{"vocab", "grammar", "kanji", "conjugation", "reading"}
	counts := make(map[string]int)
	for _, t := range types {
		c, _ := s.goalsRepo.GetDailyCount(userID, today, t)
		counts[t] = c
	}

	return &models.DailyProgressResponse{
		VocabCompleted:       counts["vocab"],
		VocabTarget:          settings.VocabTarget,
		GrammarCompleted:     counts["grammar"],
		GrammarTarget:        settings.GrammarTarget,
		KanjiCompleted:       counts["kanji"],
		KanjiTarget:          settings.KanjiTarget,
		ConjugationCompleted: counts["conjugation"],
		ConjugationTarget:    settings.ConjugationTarget,
		ReadingCompleted:     counts["reading"],
		ReadingTarget:        settings.ReadingTarget,
	}, nil
}

func (s *GoalsService) GetStreak(userID string) (*models.StreakResponse, error) {
	progress, err := s.progressRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	lastStudy := ""
	if progress.LastStudyDate != nil {
		lastStudy = progress.LastStudyDate.Format("2006-01-02")
	}
	return &models.StreakResponse{
		CurrentStreak: progress.StreakDays,
		LongestStreak: progress.StreakDays,
		LastStudyDate: lastStudy,
	}, nil
}

func (s *GoalsService) GetWeeklyProgress(userID string) ([]models.WeeklyDay, error) {
	settings, err := s.goalsRepo.GetOrCreateSettings(userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	days := make([]string, 7)
	for i := 6; i >= 0; i-- {
		days[6-i] = now.AddDate(0, 0, -i).Format("2006-01-02")
	}

	activityMap, err := s.goalsRepo.GetWeeklyActivity(userID, days)
	if err != nil {
		return nil, err
	}

	total := settings.VocabTarget + settings.GrammarTarget + settings.KanjiTarget +
		settings.ConjugationTarget + settings.ReadingTarget

	result := make([]models.WeeklyDay, 7)
	for i, d := range days {
		acts := activityMap[d]
		done := acts["vocab"] + acts["grammar"] + acts["kanji"] + acts["conjugation"] + acts["reading"]
		progress := 0
		if total > 0 {
			progress = (done * 100) / total
			if progress > 100 {
				progress = 100
			}
		}
		result[i] = models.WeeklyDay{
			Date:            d,
			OverallProgress: progress,
			IsCompleted:     progress >= 100,
		}
	}
	return result, nil
}

func (s *GoalsService) GetSettings(userID string) (*models.UserGoalSettings, error) {
	return s.goalsRepo.GetOrCreateSettings(userID)
}

func (s *GoalsService) SaveSettings(userID string, settings *models.UserGoalSettings) error {
	settings.UserID = userID
	return s.goalsRepo.SaveSettings(settings)
}

func (s *GoalsService) RecordActivity(userID, activityType string, count int) error {
	if count <= 0 {
		count = 1
	}
	today := time.Now().Format("2006-01-02")
	return s.goalsRepo.IncrementActivity(userID, today, activityType, count)
}

func (s *GoalsService) GetAllAchievements() ([]models.AchievementDef, error) {
	return s.goalsRepo.GetAllAchievements()
}

func (s *GoalsService) GetUserAchievements(userID string) ([]models.UserAchievement, error) {
	return s.goalsRepo.GetUserAchievements(userID)
}

func (s *GoalsService) CheckAndGrantAchievements(userID string) error {
	all, err := s.goalsRepo.GetAllAchievements()
	if err != nil {
		return err
	}
	earned, err := s.goalsRepo.GetUserAchievements(userID)
	if err != nil {
		return err
	}
	earnedIDs := make(map[string]bool)
	for _, e := range earned {
		earnedIDs[e.AchievementID] = true
	}

	progress, err := s.progressRepo.GetByUserID(userID)
	if err != nil {
		return err
	}

	today := time.Now().Format("2006-01-02")
	totalConj := 0
	{
		c, _ := s.goalsRepo.GetDailyCount(userID, today, "conjugation")
		totalConj = c
	}

	for _, a := range all {
		if earnedIDs[a.ID] {
			continue
		}
		var met bool
		switch a.RequirementType {
		case "streak_days":
			met = progress.StreakDays >= a.RequirementValue
		case "total_count":
			switch a.Category {
			case "vocab":
				met = progress.WordsLearnedCount >= a.RequirementValue
			case "grammar":
				met = progress.GrammarLearnedCount >= a.RequirementValue
			case "conjugation":
				met = totalConj >= a.RequirementValue
			}
		}
		if met {
			_ = s.goalsRepo.GrantAchievement(userID, a.ID)
			fmt.Printf("Granted achievement %s to user %s\n", a.Name, userID)
		}
	}
	return nil
}
