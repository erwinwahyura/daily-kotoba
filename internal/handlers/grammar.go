package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type GrammarHandler struct {
	grammarService *services.GrammarService
}

func NewGrammarHandler(grammarService *services.GrammarService) *GrammarHandler {
	return &GrammarHandler{grammarService: grammarService}
}

// GetDailyPattern returns the current grammar pattern for the user
func (h *GrammarHandler) GetDailyPattern(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	response, err := h.grammarService.GetDailyPattern(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get grammar pattern", err)
		return
	}

	utils.SendSuccess(c, 200, "Grammar pattern retrieved successfully", response)
}

// GetPatternByID returns a specific grammar pattern
func (h *GrammarHandler) GetPatternByID(c *gin.Context) {
	patternID := c.Param("id")
	if patternID == "" {
		utils.SendError(c, 400, "Pattern ID is required", nil)
		return
	}

	pattern, err := h.grammarService.GetPatternByID(patternID)
	if err != nil {
		utils.SendError(c, 404, "Grammar pattern not found", err)
		return
	}

	utils.SendSuccess(c, 200, "Grammar pattern retrieved successfully", gin.H{"pattern": pattern})
}

// GetPatternsByLevel returns grammar patterns for a specific JLPT level
func (h *GrammarHandler) GetPatternsByLevel(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		utils.SendError(c, 400, "Level is required", nil)
		return
	}

	validLevels := map[string]bool{"N5": true, "N4": true, "N3": true, "N2": true, "N1": true}
	if !validLevels[level] {
		utils.SendError(c, 400, "Invalid JLPT level", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.grammarService.GetPatternsByLevel(level, page, limit)
	if err != nil {
		utils.SendError(c, 500, "Failed to get grammar patterns", err)
		return
	}

	utils.SendSuccess(c, 200, "Grammar patterns retrieved successfully", response)
}
