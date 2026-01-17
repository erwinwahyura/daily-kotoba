package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/services"
	"github.com/yourusername/kotoba-api/internal/utils"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	user, token, err := h.authService.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user with this email already exists" {
			utils.SendError(c, 409, err.Error(), nil)
			return
		}
		utils.SendError(c, 500, "Failed to register user", err)
		return
	}

	response := models.AuthResponse{
		User:  user,
		Token: token,
	}

	utils.SendSuccess(c, 201, "User registered successfully", response)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendError(c, 400, "Invalid request body", err)
		return
	}

	user, token, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		utils.SendError(c, 401, err.Error(), nil)
		return
	}

	response := models.AuthResponse{
		User:  user,
		Token: token,
	}

	utils.SendSuccess(c, 200, "Login successful", response)
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		utils.SendError(c, 401, "User not authenticated", nil)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		utils.SendError(c, 404, "User not found", err)
		return
	}

	utils.SendSuccess(c, 200, "User retrieved successfully", gin.H{"user": user})
}
