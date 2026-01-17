package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/yourusername/kotoba-api/internal/config"
	"github.com/yourusername/kotoba-api/internal/handlers"
	"github.com/yourusername/kotoba-api/internal/middleware"
	"github.com/yourusername/kotoba-api/internal/repository"
	"github.com/yourusername/kotoba-api/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := sql.Open("postgres", cfg.GetDatabaseDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	vocabRepo := repository.NewVocabRepository(db)
	progressRepo := repository.NewProgressRepository(db)

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	vocabService := services.NewVocabService(vocabRepo, progressRepo, userRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	vocabHandler := handlers.NewVocabularyHandler(vocabService)
	progressHandler := handlers.NewProgressHandler(vocabService)

	// Set up Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Kotoba API is running",
		})
	})

	// API v1 group
	v1 := router.Group("/api")
	{
		// Public routes
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(authService), authHandler.GetMe)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Vocabulary routes
			vocab := protected.Group("/vocab")
			{
				vocab.GET("/daily", vocabHandler.GetDailyWord)
				vocab.GET("/:id", vocabHandler.GetVocabByID)
				vocab.POST("/:id/skip", vocabHandler.SkipWord)
			}

			// Vocabulary by level route (outside /vocab group for cleaner URL)
			protected.GET("/vocab/level/:level", vocabHandler.GetVocabularyByLevel)

			// Progress routes
			progress := protected.Group("/progress")
			{
				progress.GET("", progressHandler.GetProgress)
				progress.GET("/stats", progressHandler.GetStats)
			}
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Server.Env)

	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
