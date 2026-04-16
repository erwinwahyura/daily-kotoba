package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
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

// SkipPattern advances to the next grammar pattern
func (h *GrammarHandler) SkipPattern(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	patternID := c.Param("id")
	if patternID == "" {
		utils.SendError(c, 400, "Pattern ID is required", nil)
		return
	}

	var req struct {
		Status string `json:"status" binding:"required,oneof=studied skipped"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	nextPattern, err := h.grammarService.SkipToNextPattern(userID, patternID, req.Status)
	if err != nil {
		utils.SendError(c, 500, "Failed to skip to next pattern", err)
		return
	}

	utils.SendSuccess(c, 200, "Moved to next pattern successfully", nextPattern)
}

// GetComparisonPairs returns grammar patterns grouped for side-by-side comparison
func (h *GrammarHandler) GetComparisonPairs(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	level := c.DefaultQuery("level", "N4")
	pairs, err := h.grammarService.GetComparisonPairs(userID, level)
	if err != nil {
		utils.SendError(c, 500, "Failed to get comparison pairs", err)
		return
	}

	utils.SendSuccess(c, 200, "Comparison pairs retrieved successfully", gin.H{
		"pairs": pairs,
		"level": level,
	})
}

// ComparePatterns returns detailed comparison between two specific patterns
func (h *GrammarHandler) ComparePatterns(c *gin.Context) {
	patternA := c.Query("a")
	patternB := c.Query("b")

	if patternA == "" || patternB == "" {
		utils.SendError(c, 400, "Both pattern IDs required (query params: a, b)", nil)
		return
	}

	comparison, err := h.grammarService.ComparePatterns(patternA, patternB)
	if err != nil {
		utils.SendError(c, 500, "Failed to compare patterns", err)
		return
	}

	utils.SendSuccess(c, 200, "Comparison retrieved successfully", comparison)
}

// SearchGrammar searches grammar patterns by query
func (h *GrammarHandler) SearchGrammar(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		utils.SendError(c, 400, "Search query is required", nil)
		return
	}

	level := c.Query("level")
	
	results, err := h.grammarService.SearchGrammar(query, level)
	if err != nil {
		utils.SendError(c, 500, "Failed to search grammar patterns", err)
		return
	}

	utils.SendSuccess(c, 200, "Search completed", gin.H{"results": results, "count": len(results)})
}
