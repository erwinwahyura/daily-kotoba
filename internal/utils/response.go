package utils

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func SendError(c *gin.Context, statusCode int, message string, err error) {
	response := ErrorResponse{
		Code:    statusCode,
		Message: message,
	}
	if err != nil {
		response.Error = err.Error()
	}
	c.JSON(statusCode, response)
}

func SendSuccess(c *gin.Context, statusCode int, message string, data interface{}) {
	response := SuccessResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}
	c.JSON(statusCode, response)
}
