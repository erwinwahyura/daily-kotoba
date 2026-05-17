package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
)

type ReadingHandler struct {
	service *services.ReadingService
}

func NewReadingHandler(service *services.ReadingService) *ReadingHandler {
	return &ReadingHandler{service: service}
}

// GetArticles returns articles for a JLPT level
func (h *ReadingHandler) GetArticles(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		utils.SendError(c, 400, "Level is required", nil)
		return
	}

	articles, err := h.service.GetArticlesByLevel(level)
	if err != nil {
		utils.SendError(c, 500, "Failed to get articles", err)
		return
	}

	utils.SendSuccess(c, 200, "Articles retrieved", gin.H{"articles": articles})
}

// GetArticle returns a single article with questions
func (h *ReadingHandler) GetArticle(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendError(c, 400, "Article ID is required", nil)
		return
	}

	article, err := h.service.GetArticle(id)
	if err != nil {
		utils.SendError(c, 404, "Article not found", err)
		return
	}

	utils.SendSuccess(c, 200, "Article retrieved", gin.H{"article": article})
}

// GetRandomArticle returns a random article for a level
func (h *ReadingHandler) GetRandomArticle(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		utils.SendError(c, 400, "Level is required", nil)
		return
	}

	article, err := h.service.GetRandomArticle(level)
	if err != nil {
		utils.SendError(c, 404, "No articles available", err)
		return
	}

	utils.SendSuccess(c, 200, "Article retrieved", gin.H{"article": article})
}

// SubmitAnswers grades reading comprehension answers
func (h *ReadingHandler) SubmitAnswers(c *gin.Context) {
	var req models.SubmitReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	result, err := h.service.SubmitAnswers(req)
	if err != nil {
		utils.SendError(c, 500, "Failed to grade answers", err)
		return
	}

	utils.SendSuccess(c, 200, "Answers graded", result)
}
