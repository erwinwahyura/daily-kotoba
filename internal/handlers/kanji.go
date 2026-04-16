package handlers

import (
	"net/http"

	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
	"github.com/erwinwahyura/daily-kotoba/internal/utils"
	"github.com/gin-gonic/gin"
)

// KanjiHandler handles kanji writing practice HTTP requests
type KanjiHandler struct {
	service *services.KanjiService
}

// NewKanjiHandler creates a new handler
func NewKanjiHandler(service *services.KanjiService) *KanjiHandler {
	return &KanjiHandler{
		service: service,
	}
}

// GetKanjiByCharacter returns kanji details including stroke order
func (h *KanjiHandler) GetKanjiByCharacter(c *gin.Context) {
	char := c.Param("char")
	if char == "" {
		utils.SendError(c, http.StatusBadRequest, "Character is required", nil)
		return
	}

	kanji, err := h.service.GetKanjiByCharacter(char)
	if err != nil {
		utils.SendError(c, http.StatusNotFound, "Kanji not found", err)
		return
	}

	// Don't expose full stroke paths to prevent cheating
	// Only send stroke count and basic info
	simplified := gin.H{
		"character":    kanji.Character,
		"jlpt_level":   kanji.JLPTLevel,
		"meaning":      kanji.Meaning,
		"readings":     kanji.Readings,
		"stroke_count": kanji.StrokeCount,
	}

	utils.SendSuccess(c, http.StatusOK, "Kanji retrieved successfully", simplified)
}

// GetKanjiByLevel returns kanji list for a JLPT level
func (h *KanjiHandler) GetKanjiByLevel(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		level = "N5"
	}

	limit := 50
	if l := c.Query("limit"); l != "" {
		// Parse limit if provided
		// For now, use default
	}

	response, err := h.service.GetKanjiByLevel(level, limit)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve kanji", err)
		return
	}

	// Simplify response - don't expose stroke paths
	var simplified []gin.H
	for _, k := range response.Kanji {
		simplified = append(simplified, gin.H{
			"character":    k.Character,
			"jlpt_level":   k.JLPTLevel,
			"meaning":      k.Meaning,
			"readings":     k.Readings,
			"stroke_count": k.StrokeCount,
		})
	}

	utils.SendSuccess(c, http.StatusOK, "Kanji list retrieved successfully", gin.H{
		"kanji":       simplified,
		"total_count": response.TotalCount,
		"level":       response.Level,
	})
}

// StartPracticeRequest represents a practice session start request
type StartPracticeRequest struct {
	KanjiChar string `json:"kanji_char" binding:"required"`
}

// StartPracticeSession creates a new kanji writing practice session
func (h *KanjiHandler) StartPracticeSession(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	var req StartPracticeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	session, err := h.service.StartPracticeSession(userID, req.KanjiChar)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to start practice session", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Practice session started", gin.H{
		"session_id":   session.ID,
		"kanji_char":   session.KanjiChar,
		"stroke_count": len(session.UserStrokes),
		"status":       session.Status,
	})
}

// CompareStrokeRequest represents a stroke comparison request
type CompareStrokeRequest struct {
	SessionID string           `json:"session_id" binding:"required"`
	StrokeNum int              `json:"stroke_num" binding:"required,min=1"`
	Path      []models.Point   `json:"path" binding:"required"`
}

// CompareStroke compares user's stroke with reference
func (h *KanjiHandler) CompareStroke(c *gin.Context) {
	var req CompareStrokeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, http.StatusBadRequest, "Invalid request", err)
		return
	}

	result, err := h.service.CompareStroke(req.SessionID, req.StrokeNum, req.Path)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to compare stroke", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Stroke compared", result)
}

// GetPracticeSession retrieves practice session details
func (h *KanjiHandler) GetPracticeSession(c *gin.Context) {
	sessionID := c.Param("id")
	if sessionID == "" {
		utils.SendError(c, http.StatusBadRequest, "Session ID is required", nil)
		return
	}

	session, err := h.service.GetPracticeSession(sessionID)
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

	utils.SendSuccess(c, http.StatusOK, "Session retrieved", gin.H{
		"session_id":   session.ID,
		"kanji_char":   session.KanjiChar,
		"strokes":      len(session.UserStrokes),
		"accuracy":     session.Accuracy,
		"status":       session.Status,
		"completed_at": session.CompletedAt,
	})
}

// GetUserStats returns user's kanji practice statistics
func (h *KanjiHandler) GetUserStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, http.StatusUnauthorized, "User not authenticated", nil)
		return
	}

	stats, err := h.service.GetUserStats(userID)
	if err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to retrieve stats", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Stats retrieved", stats)
}

// SeedKanjiData seeds sample kanji data (admin only)
func (h *KanjiHandler) SeedKanjiData(c *gin.Context) {
	if err := h.service.SeedKanjiData(); err != nil {
		utils.SendError(c, http.StatusInternalServerError, "Failed to seed kanji data", err)
		return
	}

	utils.SendSuccess(c, http.StatusOK, "Kanji data seeded successfully", nil)
}
