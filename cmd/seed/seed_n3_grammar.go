package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/yourusername/kotoba-api/internal/config"
	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
	"github.com/yourusername/kotoba-api/internal/repository"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	sqlDB, err := cfg.GetDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	wrappedDB := db.New(sqlDB, cfg.DB.Driver)
	if cfg.DB.Driver == "sqlite" {
		wrappedDB.InitializeSQLite()
	}

	log.Println("Connected to database successfully")

	// Clear existing N3 grammar patterns
	log.Println("Clearing existing N3 grammar patterns...")
	_, err = wrappedDB.Exec("DELETE FROM grammar_patterns WHERE jlpt_level = 'N3'")
	if err != nil {
		log.Fatalf("Failed to clear existing patterns: %v", err)
	}
	log.Println("Existing N3 grammar patterns cleared")

	grammarRepo := repository.NewGrammarRepository(wrappedDB)

	// Seed N3 grammar patterns (high-priority forms that learners confuse)
	patterns := []models.GrammarPattern{
		{
			Pattern:             "〜わけにはいかない",
			PlainForm:           "わけにはいかない",
			Meaning:             "cannot afford to; must not (due to circumstances/obligation)",
			DetailedExplanation: "Expresses that an action is impossible because of circumstances, social obligations, or moral constraints. STRONGER than a simple 'cannot' — it implies 'even though I might want to, I can't because of X'. Often used when declining invitations or avoiding actions that would cause problems.",
			ConjugationRules:    "Verb (dictionary form) + わけにはいかない",
			NuanceNotes:         "This is NOT the same as 〜わけがない (no way that). Compare: 行くわけにはいかない (can't go due to circumstances) vs 行くわけがない (no way I'd go / impossible that I'd go). The first is 'shouldn't/can't due to situation' — the second is 'utterly impossible'.",
			JLPTLevel:           "N3",
			IndexPosition:       0,
			UsageExamples: models.UsageExamples{
				{
					Japanese:    "明日試験があるから、今夜遊びに行くわけにはいかない",
					Reading:     "あしたしけんがあるから、こんやあそびにいくわけにはいかない",
					Meaning:     "I have an exam tomorrow, so I can't afford to go out tonight",
					Nuance:      "External obligation (exam) prevents the action",
					Context:     "Declining a friend's invitation to party",
					Alternative: "明日試験があるから、今夜遊びに行けない (simpler, less formal)",
				},
				{
					Japanese:    "約束したから、忘れるわけにはいかない",
					Reading:     "やくそくしたから、わすれるわけにはいかない",
					Meaning:     "I made a promise, so I can't forget",
					Nuance:      "Moral/social obligation makes the action impossible",
					Context:     "Reminding yourself of an important commitment",
					Alternative: "約束したから、忘れられない (less formal nuance)",
				},
				{
					Japanese:    "上司がいるのに、先に帰るわけにはいかない",
					Reading:     "じょうしがいるのに、さきにかえるわけにはいかない",
					Meaning:     "My boss is here, so I can't leave first",
					Nuance:      "Social circumstance prevents the action",
					Context:     "Workplace situation — wanting to leave but shouldn't",
					Alternative: "",
				},
			},
			RelatedPatterns: models.RelatedPatterns{
				{Pattern: "〜わけがない", Relationship: "often confused with — COMPLETELY DIFFERENT", KeyDifference: "〜わけにはいかない = cannot afford to (circumstances); 〜わけがない = no way that (disbelief)"},
				{Pattern: "〜ないわけにはいかない", Relationship: "negative form — RARE, careful!", KeyDifference: "Double negative = 'must do' (cannot NOT do). Use 〜なければならない instead — clearer."},
				{Pattern: "〜べきではない", Relationship: "similar prohibition", KeyDifference: "〜べきではない = shouldn't (moral advice); 〜わけにはいかない = can't due to circumstances"},
			},
			CommonMistakes: "1) DO NOT confuse with 〜わけがない (no way that) — different meaning entirely. 2) Don't use with adjectives directly — needs verb stem. 3) Negative form (〜ないわけにはいかない) is tricky — avoid unless advanced.",
		},
		{
			Pattern:             "〜わけがない",
			PlainForm:           "わけがない",
			Meaning:             "no way that; impossible that; there's no reason that",
			DetailedExplanation: "Expresses strong disbelief or impossibility. Like saying 'that can't be' or 'no way'. Often used when rejecting rumors, defending someone, or expressing skepticism. This is EMOTIONAL — you're reacting to something surprising or outrageous.",
			ConjugationRules:    "Verb (dictionary/plain form) + わけがない; い-adjective + わけがない; な-adjective + な + わけがない; Noun + の + わけがない",
			NuanceNotes:         "Compare: 〜はずがない (also 'unlikely') vs 〜わけがない. 〜はずがない = logical impossibility (planning/scheduling context). 〜わけがない = emotional disbelief (reaction to claims about people/situations).",
			JLPTLevel:           "N3",
			IndexPosition:       1,
			UsageExamples: models.UsageExamples{
				{
					Japanese:    "彼がそんなことをするわけがない",
					Reading:     "かれがそんなことをするわけがない",
					Meaning:     "There's no way he'd do something like that",
					Nuance:      "Defending someone's character against accusation",
					Context:     "Someone gossips about your friend; you defend them",
					Alternative: "彼がそんなことするはずがない (more logical, less emotional)",
				},
				{
					Japanese:    "彼女は日本に行ったことがない。日本語が上手なわけがない",
					Reading:     "かのじょはにほんにいったことがない。にほんごがじょうずなわけがない",
					Meaning:     "She's never been to Japan. There's no way her Japanese is good",
					Nuance:      "Logical impossibility based on evidence",
					Context:     "Reacting to someone claiming unrealistic skill",
					Alternative: "日本語が上手なはずがない (more logical deduction)",
				},
				{
					Japanese:    "そんな高いものが買えるわけがない",
					Reading:     "そんなたかいものがかえるわけがない",
					Meaning:     "There's no way I can buy something that expensive",
					Nuance:      "Expressing impossibility due to financial reality",
					Context:     "Shopping, seeing expensive item, reacting",
					Alternative: "",
				},
			},
			RelatedPatterns: models.RelatedPatterns{
				{Pattern: "〜わけにはいかない", Relationship: "COMPLETELY DIFFERENT — most common confusion", KeyDifference: "〜わけがない = no way that (disbelief); 〜わけにはいかない = cannot afford to (circumstances)"},
				{Pattern: "〜はずがない", Relationship: "similar 'impossibility' but different nuance", KeyDifference: "〜はずがない = logical impossibility (planning); 〜わけがない = emotional disbelief (character/situation claims)"},
			},
			CommonMistakes: "THIS IS THE #1 CONFUSED PATTERN with 〜わけにはいかない. They mean OPPOSITE things. Remember: わけがない = NO WAY (strong negative); わけにはいかない = CAN'T AFFORD TO (obligation blocks action).",
		},
		{
			Pattern:             "〜ものだ",
			PlainForm:           "ものだ",
			Meaning:             "used to; express emotion; general truth; soft advice",
			DetailedExplanation: "Multiple uses: (1) Nostalgia — 子供の頃、よく公園で遊んだものだ ('I used to play...'), (2) General truths — 時間が経つのは早いものだ ('Time flies'), (3) Emotional exclamations — いいものだね ('How nice!'), (4) Soft advice — もっと勉強するものだ ('You should study more' — softer than べき). The common thread: expressing human experience, sentiment, or gentle guidance.",
			ConjugationRules:    "Verb (plain past/present) + ものだ; い-adjective + ものだ; な-adjective + な + ものだ; casual contraction: もんだ",
			NuanceNotes:         "SOFTER than 〜べきだ (obligation advice). 〜ものだ expresses general wisdom/sentiment. 〜べきだ expresses moral duty. Compare: 学生は勉強するべきだ (strong: students MUST study) vs 学生は勉強するものだ (gentle: it's what students do).",
			JLPTLevel:           "N3",
			IndexPosition:       2,
			UsageExamples: models.UsageExamples{
				{
					Japanese:    "子供の頃、夏休みになると海に行ったものだ",
					Reading:     "こどものころ、なつやすみになるとうみにいったものだ",
					Meaning:     "When I was a child, I used to go to the beach every summer vacation",
					Nuance:      "Nostalgic reflection on past habits",
					Context:     "Reminiscing with family or friends",
					Alternative: "海に行った (simple past, less nostalgic)",
				},
				{
					Japanese:    "年を取ると、体力が落ちるものだ",
					Reading:     "としをとると、たいりょくがおちるものだ",
					Meaning:     "As you get older, your physical strength declines (general truth)",
					Nuance:      "Stating universal human experience",
					Context:     "Sympathizing with someone getting older",
					Alternative: "体力が落ちる (simple statement, less empathetic)",
				},
				{
					Japanese:    "学生はもっと本を読むものだ",
					Reading:     "がくせいはもっとほんをよむものだ",
					Meaning:     "Students should read more books (gentle advice)",
					Nuance:      "Soft guidance rather than obligation",
					Context:     "Giving advice without sounding preachy",
					Alternative: "学生はもっと本を読むべきだ (stronger obligation nuance)",
				},
			},
			RelatedPatterns: models.RelatedPatterns{
				{Pattern: "〜ものではない", Relationship: "similar form, OPPOSITE meaning", KeyDifference: "〜ものだ = general truth/sentiment; 〜ものではない = should never do (strong prohibition)"},
				{Pattern: "〜べきだ", Relationship: "similar advice meaning", KeyDifference: "〜べきだ = moral obligation/should; 〜ものだ = gentle sentiment/soft advice"},
				{Pattern: "〜ことだ", Relationship: "similar advice, slightly more direct", KeyDifference: "〜ことだ = advice about specific action; 〜ものだ = general truth/wisdom about human nature"},
			},
			CommonMistakes: "1) Don't confuse with 〜ものではない (prohibition). 2) 〜もんだ is casual speech contraction, NOT a different grammar point. 3) Don't use for specific one-time events — this is for general truths/nostalgia.",
		},
		{
			Pattern:             "〜ものではない",
			PlainForm:           "ものではない",
			Meaning:             "should not; must never; is not appropriate",
			DetailedExplanation: "Strong prohibition or judgment that something is inappropriate. Used for: (1) Moral/social prohibitions — 人の悪口を言うものではない ('one shouldn't speak ill of others'), (2) Practical advice — そんなに焦るものではない ('you shouldn't rush like that'). More EMPHATIC than 〜べきではない.",
			ConjugationRules:    "Verb (dictionary form) + ものではない; な-adjective + である + ものではない (formal)",
			NuanceNotes:         "Stronger than 〜べきではない. 〜べきではない = 'shouldn't' (advice); 〜ものではない = 'one must never / it is not the done thing' (social/moral judgment). This carries cultural weight about proper behavior.",
			JLPTLevel:           "N3",
			IndexPosition:       3,
			UsageExamples: models.UsageExamples{
				{
					Japanese:    "人の失敗を笑うものではない",
					Reading:     "ひとのしっぱいをわらうものではない",
					Meaning:     "One should not laugh at others' failures",
					Nuance:      "Moral/social judgment about proper behavior",
					Context:     "Someone laughing at colleague's mistake — you correct them",
					Alternative: "人の失敗を笑うべきではない (weaker, more personal advice)",
				},
				{
					Japanese:    "そんなに焦るものではない。ゆっくりやればいい",
					Reading:     "そんなにあせるものではない。ゆっくりやればいい",
					Meaning:     "You shouldn't rush like that. Take your time",
					Nuance:      "Practical advice with gentle scolding",
					Context:     "Friend panicking before exam — calming them",
					Alternative: "焦るべきではない (more direct)",
				},
				{
					Japanese:    "約束を破るものではない",
					Reading:     "やくそくをやぶるものではない",
					Meaning:     "One must not break promises",
					Nuance:      "Strong moral statement",
					Context:     "Someone considering backing out of commitment",
					Alternative: "約束を破るべきではない (less culturally weighty)",
				},
			},
			RelatedPatterns: models.RelatedPatterns{
				{Pattern: "〜ものだ", Relationship: "similar form, OPPOSITE meaning", KeyDifference: "〜ものだ = general truth/sentiment; 〜ものではない = strong prohibition"},
				{Pattern: "〜べきではない", Relationship: "similar prohibition, weaker", KeyDifference: "〜べきではない = shouldn't (personal advice); 〜ものではない = must never (social/moral rule)"},
				{Pattern: "〜てはいけない", Relationship: "casual prohibition", KeyDifference: "〜てはいけない = don't do (casual/children); 〜ものではない = one must never (adult moral judgment)"},
			},
			CommonMistakes: "THIS IS NOT 〜ものだ. They mean opposite things. Remember: 〜ものだ = (nostalgic/gentle); 〜ものではない = (prohibition). Also: this is FORMAL — don't use in casual conversation with friends.",
		},
		{
			Pattern:             "〜ばかり",
			PlainForm:           "ばかり",
			Meaning:             "only; nothing but; keep doing; approximately",
			DetailedExplanation: "Multiple distinct uses: (1) 食べてばかりいる = 'keeps eating / does nothing but eat' (negative criticism of repeated action), (2) 学生ばかり = 'only students / nothing but students' (exclusivity, often negative), (3) 3時ばかり = 'around 3 o'clock' (approximate time). The negative connotation is key — unlike 〜だけ which is neutral.",
			ConjugationRules:    "Verb (te-form) + ばかり + いる (repeated action); Noun + ばかり (exclusivity); Time + ばかり (approximation); casual: ばっかり, ばっか",
			NuanceNotes:         "NEGATIVE nuance compared to 〜だけ. 学生だけ = only students (neutral statement); 学生ばかり = nothing but students (implies criticism, lack of diversity, too many). With verbs: ゲームしてばかり = keeps playing games (implies bad habit).",
			JLPTLevel:           "N3",
			IndexPosition:       4,
			UsageExamples: models.UsageExamples{
				{
					Japanese:    "彼はゲームをしてばかりいる",
					Reading:     "かれはゲームをしてばかりいる",
					Meaning:     "All he does is play games (negative: criticism of habit)",
					Nuance:      "Judgmental — implies wasted time, bad habit",
					Context:     "Complaining about partner/friend's gaming",
					Alternative: "ゲームをしているだけ (neutral: he just plays games — no judgment)",
				},
				{
					Japanese:    "ここは観光客ばかりだ",
					Reading:     "ここはかんこうきゃくばかりだ",
					Meaning:     "This place is nothing but tourists (implies lack of locals/too crowded)",
					Nuance:      "Negative: lack of diversity, authenticity",
					Context:     "Commenting on a tourist trap area",
					Alternative: "観光客だけ (neutral: just tourists)",
				},
				{
					Japanese:    "会議は2時間ばかり続いた",
					Reading:     "かいぎはにじかんばかりつづいた",
					Meaning:     "The meeting lasted about 2 hours",
					Nuance:      "Approximate time/amount",
					Context:     "Reporting meeting duration casually",
					Alternative: "2時間ぐらい (more common for time approximation)",
				},
			},
			RelatedPatterns: models.RelatedPatterns{
				{Pattern: "〜だけ", Relationship: "neutral 'only' — MOST IMPORTANT DISTINCTION", KeyDifference: "〜だけ = neutral 'only' (just stating fact); 〜ばかり = negative 'nothing but' (implies criticism, excess, bad habit)"},
				{Pattern: "〜のみ", Relationship: "formal written 'only'", KeyDifference: "〜のみ = formal written only; 〜ばかり = spoken, often negative/critical"},
				{Pattern: "〜ばかりか", Relationship: "advanced 'not only...but also'", KeyDifference: "〜ばかりか = not only X but also Y (emphasis); different grammar point entirely"},
			},
			CommonMistakes: "1) BIGGEST MISTAKE: Using 〜ばかり when you mean neutral 'only'. If you don't want to sound critical, use 〜だけ. 2) Position matters: ばかり食べる (rare, 'only eat') vs 食べてばかり (common, 'keep eating excessively'). 3) Time approximation: ばかり = older/literary style; 〜ぐらい more common now.",
		},
	}

	// Bulk insert grammar patterns
	log.Printf("Inserting %d N3 grammar patterns...", len(patterns))
	if err := grammarRepo.BulkCreate(patterns); err != nil {
		log.Fatalf("Failed to insert grammar patterns: %v", err)
	}

	log.Println("N3 grammar patterns seeded successfully!")
}
