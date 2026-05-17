package services

import (
	"fmt"

	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
)

type ReadingService struct {
	repo *repository.ReadingRepository
}

func NewReadingService(repo *repository.ReadingRepository) *ReadingService {
	return &ReadingService{repo: repo}
}

func (s *ReadingService) GetArticlesByLevel(level string) ([]models.ReadingArticle, error) {
	if level == "" {
		return nil, fmt.Errorf("level is required")
	}
	return s.repo.GetArticlesByLevel(level)
}

func (s *ReadingService) GetArticle(id string) (*models.ReadingArticle, error) {
	if id == "" {
		return nil, fmt.Errorf("article id is required")
	}
	return s.repo.GetArticleByID(id)
}

func (s *ReadingService) GetRandomArticle(level string) (*models.ReadingArticle, error) {
	if level == "" {
		return nil, fmt.Errorf("level is required")
	}
	return s.repo.GetRandomArticle(level)
}

func (s *ReadingService) SubmitAnswers(req models.SubmitReadingRequest) (map[string]interface{}, error) {
	article, err := s.repo.GetArticleByID(req.ArticleID)
	if err != nil {
		return nil, err
	}

	correct := 0
	results := make([]map[string]interface{}, 0, len(article.Questions))
	for _, q := range article.Questions {
		userAnswer, answered := req.Answers[q.ID]
		isCorrect := answered && userAnswer == q.CorrectIndex
		if isCorrect {
			correct++
		}
		results = append(results, map[string]interface{}{
			"question_id":   q.ID,
			"correct_index": q.CorrectIndex,
			"user_answer":   userAnswer,
			"is_correct":    isCorrect,
			"explanation":   q.Explanation,
		})
	}

	total := len(article.Questions)
	score := 0
	if total > 0 {
		score = correct * 100 / total
	}

	return map[string]interface{}{
		"article_id": req.ArticleID,
		"score":      score,
		"correct":    correct,
		"total":      total,
		"results":    results,
	}, nil
}
