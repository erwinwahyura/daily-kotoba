package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/yourusername/kotoba-api/internal/config"
	"github.com/yourusername/kotoba-api/internal/models"
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

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to database successfully")

	// Initialize repositories and services
	vocabRepo := repository.NewVocabRepository(db)
	progressRepo := repository.NewProgressRepository(db)
	userRepo := repository.NewUserRepository(db)
	vocabService := services.NewVocabService(vocabRepo, progressRepo, userRepo)

	// Seed N4 vocabulary
	n4Vocab := []models.Vocabulary{
		{
			Word:                "諦める",
			Reading:             "あきらめる",
			ShortMeaning:        "to give up",
			DetailedExplanation: "Used when abandoning an effort or goal. Often implies acceptance of failure or impossibility. This word carries a sense of finality and resignation.",
			ExampleSentences: models.ExampleSentences{
				"試験に合格するのを諦めた - I gave up on passing the exam",
				"彼女は夢を諦めなかった - She didn't give up on her dream",
				"もう諦めたほうがいい - You should give up already",
			},
			UsageNotes:    "Commonly used in both casual and formal contexts. Often paired with のを when followed by a noun phrase.",
			JLPTLevel:     "N4",
			IndexPosition: 0,
		},
		{
			Word:                "相変わらず",
			Reading:             "あいかわらず",
			ShortMeaning:        "as usual, as always",
			DetailedExplanation: "An expression meaning that something or someone remains the same as before. Often used in greetings or when commenting on unchanged situations.",
			ExampleSentences: models.ExampleSentences{
				"相変わらず元気ですか - Are you well as always?",
				"彼は相変わらず忙しい - He's busy as usual",
				"この街は相変わらず静かだ - This town is quiet as always",
			},
			UsageNotes:    "Commonly used in daily conversation. Can be used positively or negatively depending on context.",
			JLPTLevel:     "N4",
			IndexPosition: 1,
		},
		{
			Word:                "合う",
			Reading:             "あう",
			ShortMeaning:        "to match, to fit",
			DetailedExplanation: "Indicates compatibility or suitability. Can refer to physical fit, matching of opinions, or harmonious relationships.",
			ExampleSentences: models.ExampleSentences{
				"この服は私に合う - This outfit suits me",
				"意見が合わない - Our opinions don't match",
				"サイズが合います - The size fits",
			},
			UsageNotes:    "Often paired with に particle to indicate what something matches or fits with.",
			JLPTLevel:     "N4",
			IndexPosition: 2,
		},
		{
			Word:                "青い",
			Reading:             "あおい",
			ShortMeaning:        "blue, green (for traffic lights)",
			DetailedExplanation: "An i-adjective meaning blue. Also used for green in certain contexts like traffic lights and young/unripe things.",
			ExampleSentences: models.ExampleSentences{
				"空が青い - The sky is blue",
				"信号が青になった - The light turned green",
				"青いりんご - Green apple",
			},
			UsageNotes:    "Note that 青 is used for both blue and green in Japanese, particularly for traffic lights and unripe fruit.",
			JLPTLevel:     "N4",
			IndexPosition: 3,
		},
		{
			Word:                "赤ちゃん",
			Reading:             "あかちゃん",
			ShortMeaning:        "baby",
			DetailedExplanation: "A common word for baby or infant. Used affectionately and is appropriate in both casual and formal contexts.",
			ExampleSentences: models.ExampleSentences{
				"赤ちゃんが泣いている - The baby is crying",
				"可愛い赤ちゃんですね - What a cute baby",
				"赤ちゃんが生まれた - A baby was born",
			},
			UsageNotes:    "Can be used for babies of any gender. Very common in everyday conversation.",
			JLPTLevel:     "N4",
			IndexPosition: 4,
		},
	}

	log.Printf("Seeding %d N4 vocabulary words...", len(n4Vocab))

	err = vocabService.BulkCreateVocabulary(n4Vocab)
	if err != nil {
		log.Fatalf("Failed to seed vocabulary: %v", err)
	}

	log.Println("Successfully seeded vocabulary data!")
	log.Printf("Total words inserted: %d", len(n4Vocab))
}
