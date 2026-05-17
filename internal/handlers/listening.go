package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type ListeningHandler struct {
	listeningService *services.ListeningService
}

func NewListeningHandler(listeningService *services.ListeningService) *ListeningHandler {
	return &ListeningHandler{listeningService: listeningService}
}

// GET /listening/exercises/:level
func (h *ListeningHandler) GetExercises(c *gin.Context) {
	level := c.Param("level")
	exercises, err := h.listeningService.GetExercises(level)
	if err != nil {
		utils.SendError(c, 500, "Failed to load exercises", err)
		return
	}
	utils.SendSuccess(c, 200, "", gin.H{"exercises": exercises})
}

// GET /listening/exercise/:id
func (h *ListeningHandler) GetExercise(c *gin.Context) {
	id := c.Param("id")
	detail, err := h.listeningService.GetExerciseDetail(id)
	if err != nil {
		utils.SendError(c, 404, "Exercise not found", err)
		return
	}
	utils.SendSuccess(c, 200, "", detail)
}

// POST /listening/session/start
func (h *ListeningHandler) StartSession(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	var req struct {
		ExerciseID string `json:"exercise_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "exercise_id required", err)
		return
	}
	session, err := h.listeningService.StartSession(userID, req.ExerciseID)
	if err != nil {
		utils.SendError(c, 500, "Failed to start session", err)
		return
	}
	utils.SendSuccess(c, 201, "Session started", session)
}

// POST /listening/session/answer
func (h *ListeningHandler) SubmitAnswer(c *gin.Context) {
	var req models.SubmitListeningAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}
	result, err := h.listeningService.SubmitAnswer(req.SessionID, req.QuestionID, req.Answer)
	if err != nil {
		utils.SendError(c, 400, err.Error(), err)
		return
	}
	utils.SendSuccess(c, 200, "", result)
}
