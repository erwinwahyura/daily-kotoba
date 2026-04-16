package handlers

import (
	"net/http"

	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
	"github.com/gin-gonic/gin"
)

// ListeningHandler handles listening practice HTTP requests
type ListeningHandler struct {
	service *services.ListeningService
}

// NewListeningHandler creates a new handler
func NewListeningHandler(service *services.ListeningService) *ListeningHandler {
	return &ListeningHandler{
		service: service,
	}
}

// GetExercises returns listening exercises for a JLPT level
func (h *ListeningHandler) GetExercises(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		level = "N5"
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		// Could parse limit here
	}

	response, err := h.service.GetExercisesByLevel(level, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get exercises", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Exercises retrieved", response)
}

// GetExercise returns a specific listening exercise
func (h *ListeningHandler) GetExercise(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		utils.SendError(c, http.StatusBadRequest, "Exercise ID is required", nil)
		return
	}

	exercise, err := h.service.GetExercise(id)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Exercise not found", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Exercise retrieved", exercise)
}

// StartSessionRequest represents a session start request
type StartSessionRequest struct {
	ExerciseID string `json:"exercise_id" binding:"required"`
}

// StartSession creates a new listening session
func (h *ListeningHandler) StartSession(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req StartSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	session, err := h.service.StartSession(userID, req.ExerciseID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to start session", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Session started", gin.H{
		"session_id": session.ID,
		"exercise_id": session.ExerciseID,
		"status": session.Status,
	})
}

// SubmitAnswerRequest represents an answer submission
type SubmitAnswerRequest struct {
	SessionID     string `json:"session_id" binding:"required"`
	QuestionID    string `json:"question_id" binding:"required"`
	Answer        int    `json:"answer" binding:"required,min=0"`
	AudioPosition int    `json:"audio_position"`
}

// SubmitAnswer processes a user's answer
func (h *ListeningHandler) SubmitAnswer(c *gin.Context) {
	var req SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	result, err := h.service.SubmitAnswer(req.SessionID, req.QuestionID, req.Answer, req.AudioPosition)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to submit answer", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Answer submitted", gin.H{
		"correct": result.IsCorrect,
		"answer":  result.Answer,
	})
}

// GetSession retrieves session progress
func (h *ListeningHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		utils.SendError(c, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, err := h.service.GetSession(sessionID)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Session not found", err)
		return
	}

	// Verify user owns this session
	userID, exists := middleware.GetUserID(c)
	if !exists || session.UserID != userID {
		utils.SendError(c, http.StatusForbidden, "Access denied", nil)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Session retrieved", session)
}

// GetStats returns user's listening statistics
func (h *ListeningHandler) GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	stats, err := h.service.GetUserStats(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get stats", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Stats retrieved", stats)
}

// SeedExercises seeds sample listening exercises (admin only)
func (h *ListeningHandler) SeedExercises(c *gin.Context) {
	if err := h.service.SeedExercises(); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to seed exercises", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Exercises seeded successfully", nil)
}
