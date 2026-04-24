package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
)

type ConjugationHandler struct {
	service *services.ConjugationService
}

func NewConjugationHandler(service *services.ConjugationService) *ConjugationHandler {
	return &ConjugationHandler{service: service}
}

// StartSession starts a new conjugation drill session
func (h *ConjugationHandler) StartSession(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	form := c.DefaultQuery("form", "te")
	maxLevel := c.DefaultQuery("max_level", "N5")

	response, err := h.service.StartDrillSessionWithLevel(userID, form, maxLevel)
	if err != nil {
		utils.SendError(c, 500, "Failed to start conjugation session", err)
		return
	}

	utils.SendSuccess(c, 200, "Conjugation session started", gin.H{
		"session":    response.Session,
		"challenges": response.Challenges,
		"progress":   response.Progress,
		"form_info":  response.FormInfo,
	})
}

// SubmitAnswer checks user's conjugation answer
func (h *ConjugationHandler) SubmitAnswer(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	var req struct {
		SessionID   string `json:"session_id" binding:"required"`
		ChallengeID string `json:"challenge_id" binding:"required"`
		Answer      string `json:"answer" binding:"required"`
		TimeSpentMs int    `json:"time_spent_ms"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	response, err := h.service.SubmitAnswer(userID, req.SessionID, req.ChallengeID, req.Answer, req.TimeSpentMs)
	if err != nil {
		utils.SendError(c, 500, "Failed to submit answer", err)
		return
	}

	utils.SendSuccess(c, 200, "Answer processed", response)
}

// GetProgress retrieves user's conjugation progress
func (h *ConjugationHandler) GetProgress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	progress, err := h.service.GetProgress(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get progress", err)
		return
	}

	utils.SendSuccess(c, 200, "Progress retrieved", progress)
}

// GetWeakPoints retrieves user's weak conjugation forms
func (h *ConjugationHandler) GetWeakPoints(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	analysis, err := h.service.GetWeakPointsAnalysis(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to analyze weak points", err)
		return
	}

	utils.SendSuccess(c, 200, "Weak points analysis", analysis)
}

// StartWeakPointDrill starts a focused drill for weak forms
func (h *ConjugationHandler) StartWeakPointDrill(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	response, err := h.service.GenerateWeakPointDrill(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to start weak point drill", err)
		return
	}

	utils.SendSuccess(c, 200, "Weak point drill started", gin.H{
		"session":    response.Session,
		"challenges": response.Challenges,
		"progress":   response.Progress,
		"form_info":  response.FormInfo,
	})
}
