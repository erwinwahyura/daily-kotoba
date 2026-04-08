package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handles cross-origin requests
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		
		// Check if origin should be allowed (erwarx.com subdomains + specific origins)
		allowOrigin := false
		if strings.HasSuffix(origin, ".erwarx.com") || origin == "https://erwarx.com" {
			allowOrigin = true
		} else {
			switch origin {
			case "https://erwinwahyura.github.io",
				"https://kotoba-web.erwinwahyura.workers.dev",
				"http://localhost:3000",
				"http://localhost:8080",
				"http://localhost:5173":
				allowOrigin = true
			}
		}
		
		if allowOrigin {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Vary", "Origin")
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
