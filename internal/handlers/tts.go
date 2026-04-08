package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type TTSHandler struct {
	ttsService *services.TTSService
}

func NewTTSHandler(ttsService *services.TTSService) *TTSHandler {
	return &TTSHandler{ttsService: ttsService}
}

// GenerateTTS generates audio for given text
func (h *TTSHandler) GenerateTTS(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}
	_ = userID // Track usage per user if needed later

	var req models.TTSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	response, err := h.ttsService.GenerateAudio(req.Text, req.VoiceID)
	if err != nil {
		utils.SendError(c, 500, "Failed to generate audio", err)
		return
	}

	utils.SendSuccess(c, 200, "Audio generated successfully", response)
}

// GetAudio retrieves audio file by ID
func (h *TTSHandler) GetAudio(c *gin.Context) {
	audioID := c.Param("id")
	if audioID == "" {
		utils.SendError(c, 400, "Audio ID required", nil)
		return
	}

	audioData, contentType, err := h.ttsService.GetAudioByID(audioID)
	if err != nil {
		utils.SendError(c, 404, "Audio not found", err)
		return
	}

	c.Data(http.StatusOK, contentType, audioData)
}

// GetVoices returns available TTS voices
func (h *TTSHandler) GetVoices(c *gin.Context) {
	voices := h.ttsService.GetAvailableVoices()
	utils.SendSuccess(c, 200, "Available voices retrieved", gin.H{
		"voices": voices,
	})
}

// GetCacheStats returns TTS cache statistics
func (h *TTSHandler) GetCacheStats(c *gin.Context) {
	stats, err := h.ttsService.GetCacheStats()
	if err != nil {
		utils.SendError(c, 500, "Failed to get cache stats", err)
		return
	}
	utils.SendSuccess(c, 200, "Cache statistics retrieved", stats)
}
