package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type JLPTHandler struct {
	jlptService *services.JLPTService
}

func NewJLPTHandler(jlptService *services.JLPTService) *JLPTHandler {
	return &JLPTHandler{jlptService: jlptService}
}

// GetLevels returns JLPT level information
func (h *JLPTHandler) GetLevels(c *gin.Context) {
	levels := h.jlptService.GetLevelInfo()
	utils.SendSuccess(c, 200, "JLPT levels retrieved", gin.H{"levels": levels})
}

// GetTests returns available tests for a level
func (h *JLPTHandler) GetTests(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		utils.SendError(c, 400, "Level is required", nil)
		return
	}

	tests, err := h.jlptService.GetAvailableTests(level)
	if err != nil {
		utils.SendError(c, 500, "Failed to get tests", err)
		return
	}

	utils.SendSuccess(c, 200, "Tests retrieved", gin.H{"tests": tests})
}

// StartTest begins a new test session
func (h *JLPTHandler) StartTest(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	var req models.StartTestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	session, questions, err := h.jlptService.StartTest(userID, req.Level, req.Section)
	if err != nil {
		utils.SendError(c, 500, "Failed to start test", err)
		return
	}

	utils.SendSuccess(c, 200, "Test started", gin.H{
		"session":   session,
		"questions": questions,
	})
}

// SubmitAnswer records an answer during a test
func (h *JLPTHandler) SubmitAnswer(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}
	_ = userID

	var req models.SubmitAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	if err := h.jlptService.SubmitAnswer(req.SessionID, req.QuestionID, req.AnswerIndex); err != nil {
		utils.SendError(c, 500, "Failed to submit answer", err)
		return
	}

	utils.SendSuccess(c, 200, "Answer recorded", nil)
}

// CompleteTest finishes a test and returns results
func (h *JLPTHandler) CompleteTest(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}
	_ = userID

	sessionID := c.Param("session_id")
	if sessionID == "" {
		utils.SendError(c, 400, "Session ID is required", nil)
		return
	}

	var req struct {
		Answers map[string]int `json:"answers"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		// Allow empty body - will use stored answers
		req.Answers = nil
	}

	result, err := h.jlptService.CompleteTest(sessionID, req.Answers)
	if err != nil {
		utils.SendError(c, 500, "Failed to complete test", err)
		return
	}

	utils.SendSuccess(c, 200, "Test completed", result)
}

// GetHistory returns user's test history
func (h *JLPTHandler) GetHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	history, err := h.jlptService.GetUserHistory(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get history", err)
		return
	}

	utils.SendSuccess(c, 200, "History retrieved", gin.H{"history": history})
}

// GetProgress returns current test progress
func (h *JLPTHandler) GetProgress(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		utils.SendError(c, 400, "Session ID is required", nil)
		return
	}

	answered, total, err := h.jlptService.GetTestProgress(sessionID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get progress", err)
		return
	}

	utils.SendSuccess(c, 200, "Progress retrieved", gin.H{
		"answered": answered,
		"total":    total,
		"remaining": total - answered,
	})
}