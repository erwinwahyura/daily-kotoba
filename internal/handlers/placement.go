package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type PlacementHandler struct {
	placementService *services.PlacementService
}

func NewPlacementHandler(placementService *services.PlacementService) *PlacementHandler {
	return &PlacementHandler{
		placementService: placementService,
	}
}

// GetPlacementTest handles GET /api/placement-test
// Returns all placement test questions (without correct answers)
func (h *PlacementHandler) GetPlacementTest(c *gin.Context) {
	questions, err := h.placementService.GetPlacementTest()
	if err != nil {
		utils.SendError(c, 500, "Failed to retrieve placement test", err)
		return
	}

	utils.SendSuccess(c, 200, "Placement test retrieved successfully", gin.H{
		"questions": questions,
	})
}

// SubmitPlacementTest handles POST /api/placement-test/submit
// Evaluates the test and assigns a JLPT level to the user
func (h *PlacementHandler) SubmitPlacementTest(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}

	// Parse request body
	var submission models.PlacementTestSubmission
	if err := c.ShouldBindJSON(&submission); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	// Validate that answers were provided
	if len(submission.Answers) == 0 {
		utils.SendError(c, 400, "No answers provided", nil)
		return
	}

	// Submit test and get results
	result, err := h.placementService.SubmitPlacementTest(userID, &submission)
	if err != nil {
		utils.SendError(c, 500, "Failed to submit placement test", err)
		return
	}

	utils.SendSuccess(c, 200, "Placement test submitted successfully", result)
}

// GetUserTestResult handles GET /api/placement-test/result
// Returns the user's most recent placement test result
func (h *PlacementHandler) GetUserTestResult(c *gin.Context) {
	// Get user ID from context
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}

	result, err := h.placementService.GetUserTestResult(userID)
	if err != nil {
		utils.SendError(c, 500, "Failed to retrieve test result", err)
		return
	}

	if result == nil {
		utils.SendError(c, 404, "No placement test result found", nil)
		return
	}

	utils.SendSuccess(c, 200, "Placement test result retrieved successfully", result)
}
