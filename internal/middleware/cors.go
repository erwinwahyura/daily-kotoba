package middleware

import (
	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Allow specific origins
		allowedOrigins := []string{
			"https://erwinwahyura.github.io",
			"https://kotoba-web.erwinwahyura.workers.dev",
			"http://localhost:3000",
			"http://localhost:8080",
			"http://localhost:5173", // Vite dev server
		}

		origin := c.GetHeader("Origin")
		
		// Check if origin is allowed
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		// Default to wildcard if no specific origin matched (for public health checks)
		if c.Writer.Header().Get("Access-Control-Allow-Origin") == "" && origin == "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
