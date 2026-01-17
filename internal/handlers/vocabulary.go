package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type VocabularyHandler struct {
	vocabService *services.VocabService
}

func NewVocabularyHandler(vocabService *services.VocabService) *VocabularyHandler {
	return &VocabularyHandler{vocabService: vocabService}
}

func (h *VocabularyHandler) GetDailyWord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	vocab, err := h.vocabService.GetDailyWord(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get daily word", err)
		return
	}

	utils.SendSuccess(c, 200, "Daily word retrieved successfully", vocab)
}

func (h *VocabularyHandler) GetVocabByID(c *gin.Context) {
	vocabID := c.Param("id")
	if vocabID == "" {
		utils.SendError(c, 400, "Vocabulary ID is required", nil)
		return
	}

	vocab, err := h.vocabService.GetVocabByID(vocabID)
	if err != nil {
		utils.SendError(c, 404, "Vocabulary not found", err)
		return
	}

	utils.SendSuccess(c, 200, "Vocabulary retrieved successfully", gin.H{"vocabulary": vocab})
}

func (h *VocabularyHandler) SkipWord(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	vocabID := c.Param("id")
	if vocabID == "" {
		utils.SendError(c, 400, "Vocabulary ID is required", nil)
		return
	}

	var req models.SkipRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	nextVocab, err := h.vocabService.SkipToNextWord(userID, vocabID, req.Status)
	if err != nil {
		utils.SendError(c, 500, "Failed to skip to next word", err)
		return
	}

	utils.SendSuccess(c, 200, "Moved to next word successfully", nextVocab)
}

func (h *VocabularyHandler) GetVocabularyByLevel(c *gin.Context) {
	level := c.Param("level")
	if level == "" {
		utils.SendError(c, 400, "Level is required", nil)
		return
	}

	// Validate level
	validLevels := map[string]bool{"N5": true, "N4": true, "N3": true, "N2": true, "N1": true}
	if !validLevels[level] {
		utils.SendError(c, 400, "Invalid JLPT level", nil)
		return
	}

	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	response, err := h.vocabService.GetVocabularyByLevel(level, page, limit)
	if err != nil {
		utils.SendError(c, 500, "Failed to get vocabulary list", err)
		return
	}

	utils.SendSuccess(c, 200, "Vocabulary list retrieved successfully", response)
}
