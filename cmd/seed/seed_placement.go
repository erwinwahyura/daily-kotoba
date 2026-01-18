package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/lib/pq"
	"github.com/yourusername/kotoba-api/internal/config"
)

type PlacementQuestion struct {
	QuestionText   string   `json:"question_text"`
	CorrectAnswer  string   `json:"correct_answer"`
	WrongAnswers   []string `json:"wrong_answers"`
	DifficultyLevel string   `json:"difficulty_level"`
	OrderIndex     int      `json:"order_index"`
}

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

	// Clear existing placement test questions
	log.Println("Clearing existing placement test questions...")
	_, err = db.Exec("DELETE FROM placement_questions")
	if err != nil {
		log.Fatalf("Failed to clear existing questions: %v", err)
	}
	log.Println("Existing placement test questions cleared")

	// Placement test questions (20 total)
	questions := []PlacementQuestion{
		// N5 Level Questions (5 questions - easiest)
		{
			QuestionText:   "Choose the correct reading for 食べる",
			CorrectAnswer:  "たべる",
			WrongAnswers:   []string{"のべる", "かべる", "よべる"},
			DifficultyLevel: "N5",
			OrderIndex:     1,
		},
		{
			QuestionText:   "What does 見る mean?",
			CorrectAnswer:  "to see, to watch",
			WrongAnswers:   []string{"to eat", "to go", "to sleep"},
			DifficultyLevel: "N5",
			OrderIndex:     2,
		},
		{
			QuestionText:   "Choose the correct particle: 私___学生です",
			CorrectAnswer:  "は",
			WrongAnswers:   []string{"が", "を", "に"},
			DifficultyLevel: "N5",
			OrderIndex:     3,
		},
		{
			QuestionText:   "What is the meaning of 大きい?",
			CorrectAnswer:  "big, large",
			WrongAnswers:   []string{"small", "old", "new"},
			DifficultyLevel: "N5",
			OrderIndex:     4,
		},
		{
			QuestionText:   "Choose the correct reading for 今日",
			CorrectAnswer:  "きょう",
			WrongAnswers:   []string{"こんにち", "いまひ", "きょうじつ"},
			DifficultyLevel: "N5",
			OrderIndex:     5,
		},

		// N4 Level Questions (5 questions)
		{
			QuestionText:   "Choose the correct meaning for 諦める",
			CorrectAnswer:  "to give up",
			WrongAnswers:   []string{"to celebrate", "to apologize", "to continue"},
			DifficultyLevel: "N4",
			OrderIndex:     6,
		},
		{
			QuestionText:   "What does 相変わらず mean?",
			CorrectAnswer:  "as usual, as always",
			WrongAnswers:   []string{"suddenly", "fortunately", "rarely"},
			DifficultyLevel: "N4",
			OrderIndex:     7,
		},
		{
			QuestionText:   "Choose the correct particle: 電車___遅れた",
			CorrectAnswer:  "が",
			WrongAnswers:   []string{"を", "に", "で"},
			DifficultyLevel: "N4",
			OrderIndex:     8,
		},
		{
			QuestionText:   "What is the meaning of 以上?",
			CorrectAnswer:  "more than, above",
			WrongAnswers:   []string{"less than", "exactly", "between"},
			DifficultyLevel: "N4",
			OrderIndex:     9,
		},
		{
			QuestionText:   "Choose the correct reading for 安全",
			CorrectAnswer:  "あんぜん",
			WrongAnswers:   []string{"あんせん", "やすぜん", "やすいぜん"},
			DifficultyLevel: "N4",
			OrderIndex:     10,
		},

		// N3 Level Questions (5 questions)
		{
			QuestionText:   "What does 思い出す mean?",
			CorrectAnswer:  "to recall, to remember",
			WrongAnswers:   []string{"to forget", "to imagine", "to regret"},
			DifficultyLevel: "N3",
			OrderIndex:     11,
		},
		{
			QuestionText:   "Choose the correct meaning for 恥ずかしい",
			CorrectAnswer:  "embarrassing, ashamed",
			WrongAnswers:   []string{"proud", "happy", "angry"},
			DifficultyLevel: "N3",
			OrderIndex:     12,
		},
		{
			QuestionText:   "What is the meaning of せいで?",
			CorrectAnswer:  "because of (negative reason)",
			WrongAnswers:   []string{"despite", "in order to", "even though"},
			DifficultyLevel: "N3",
			OrderIndex:     13,
		},
		{
			QuestionText:   "Choose the correct reading for 都合",
			CorrectAnswer:  "つごう",
			WrongAnswers:   []string{"とごう", "みやこあい", "としあい"},
			DifficultyLevel: "N3",
			OrderIndex:     14,
		},
		{
			QuestionText:   "What does 途中 mean?",
			CorrectAnswer:  "on the way, in the middle",
			WrongAnswers:   []string{"beginning", "ending", "outside"},
			DifficultyLevel: "N3",
			OrderIndex:     15,
		},

		// N2 Level Questions (3 questions)
		{
			QuestionText:   "What is the meaning of 思い切って?",
			CorrectAnswer:  "daringly, boldly, resolutely",
			WrongAnswers:   []string{"carefully", "slowly", "repeatedly"},
			DifficultyLevel: "N2",
			OrderIndex:     16,
		},
		{
			QuestionText:   "Choose the correct meaning for 〜にしては",
			CorrectAnswer:  "considering, for",
			WrongAnswers:   []string{"instead of", "because of", "in addition to"},
			DifficultyLevel: "N2",
			OrderIndex:     17,
		},
		{
			QuestionText:   "What does 遠慮なく mean?",
			CorrectAnswer:  "without hesitation, freely",
			WrongAnswers:   []string{"carefully", "politely", "secretly"},
			DifficultyLevel: "N2",
			OrderIndex:     18,
		},

		// N1 Level Questions (2 questions - hardest)
		{
			QuestionText:   "What is the meaning of 〜ならでは?",
			CorrectAnswer:  "unique to, only possible with",
			WrongAnswers:   []string{"except for", "if not", "even though"},
			DifficultyLevel: "N1",
			OrderIndex:     19,
		},
		{
			QuestionText:   "Choose the correct meaning for 言うまでもない",
			CorrectAnswer:  "needless to say, it goes without saying",
			WrongAnswers:   []string{"unable to say", "must be said", "said before"},
			DifficultyLevel: "N1",
			OrderIndex:     20,
		},
	}

	log.Printf("Seeding %d placement test questions...", len(questions))

	// Insert questions
	for _, q := range questions {
		wrongAnswersJSON, err := json.Marshal(q.WrongAnswers)
		if err != nil {
			log.Fatalf("Failed to marshal wrong answers for question %d: %v", q.OrderIndex, err)
		}

		_, err = db.Exec(`
			INSERT INTO placement_questions
			(question_text, correct_answer, wrong_answers, difficulty_level, order_index)
			VALUES ($1, $2, $3, $4, $5)
		`, q.QuestionText, q.CorrectAnswer, wrongAnswersJSON, q.DifficultyLevel, q.OrderIndex)

		if err != nil {
			log.Fatalf("Failed to insert question %d: %v", q.OrderIndex, err)
		}
	}

	log.Println("Successfully seeded placement test questions!")
	log.Printf("Total questions inserted: %d", len(questions))
	log.Println("Distribution: 5 N5, 5 N4, 5 N3, 3 N2, 2 N1")
}
