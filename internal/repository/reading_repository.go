package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

type ReadingRepository struct {
	db *db.DB
}

func NewReadingRepository(db *db.DB) *ReadingRepository {
	return &ReadingRepository{db: db}
}

func (r *ReadingRepository) GetArticlesByLevel(level string) ([]models.ReadingArticle, error) {
	query := `SELECT id, jlpt_level, title, text, topic, difficulty, word_count, created_at
		FROM reading_articles WHERE jlpt_level = ? ORDER BY difficulty, created_at`

	rows, err := r.db.Query(query, level)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var articles []models.ReadingArticle
	for rows.Next() {
		var a models.ReadingArticle
		if err := rows.Scan(&a.ID, &a.JLPTLevel, &a.Title, &a.Text, &a.Topic, &a.Difficulty, &a.WordCount, &a.CreatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}

func (r *ReadingRepository) GetArticleByID(id string) (*models.ReadingArticle, error) {
	query := `SELECT id, jlpt_level, title, text, topic, difficulty, word_count, created_at
		FROM reading_articles WHERE id = ?`

	var a models.ReadingArticle
	err := r.db.QueryRow(query, id).Scan(&a.ID, &a.JLPTLevel, &a.Title, &a.Text, &a.Topic, &a.Difficulty, &a.WordCount, &a.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}
	if err != nil {
		return nil, err
	}

	questions, err := r.GetQuestionsByArticle(id)
	if err != nil {
		return nil, err
	}
	a.Questions = questions
	return &a, nil
}

func (r *ReadingRepository) GetRandomArticle(level string) (*models.ReadingArticle, error) {
	articles, err := r.GetArticlesByLevel(level)
	if err != nil {
		return nil, err
	}
	if len(articles) == 0 {
		return nil, fmt.Errorf("no articles found for level %s", level)
	}

	picked := articles[rand.Intn(len(articles))]
	return r.GetArticleByID(picked.ID)
}

func (r *ReadingRepository) GetQuestionsByArticle(articleID string) ([]models.ReadingQuestion, error) {
	query := `SELECT id, article_id, question_num, question, options, correct_index, explanation
		FROM reading_questions WHERE article_id = ? ORDER BY question_num`

	rows, err := r.db.Query(query, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.ReadingQuestion
	for rows.Next() {
		var q models.ReadingQuestion
		var optionsJSON string
		if err := rows.Scan(&q.ID, &q.ArticleID, &q.QuestionNum, &q.Question, &optionsJSON, &q.CorrectIndex, &q.Explanation); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(optionsJSON), &q.Options); err != nil {
			q.Options = []string{}
		}
		questions = append(questions, q)
	}
	return questions, rows.Err()
}
