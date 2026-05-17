package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/erwinwahyura/daily-kotoba/internal/config"
	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/handlers"
	"github.com/erwinwahyura/daily-kotoba/internal/middleware"
	"github.com/erwinwahyura/daily-kotoba/internal/repository"
	"github.com/erwinwahyura/daily-kotoba/internal/services"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	sqlDB, err := cfg.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	// Configure connection pool (only for Postgres)
	if cfg.DB.Driver == "postgres" {
		sqlDB.SetMaxOpenConns(25)
		sqlDB.SetMaxIdleConns(5)
	}

	// Test database connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Printf("Successfully connected to %s database", cfg.DB.Driver)

	// Wrap database with driver-aware abstraction
	wrappedDB := db.New(sqlDB, cfg.DB.Driver)
	
	// Initialize SQLite if needed
	if cfg.DB.Driver == "sqlite" {
		if err := wrappedDB.InitializeSQLite(); err != nil {
			log.Fatalf("Failed to initialize SQLite: %v", err)
		}
		log.Println("SQLite initialized with WAL mode")
	}

	// Run database migrations (idempotent - only applies pending migrations)
	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "./migrations" // Default to local migrations folder
	}
	
	if err := wrappedDB.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Database migrations completed")

	// Run auto-seeding (idempotent - only inserts data once)
	seedsDir := os.Getenv("SEEDS_DIR")
	if seedsDir == "" {
		seedsDir = "./seeds" // Default to local seeds folder
	}
	
	if err := wrappedDB.RunAutoSeeding(seedsDir); err != nil {
		log.Fatalf("Failed to run auto-seeding: %v", err)
	}
	log.Println("Auto-seeding completed (only applied new seeds)")

	// Initialize repositories
	userRepo := repository.NewUserRepository(wrappedDB)
	vocabRepo := repository.NewVocabRepository(wrappedDB)
	progressRepo := repository.NewProgressRepository(wrappedDB)
	placementRepo := repository.NewPlacementRepository(wrappedDB)
	grammarRepo := repository.NewGrammarRepository(wrappedDB)
	srsRepo := repository.NewSRSRepository(wrappedDB)
	conjRepo := repository.NewConjugationRepository(wrappedDB)
	ttsRepo := repository.NewTTSRepository(wrappedDB)
	jlptRepo := repository.NewJLPTRepository(wrappedDB)
	kanjiRepo := repository.NewKanjiRepository(wrappedDB)
	goalsRepo := repository.NewGoalsRepository(wrappedDB)
	listeningRepo := repository.NewListeningRepository(wrappedDB)
	conversationRepo := repository.NewConversationRepository(wrappedDB)

	// Seed static data (kanji, listening exercises, conversation scenarios)
	log.Println("Seeding static data...")
	if err := kanjiRepo.SeedSampleKanji(); err != nil {
		log.Printf("Warning: failed to seed kanji: %v", err)
	} else {
		log.Println("Kanji data seeded successfully")
	}
	if err := listeningRepo.SeedSampleExercises(); err != nil {
		log.Printf("Warning: failed to seed listening exercises: %v", err)
	} else {
		log.Println("Listening exercises seeded successfully")
	}
	if err := conversationRepo.SeedScenarios(); err != nil {
		log.Printf("Warning: failed to seed conversation scenarios: %v", err)
	} else {
		log.Println("Conversation scenarios seeded successfully")
	}

	// Initialize services
	authService := services.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpirationHours)
	vocabService := services.NewVocabService(vocabRepo, progressRepo, userRepo)
	placementService := services.NewPlacementService(placementRepo, userRepo)
	grammarService := services.NewGrammarService(grammarRepo, progressRepo, userRepo)
	srsService := services.NewSRSService(srsRepo, vocabRepo, grammarRepo, userRepo)
	conjService := services.NewConjugationService(conjRepo)
	ttsService := services.NewTTSService(ttsRepo)
	jlptService := services.NewJLPTService(jlptRepo)
	kanjiService := services.NewKanjiService(kanjiRepo)
	goalsService := services.NewGoalsService(goalsRepo)
	listeningService := services.NewListeningService(listeningRepo)
	conversationService := services.NewConversationService(conversationRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	ttsHandler := handlers.NewTTSHandler(ttsService)
	jlptHandler := handlers.NewJLPTHandler(jlptService)
	vocabHandler := handlers.NewVocabularyHandler(vocabService)
	progressHandler := handlers.NewProgressHandler(vocabService)
	placementHandler := handlers.NewPlacementHandler(placementService)
	grammarHandler := handlers.NewGrammarHandler(grammarService)
	srsHandler := handlers.NewSRShandler(srsService)
	conjHandler := handlers.NewConjugationHandler(conjService)
	kanjiHandler := handlers.NewKanjiHandler(kanjiService)
	goalsHandler := handlers.NewGoalsHandler(goalsService)
	listeningHandler := handlers.NewListeningHandler(listeningService)
	conversationHandler := handlers.NewConversationHandler(conversationService)

	// Set up Gin router
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	router.Use(middleware.CORSMiddleware())

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

		// Placement test routes
		placement := v1.Group("/placement-test")
		{
			placement.GET("", placementHandler.GetPlacementTest) // Public: get questions
			placement.POST("/submit", middleware.AuthMiddleware(authService), placementHandler.SubmitPlacementTest)
			placement.GET("/result", middleware.AuthMiddleware(authService), placementHandler.GetUserTestResult)
		}

		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(authService))
		{
			// Vocabulary routes
			vocab := protected.Group("/vocab")
			{
				vocab.GET("/daily", vocabHandler.GetDailyWord)
				vocab.GET("/level/:level", vocabHandler.GetVocabularyByLevel)
				vocab.GET("/search", vocabHandler.SearchVocab)
				vocab.GET("/:id", vocabHandler.GetVocabByID)
				vocab.POST("/:id/skip", vocabHandler.SkipWord)
			}

			// Progress routes
			progress := protected.Group("/progress")
			{
				progress.GET("", progressHandler.GetProgress)
				progress.GET("/stats", progressHandler.GetStats)
			}

			// Grammar pattern routes
			grammar := protected.Group("/grammar")
			{
				grammar.GET("/daily", grammarHandler.GetDailyPattern)
				grammar.GET("/level/:level", grammarHandler.GetPatternsByLevel)
				grammar.GET("/search", grammarHandler.SearchGrammar)
				grammar.GET("/compare/pairs", grammarHandler.GetComparisonPairs)
				grammar.GET("/compare/detail", grammarHandler.ComparePatterns)
				grammar.GET("/:id", grammarHandler.GetPatternByID)
				grammar.POST("/:id/skip", grammarHandler.SkipPattern)
			}

			// SRS (Spaced Repetition) routes
			srs := protected.Group("/srs")
			{
				srs.GET("/queue", srsHandler.GetReviewQueue)
				srs.POST("/review", srsHandler.SubmitReview)
				srs.GET("/stats", srsHandler.GetSRSStats)
				srs.POST("/init", srsHandler.InitializeItem)
			}

			// Conjugation Drill routes
			conjugation := protected.Group("/conjugation")
			{
				conjugation.GET("/start", conjHandler.StartSession)
				conjugation.POST("/answer", conjHandler.SubmitAnswer)
				conjugation.GET("/progress", conjHandler.GetProgress)
				conjugation.GET("/weak-points", conjHandler.GetWeakPoints)
				conjugation.POST("/weak-points/drill", conjHandler.StartWeakPointDrill)
			}

			// TTS (Text-to-Speech) routes
			tts := protected.Group("/tts")
			{
				tts.POST("/generate", ttsHandler.GenerateTTS)
				tts.GET("/voices", ttsHandler.GetVoices)
				tts.GET("/stats", ttsHandler.GetCacheStats)
			}
			v1.GET("/tts/audio/:id", ttsHandler.GetAudio)

			// JLPT Mock Test routes
			jlpt := protected.Group("/jlpt")
			{
				jlpt.GET("/levels", jlptHandler.GetLevels)
				jlpt.GET("/tests/:level", jlptHandler.GetTests)
				jlpt.POST("/start", jlptHandler.StartTest)
				jlpt.POST("/answer", jlptHandler.SubmitAnswer)
				jlpt.POST("/complete/:session_id", jlptHandler.CompleteTest)
				jlpt.GET("/progress/:session_id", jlptHandler.GetProgress)
				jlpt.GET("/history", jlptHandler.GetHistory)
			}

			// Goals, streaks & achievements
			goals := protected.Group("/goals")
			{
				goals.GET("/daily", goalsHandler.GetDailyProgress)
				goals.POST("/progress", goalsHandler.UpdateProgress)
				goals.GET("/settings", goalsHandler.GetSettings)
				goals.PUT("/settings", goalsHandler.UpdateSettings)
				goals.GET("/streak", goalsHandler.GetStreak)
				goals.GET("/achievements", goalsHandler.GetAchievements)
				goals.GET("/achievements/all", goalsHandler.GetAllAchievements)
				goals.GET("/weekly", goalsHandler.GetWeeklyProgress)
			}

			// Kanji writing practice
			kanji := protected.Group("/kanji")
			{
				kanji.GET("/level/:level", kanjiHandler.GetKanjiByLevel)
				kanji.GET("/character/:char", kanjiHandler.GetKanjiByCharacter)
				kanji.POST("/practice/start", kanjiHandler.StartPracticeSession)
				kanji.POST("/practice/compare", kanjiHandler.CompareStroke)
				kanji.GET("/practice/:id", kanjiHandler.GetPracticeSession)
				kanji.GET("/stats", kanjiHandler.GetUserStats)
			}
			protected.POST("/kanji/seed", kanjiHandler.SeedKanjiData)

			// Listening practice
			listening := protected.Group("/listening")
			{
				listening.GET("/exercises/:level", listeningHandler.GetExercises)
				listening.GET("/exercise/:id", listeningHandler.GetExercise)
				listening.POST("/session/start", listeningHandler.StartSession)
				listening.POST("/session/answer", listeningHandler.SubmitAnswer)
				listening.GET("/session/:id", listeningHandler.GetSession)
				listening.GET("/stats", listeningHandler.GetStats)
			}
			protected.POST("/listening/seed", listeningHandler.SeedExercises)

			// Nichijou Conversation routes
			nichijou := protected.Group("/nichijou")
			{
				nichijou.GET("/scenarios", conversationHandler.GetScenarios)
				nichijou.POST("/chat/start", conversationHandler.StartChat)
				nichijou.POST("/chat/message", conversationHandler.SendMessage)
				nichijou.POST("/chat/end/:id", conversationHandler.EndSession)
				nichijou.GET("/chat/history/:id", conversationHandler.GetSessionHistory)
				nichijou.GET("/stats", conversationHandler.GetUserStats)
			}
		}
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s (environment: %s)", addr, cfg.Server.Env)
	log.Printf("Database: %s", cfg.DB.Driver)

	// Run server in goroutine
	go func() {
		if err := router.Run(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// For SQLite: periodic WAL checkpoint to prevent unbounded growth
	if cfg.DB.Driver == "sqlite" {
		go func() {
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()
			for range ticker.C {
				if err := wrappedDB.CheckpointWAL("passive"); err != nil {
					log.Printf("WAL checkpoint error: %v", err)
				}
			}
		}()
		log.Println("Started periodic WAL checkpointing (every 5min)")
	}

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// SQLite optimization before shutdown
	if cfg.DB.Driver == "sqlite" {
		log.Println("Running SQLite optimization...")
		if err := wrappedDB.Optimize(); err != nil {
			log.Printf("Optimization error: %v", err)
		}
		// Final checkpoint
		if err := wrappedDB.CheckpointWAL("full"); err != nil {
			log.Printf("Final checkpoint error: %v", err)
		}
	}

	// Close database
	if err := sqlDB.Close(); err != nil {
		log.Printf("Database close error: %v", err)
	}

	log.Println("Server stopped")
}