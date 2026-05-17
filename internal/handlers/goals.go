package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type GoalsHandler struct {
	goalsService *services.GoalsService
}

func NewGoalsHandler(goalsService *services.GoalsService) *GoalsHandler {
	return &GoalsHandler{goalsService: goalsService}
}

func (h *GoalsHandler) GetDailyProgress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	data, err := h.goalsService.GetDailyProgress(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to load daily progress", err)
		return
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) GetStreak(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	data, err := h.goalsService.GetStreak(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to load streak", err)
		return
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) GetUserAchievements(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	data, err := h.goalsService.GetUserAchievements(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to load achievements", err)
		return
	}
	if data == nil {
		data = []models.UserAchievement{}
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) GetAllAchievements(c *gin.Context) {
	data, err := h.goalsService.GetAllAchievements()
	if err != nil {
		utils.SendError(c, 500, "Failed to load achievements", err)
		return
	}
	if data == nil {
		data = []models.AchievementDef{}
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) GetWeeklyProgress(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	data, err := h.goalsService.GetWeeklyProgress(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to load weekly progress", err)
		return
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) GetSettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	data, err := h.goalsService.GetSettings(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to load settings", err)
		return
	}
	utils.SendSuccess(c, 200, "", data)
}

func (h *GoalsHandler) SaveSettings(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	var settings models.UserGoalSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}
	if err := h.goalsService.SaveSettings(userID, &settings); err != nil {
		utils.SendError(c, 500, "Failed to save settings", err)
		return
	}
	utils.SendSuccess(c, 200, "Settings saved", nil)
}

func (h *GoalsHandler) RecordActivity(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	var req models.RecordActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}
	if err := h.goalsService.RecordActivity(userID, req.ActivityType, req.Count); err != nil {
		utils.SendError(c, 500, "Failed to record activity", err)
		return
	}
	// Opportunistically check achievements
	go func() { _ = h.goalsService.CheckAndGrantAchievements(userID) }()
	utils.SendSuccess(c, 200, "Activity recorded", nil)
}
