package handlers

import (
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type KanjiHandler struct {
	kanjiService *services.KanjiService
}

func NewKanjiHandler(kanjiService *services.KanjiService) *KanjiHandler {
	return &KanjiHandler{kanjiService: kanjiService}
}

// GET /kanji/character/:kanji
func (h *KanjiHandler) GetCharacter(c *gin.Context) {
	encoded := c.Param("kanji")
	char, err := url.QueryUnescape(encoded)
	if err != nil || char == "" {
		utils.SendError(c, 400, "Invalid kanji character", nil)
		return
	}
	k, err := h.kanjiService.GetCharacter(char)
	if err != nil {
		utils.SendError(c, 404, "Kanji not found", err)
		return
	}
	utils.SendSuccess(c, 200, "", k)
}

// POST /kanji/practice/start
func (h *KanjiHandler) StartPractice(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.SendError(c, 401, "Unauthorized", nil)
		return
	}
	var req struct {
		KanjiChar string `json:"kanji_char" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "kanji_char required", err)
		return
	}
	session, err := h.kanjiService.StartSession(userID, req.KanjiChar)
	if err != nil {
		utils.SendError(c, 400, err.Error(), err)
		return
	}
	utils.SendSuccess(c, 201, "Practice session started", session)
}

// POST /kanji/practice/compare
func (h *KanjiHandler) CompareStroke(c *gin.Context) {
	var req models.KanjiCompareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}
	result, err := h.kanjiService.CompareStroke(req.SessionID, req.StrokeNum, req.UserPath)
	if err != nil {
		utils.SendError(c, 400, err.Error(), err)
		return
	}
	utils.SendSuccess(c, 200, "", result)
}
