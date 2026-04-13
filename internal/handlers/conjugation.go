package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
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

	response, err := h.service.StartDrillSession(userID, form)
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
