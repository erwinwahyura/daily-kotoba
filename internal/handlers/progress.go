package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type ProgressHandler struct {
	vocabService *services.VocabService
}

func NewProgressHandler(vocabService *services.VocabService) *ProgressHandler {
	return &ProgressHandler{vocabService: vocabService}
}

func (h *ProgressHandler) GetProgress(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	stats, err := h.vocabService.GetProgress(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get progress", err)
		return
	}

	utils.SendSuccess(c, 200, "Progress retrieved successfully", gin.H{"progress": stats})
}

func (h *ProgressHandler) GetStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	stats, err := h.vocabService.GetProgress(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get statistics", err)
		return
	}

	utils.SendSuccess(c, 200, "Statistics retrieved successfully", gin.H{"stats": stats})
}
