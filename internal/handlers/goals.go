package handlers

import (
	"net/http"

	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
	"github.com/gin-gonic/gin"
)

// GoalsHandler handles goals, streaks, and achievements HTTP requests
type GoalsHandler struct {
	service *services.GoalsService
}

// NewGoalsHandler creates a new handler
func NewGoalsHandler(service *services.GoalsService) *GoalsHandler {
	return &GoalsHandler{
		service: service,
	}
}

// GetDailyProgress returns today's progress summary
func (h *GoalsHandler) GetDailyProgress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	progress, err := h.service.GetDailyProgress(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get daily progress", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Daily progress retrieved", progress)
}

// UpdateProgressRequest represents a progress update request
type UpdateProgressRequest struct {
	ActivityType string `json:"activity_type" binding:"required,oneof=vocab grammar kanji conjugation reading"`
	Count        int    `json:"count" binding:"required,min=1"`
}

// UpdateProgress updates progress for an activity
func (h *GoalsHandler) UpdateProgress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	if err := h.service.UpdateProgress(userID, req.ActivityType, req.Count); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update progress", err)
		return
	}

	// Check for new achievements
	newAchievements, err := h.service.CheckAndAwardAchievements(userID)
	if err != nil {
		// Log but don't fail the request
		newAchievements = nil
	}

	utils.SendSuccess(c, http.StatusOK, "Progress updated", gin.H{
		"new_achievements": newAchievements,
	})
}

// GetSettings returns user's goal settings
func (h *GoalsHandler) GetSettings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	settings, err := h.service.GetGoalSettings(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get settings", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Settings retrieved", settings)
}

// UpdateSettingsRequest represents settings update request
type UpdateSettingsRequest struct {
	VocabTarget       int    `json:"vocab_target" binding:"min=0,max=100"`
	GrammarTarget     int    `json:"grammar_target" binding:"min=0,max=50"`
	KanjiTarget       int    `json:"kanji_target" binding:"min=0,max=50"`
	ConjugationTarget int    `json:"conjugation_target" binding:"min=0,max=100"`
	ReadingTarget     int    `json:"reading_target" binding:"min=0,max=10"`
	EnableReminders   bool   `json:"enable_reminders"`
	ReminderTime      string `json:"reminder_time"`
}

// UpdateSettings updates user's goal settings
func (h *GoalsHandler) UpdateSettings(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	settings := &models.GoalSettings{
		VocabTarget:       req.VocabTarget,
		GrammarTarget:     req.GrammarTarget,
		KanjiTarget:       req.KanjiTarget,
		ConjugationTarget: req.ConjugationTarget,
		ReadingTarget:     req.ReadingTarget,
		EnableReminders:   req.EnableReminders,
		ReminderTime:      req.ReminderTime,
	}

	if err := h.service.UpdateGoalSettings(userID, settings); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to update settings", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Settings updated", settings)
}

// GetStreak returns user's streak information
func (h *GoalsHandler) GetStreak(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	streak, err := h.service.GetUserStreak(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get streak", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Streak retrieved", streak)
}

// GetAchievements returns user's earned achievements
func (h *GoalsHandler) GetAchievements(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	achievements, err := h.service.GetAchievements(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get achievements", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Achievements retrieved", achievements)
}

// GetAllAchievements returns all available achievement definitions
func (h *GoalsHandler) GetAllAchievements(c *gin.Context) {
	definitions := h.service.GetAllAchievementDefinitions()
	utils.SendSuccess(c, http.StatusOK, "Achievement definitions retrieved", definitions)
}

// GetWeeklyProgress returns progress for the last 7 days
func (h *GoalsHandler) GetWeeklyProgress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	progress, err := h.service.GetWeeklyProgress(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to get weekly progress", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Weekly progress retrieved", progress)
}
