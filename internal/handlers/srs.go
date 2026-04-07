package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type SRShandler struct {
	srsService *services.SRSService
}

func NewSRShandler(srsService *services.SRSService) *SRShandler {
	return &SRShandler{srsService: srsService}
}

// GetReviewQueue returns items due for SRS review
func (h *SRShandler) GetReviewQueue(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	limit := 20
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 && n <= 100 {
			limit = n
		}
	}

	queue, err := h.srsService.GetReviewQueue(userID, limit)
	if err != nil {
		utils.SendError(c, 500, "Failed to get review queue", err)
		return
	}

	utils.SendSuccess(c, 200, "Review queue retrieved", queue)
}

// SubmitReview processes a review and updates SRS
func (h *SRShandler) SubmitReview(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	var req models.SRSReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	result, err := h.srsService.SubmitReview(userID, &req)
	if err != nil {
		utils.SendError(c, 500, "Failed to process review", err)
		return
	}

	utils.SendSuccess(c, 200, "Review processed", result)
}

// GetSRSStats returns SRS statistics
func (h *SRShandler) GetSRSStats(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	stats, err := h.srsService.GetStats(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to get SRS stats", err)
		return
	}

	utils.SendSuccess(c, 200, "SRS statistics", stats)
}

// InitializeItem adds an item to SRS tracking when first learned
func (h *SRShandler) InitializeItem(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	var req struct {
		ItemID   string `json:"item_id" binding:"required"`
		ItemType string `json:"item_type" binding:"required,oneof=vocabulary grammar"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request", err)
		return
	}

	if err := h.srsService.InitializeItem(userID, req.ItemID, req.ItemType); err != nil {
		utils.SendError(c, 500, "Failed to initialize item", err)
		return
	}

	utils.SendSuccess(c, 201, "Item added to SRS", nil)
}