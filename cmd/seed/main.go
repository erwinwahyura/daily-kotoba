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

	// Clear existing N4 vocabulary to avoid duplicates
	log.Println("Clearing existing N4 vocabulary...")
	_, err = db.Exec("DELETE FROM vocabulary WHERE jlpt_level = 'N4'")
	if err != nil {
		log.Fatalf("Failed to clear existing vocabulary: %v", err)
	}
	log.Println("Existing N4 vocabulary cleared")

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
		{
			Word:                "明るい",
			Reading:             "あかるい",
			ShortMeaning:        "bright, cheerful",
			DetailedExplanation: "Can describe physical brightness (light, colors) or personality (cheerful, optimistic). Common i-adjective in everyday conversation.",
			ExampleSentences: models.ExampleSentences{
				"この部屋は明るい - This room is bright",
				"彼女は明るい性格だ - She has a cheerful personality",
				"明るい未来 - A bright future",
			},
			UsageNotes:    "When describing personality, it means cheerful or positive. Very commonly used.",
			JLPTLevel:     "N4",
			IndexPosition: 5,
		},
		{
			Word:                "上がる",
			Reading:             "あがる",
			ShortMeaning:        "to rise, to go up",
			DetailedExplanation: "Indicates upward movement or increase. Can be used for prices, temperatures, stairs, or entering someone's home.",
			ExampleSentences: models.ExampleSentences{
				"温度が上がった - The temperature rose",
				"階段を上がる - To go up the stairs",
				"家に上がってください - Please come into my house",
			},
			UsageNotes:    "上がる is intransitive. The transitive version is 上げる (to raise).",
			JLPTLevel:     "N4",
			IndexPosition: 6,
		},
		{
			Word:                "集まる",
			Reading:             "あつまる",
			ShortMeaning:        "to gather, to collect",
			DetailedExplanation: "Describes people or things coming together in one place. Intransitive verb commonly used for meetings, crowds, or collections.",
			ExampleSentences: models.ExampleSentences{
				"みんなが公園に集まった - Everyone gathered at the park",
				"お金が集まらない - Money isn't being collected",
				"駅前に人が集まっている - People are gathering in front of the station",
			},
			UsageNotes:    "Intransitive verb. The transitive form is 集める (to gather/collect something).",
			JLPTLevel:     "N4",
			IndexPosition: 7,
		},
		{
			Word:                "謝る",
			Reading:             "あやまる",
			ShortMeaning:        "to apologize",
			DetailedExplanation: "To offer an apology or express regret. Essential in Japanese culture for maintaining social harmony.",
			ExampleSentences: models.ExampleSentences{
				"彼に謝った - I apologized to him",
				"謝る必要はない - There's no need to apologize",
				"素直に謝りなさい - Apologize honestly",
			},
			UsageNotes:    "Often used with に particle for the person being apologized to.",
			JLPTLevel:     "N4",
			IndexPosition: 8,
		},
		{
			Word:                "安全",
			Reading:             "あんぜん",
			ShortMeaning:        "safety, security",
			DetailedExplanation: "Na-adjective and noun meaning safe or safety. Commonly used in warnings, instructions, and discussions about security.",
			ExampleSentences: models.ExampleSentences{
				"この地域は安全です - This area is safe",
				"安全運転を心がける - Be mindful of safe driving",
				"安全を確認する - Confirm safety",
			},
			UsageNotes:    "Can be used as both noun and na-adjective (安全な場所 - a safe place).",
			JLPTLevel:     "N4",
			IndexPosition: 9,
		},
		{
			Word:                "以上",
			Reading:             "いじょう",
			ShortMeaning:        "more than, above",
			DetailedExplanation: "Indicates 'more than' or 'at least' a certain amount. Also used to conclude speeches or presentations (meaning 'that's all').",
			ExampleSentences: models.ExampleSentences{
				"18歳以上 - 18 years old or older",
				"100人以上 - More than 100 people",
				"以上です - That's all (ending a speech)",
			},
			UsageNotes:    "When placed after a number, it includes that number and higher.",
			JLPTLevel:     "N4",
			IndexPosition: 10,
		},
		{
			Word:                "以内",
			Reading:             "いない",
			ShortMeaning:        "within, inside",
			DetailedExplanation: "Indicates a limit or boundary of time, distance, or quantity that should not be exceeded.",
			ExampleSentences: models.ExampleSentences{
				"3日以内に - Within 3 days",
				"100メートル以内 - Within 100 meters",
				"予算以内で買う - Buy within the budget",
			},
			UsageNotes:    "The limit mentioned is included (e.g., 3日以内 includes day 3).",
			JLPTLevel:     "N4",
			IndexPosition: 11,
		},
		{
			Word:                "祝う",
			Reading:             "いわう",
			ShortMeaning:        "to celebrate",
			DetailedExplanation: "To celebrate or commemorate a happy occasion. Used for birthdays, holidays, achievements, etc.",
			ExampleSentences: models.ExampleSentences{
				"誕生日を祝う - Celebrate a birthday",
				"合格を祝ってパーティーをした - Had a party to celebrate passing",
				"お祝いする - To celebrate",
			},
			UsageNotes:    "Often paired with を particle for what is being celebrated.",
			JLPTLevel:     "N4",
			IndexPosition: 12,
		},
		{
			Word:                "植える",
			Reading:             "うえる",
			ShortMeaning:        "to plant",
			DetailedExplanation: "To plant seeds, flowers, or trees in the ground. Can also mean to implant or embed something.",
			ExampleSentences: models.ExampleSentences{
				"庭に花を植える - Plant flowers in the garden",
				"木を植えた - Planted a tree",
				"種を植える - Plant seeds",
			},
			UsageNotes:    "Common in gardening contexts. Transitive verb.",
			JLPTLevel:     "N4",
			IndexPosition: 13,
		},
		{
			Word:                "受ける",
			Reading:             "うける",
			ShortMeaning:        "to receive, to take (exam)",
			DetailedExplanation: "Has multiple meanings: to receive, to take (a test), to accept, or to undergo. Context determines the specific meaning.",
			ExampleSentences: models.ExampleSentences{
				"試験を受ける - Take an exam",
				"授業を受ける - Attend a class",
				"影響を受ける - Be influenced",
			},
			UsageNotes:    "Very versatile verb. Pay attention to context for correct meaning.",
			JLPTLevel:     "N4",
			IndexPosition: 14,
		},
		{
			Word:                "動く",
			Reading:             "うごく",
			ShortMeaning:        "to move",
			DetailedExplanation: "Indicates movement or motion. Can refer to physical movement, operating (machines), or being touched emotionally.",
			ExampleSentences: models.ExampleSentences{
				"車が動き出した - The car started moving",
				"心が動いた - My heart was moved",
				"機械が動かない - The machine doesn't work",
			},
			UsageNotes:    "Intransitive verb. Transitive form is 動かす (to move something).",
			JLPTLevel:     "N4",
			IndexPosition: 15,
		},
		{
			Word:                "打つ",
			Reading:             "うつ",
			ShortMeaning:        "to hit, to type",
			DetailedExplanation: "Versatile verb: to hit, strike, type (keyboard), or inject. Meaning depends heavily on context.",
			ExampleSentences: models.ExampleSentences{
				"ボールを打つ - Hit the ball",
				"メールを打つ - Type an email",
				"注射を打つ - Give an injection",
			},
			UsageNotes:    "Very common in sports and computer contexts.",
			JLPTLevel:     "N4",
			IndexPosition: 16,
		},
		{
			Word:                "映る",
			Reading:             "うつる",
			ShortMeaning:        "to be reflected, to appear",
			DetailedExplanation: "To be reflected (in mirror, water) or to appear on screen. Intransitive verb.",
			ExampleSentences: models.ExampleSentences{
				"鏡に映る - Be reflected in the mirror",
				"画面に映っている - Appearing on the screen",
				"写真に映る - Appear in a photo",
			},
			UsageNotes:    "Different kanji from 移る (to move, transfer). Be careful with pronunciation context.",
			JLPTLevel:     "N4",
			IndexPosition: 17,
		},
		{
			Word:                "選ぶ",
			Reading:             "えらぶ",
			ShortMeaning:        "to choose, to select",
			DetailedExplanation: "To make a choice or selection from multiple options. Common in shopping, decision-making contexts.",
			ExampleSentences: models.ExampleSentences{
				"好きなものを選んでください - Please choose what you like",
				"色を選ぶ - Choose a color",
				"代表を選ぶ - Select a representative",
			},
			UsageNotes:    "Often paired with を particle for what is being chosen.",
			JLPTLevel:     "N4",
			IndexPosition: 18,
		},
		{
			Word:                "遠慮",
			Reading:             "えんりょ",
			ShortMeaning:        "reserve, restraint",
			DetailedExplanation: "To hold back, refrain, or be reserved out of politeness. Important cultural concept in Japanese social interactions.",
			ExampleSentences: models.ExampleSentences{
				"遠慮しないで食べてください - Please don't hold back, eat freely",
				"遠慮なくどうぞ - Please, without hesitation",
				"遠慮する - To refrain/hold back",
			},
			UsageNotes:    "Often used in negative form (遠慮しないで) to encourage someone not to hold back.",
			JLPTLevel:     "N4",
			IndexPosition: 19,
		},
		{
			Word:                "億",
			Reading:             "おく",
			ShortMeaning:        "hundred million",
			DetailedExplanation: "Counter word for hundred million (100,000,000). Essential for understanding large numbers in Japanese.",
			ExampleSentences: models.ExampleSentences{
				"一億円 - One hundred million yen",
				"三億人 - Three hundred million people",
				"数億年前 - Hundreds of millions of years ago",
			},
			UsageNotes:    "Japanese counting system groups by 万 (ten thousand) and 億 (hundred million).",
			JLPTLevel:     "N4",
			IndexPosition: 20,
		},
		{
			Word:                "贈る",
			Reading:             "おくる",
			ShortMeaning:        "to give, to present",
			DetailedExplanation: "To give a gift or present formally. More formal than あげる. Used for meaningful gifts.",
			ExampleSentences: models.ExampleSentences{
				"プレゼントを贈る - Give a present",
				"花を贈った - Presented flowers",
				"言葉を贈る - Offer words (of encouragement)",
			},
			UsageNotes:    "More formal than あげる. Often used for special occasions.",
			JLPTLevel:     "N4",
			IndexPosition: 21,
		},
		{
			Word:                "怒る",
			Reading:             "おこる",
			ShortMeaning:        "to get angry",
			DetailedExplanation: "To become angry or mad. Can be intransitive (to get angry) or transitive (to scold).",
			ExampleSentences: models.ExampleSentences{
				"先生に怒られた - Was scolded by the teacher",
				"彼は怒っている - He is angry",
				"怒らないでください - Please don't get angry",
			},
			UsageNotes:    "怒る (to get angry) vs 起こる (to happen/occur) - different kanji, same reading.",
			JLPTLevel:     "N4",
			IndexPosition: 22,
		},
		{
			Word:                "遅れる",
			Reading:             "おくれる",
			ShortMeaning:        "to be late",
			DetailedExplanation: "To be late, delayed, or behind schedule. Very commonly used in daily life.",
			ExampleSentences: models.ExampleSentences{
				"電車が遅れている - The train is delayed",
				"授業に遅れた - Was late for class",
				"遅れてすみません - Sorry for being late",
			},
			UsageNotes:    "Intransitive verb. Often used with に particle for what you're late for.",
			JLPTLevel:     "N4",
			IndexPosition: 23,
		},
		{
			Word:                "お嬢さん",
			Reading:             "おじょうさん",
			ShortMeaning:        "daughter, young lady",
			DetailedExplanation: "Respectful term for someone else's daughter or a young woman. Cannot be used for your own daughter.",
			ExampleSentences: models.ExampleSentences{
				"お嬢さんはおいくつですか - How old is your daughter?",
				"可愛いお嬢さんですね - What a lovely daughter you have",
				"社長のお嬢さん - The president's daughter",
			},
			UsageNotes:    "Honorific term. Use 娘 (musume) for your own daughter.",
			JLPTLevel:     "N4",
			IndexPosition: 24,
		},
		{
			Word:                "お祝い",
			Reading:             "おいわい",
			ShortMeaning:        "celebration, congratulations",
			DetailedExplanation: "A celebration, congratulatory gift, or the act of celebrating. Noun form of 祝う.",
			ExampleSentences: models.ExampleSentences{
				"結婚のお祝い - Wedding celebration/gift",
				"お祝いを渡す - Give a congratulatory gift",
				"お祝いの言葉 - Words of congratulation",
			},
			UsageNotes:    "Often used for gift-giving occasions (weddings, graduations, etc.).",
			JLPTLevel:     "N4",
			IndexPosition: 25,
		},
		{
			Word:                "落ちる",
			Reading:             "おちる",
			ShortMeaning:        "to fall, to drop",
			DetailedExplanation: "To fall down, drop, or fail (an exam). Multiple meanings depending on context.",
			ExampleSentences: models.ExampleSentences{
				"りんごが木から落ちた - An apple fell from the tree",
				"試験に落ちた - Failed the exam",
				"携帯を落とした - Dropped my phone",
			},
			UsageNotes:    "Intransitive. Transitive form is 落とす. Also means 'to fail' for tests.",
			JLPTLevel:     "N4",
			IndexPosition: 26,
		},
		{
			Word:                "踊る",
			Reading:             "おどる",
			ShortMeaning:        "to dance",
			DetailedExplanation: "To dance or perform a dance. Used for all types of dancing.",
			ExampleSentences: models.ExampleSentences{
				"音楽に合わせて踊る - Dance to the music",
				"ダンスを踊る - Perform a dance",
				"みんなで踊った - Everyone danced together",
			},
			UsageNotes:    "Common verb for any type of dancing.",
			JLPTLevel:     "N4",
			IndexPosition: 27,
		},
		{
			Word:                "驚く",
			Reading:             "おどろく",
			ShortMeaning:        "to be surprised",
			DetailedExplanation: "To be surprised, shocked, or astonished. Intransitive verb expressing sudden emotional reaction.",
			ExampleSentences: models.ExampleSentences{
				"結果に驚いた - Was surprised at the result",
				"とても驚いています - I'm very surprised",
				"驚くことはない - There's nothing to be surprised about",
			},
			UsageNotes:    "Often used with に particle for the cause of surprise.",
			JLPTLevel:     "N4",
			IndexPosition: 28,
		},
		{
			Word:                "お礼",
			Reading:             "おれい",
			ShortMeaning:        "thanks, gratitude",
			DetailedExplanation: "Expression of thanks or gratitude. More formal than ありがとう. Can also mean a thank-you gift.",
			ExampleSentences: models.ExampleSentences{
				"お礼を言う - Express thanks",
				"お礼の品 - A thank-you gift",
				"お礼に何かしたい - I want to do something to thank you",
			},
			UsageNotes:    "More formal than ありがとう. Often used with を言う (to say thanks).",
			JLPTLevel:     "N4",
			IndexPosition: 29,
		},
		{
			Word:                "会場",
			Reading:             "かいじょう",
			ShortMeaning:        "venue, hall",
			DetailedExplanation: "A venue, hall, or location where an event takes place. Common for concerts, meetings, ceremonies.",
			ExampleSentences: models.ExampleSentences{
				"会場に着いた - Arrived at the venue",
				"会場はどこですか - Where is the venue?",
				"会場を予約する - Reserve a venue",
			},
			UsageNotes:    "Common in event contexts (concerts, weddings, conferences).",
			JLPTLevel:     "N4",
			IndexPosition: 30,
		},
		{
			Word:                "飼う",
			Reading:             "かう",
			ShortMeaning:        "to keep (pet), to raise",
			DetailedExplanation: "To keep or raise animals as pets or livestock. Different from 育てる which is for children/plants.",
			ExampleSentences: models.ExampleSentences{
				"犬を飼っている - I have a dog (as a pet)",
				"猫を飼いたい - I want to keep a cat",
				"ペットを飼う - Keep a pet",
			},
			UsageNotes:    "Specifically for animals. Use 育てる for children or plants.",
			JLPTLevel:     "N4",
			IndexPosition: 31,
		},
		{
			Word:                "変える",
			Reading:             "かえる",
			ShortMeaning:        "to change",
			DetailedExplanation: "To change something (transitive). To alter or modify something intentionally.",
			ExampleSentences: models.ExampleSentences{
				"予定を変える - Change the schedule",
				"考えを変えた - Changed my mind",
				"服を着替える - Change clothes",
			},
			UsageNotes:    "Transitive verb. Intransitive form is 変わる (to change/be changed).",
			JLPTLevel:     "N4",
			IndexPosition: 32,
		},
		{
			Word:                "科学",
			Reading:             "かがく",
			ShortMeaning:        "science",
			DetailedExplanation: "Science as a field of study. Encompasses natural sciences, physics, chemistry, biology, etc.",
			ExampleSentences: models.ExampleSentences{
				"科学の授業 - Science class",
				"科学技術 - Science and technology",
				"科学者 - Scientist",
			},
			UsageNotes:    "Different from 化学 (chemistry) - same pronunciation, different kanji.",
			JLPTLevel:     "N4",
			IndexPosition: 33,
		},
		{
			Word:                "掛ける",
			Reading:             "かける",
			ShortMeaning:        "to hang, to call",
			DetailedExplanation: "Very versatile verb: to hang, to make a phone call, to multiply, to sit down, to wear (glasses), etc.",
			ExampleSentences: models.ExampleSentences{
				"壁に絵を掛ける - Hang a picture on the wall",
				"電話を掛ける - Make a phone call",
				"眼鏡を掛ける - Wear glasses",
			},
			UsageNotes:    "Extremely versatile verb with many meanings. Context is crucial.",
			JLPTLevel:     "N4",
			IndexPosition: 34,
		},
		{
			Word:                "片付ける",
			Reading:             "かたづける",
			ShortMeaning:        "to tidy up, to put away",
			DetailedExplanation: "To clean up, tidy, or put things in order. Essential daily life verb.",
			ExampleSentences: models.ExampleSentences{
				"部屋を片付ける - Clean up the room",
				"食器を片付けた - Put away the dishes",
				"机を片付けてください - Please tidy your desk",
			},
			UsageNotes:    "Very common in daily life for cleaning and organizing.",
			JLPTLevel:     "N4",
			IndexPosition: 35,
		},
		{
			Word:                "硬い",
			Reading:             "かたい",
			ShortMeaning:        "hard, stiff",
			DetailedExplanation: "Describes physical hardness (opposite of soft) or rigidity in rules/thinking. Can also mean formal or serious.",
			ExampleSentences: models.ExampleSentences{
				"このパンは硬い - This bread is hard",
				"硬い表情 - A stiff expression",
				"頭が硬い - Rigid thinking",
			},
			UsageNotes:    "Can be written with different kanji (固い、堅い) with slightly different nuances.",
			JLPTLevel:     "N4",
			IndexPosition: 36,
		},
		{
			Word:                "悲しい",
			Reading:             "かなしい",
			ShortMeaning:        "sad",
			DetailedExplanation: "Feeling of sadness or sorrow. Common i-adjective for expressing sad emotions.",
			ExampleSentences: models.ExampleSentences{
				"悲しいニュース - Sad news",
				"とても悲しかった - I was very sad",
				"悲しい気持ち - A sad feeling",
			},
			UsageNotes:    "Basic emotion adjective. Often paired with です/だった.",
			JLPTLevel:     "N4",
			IndexPosition: 37,
		},
		{
			Word:                "彼女",
			Reading:             "かのじょ",
			ShortMeaning:        "she, girlfriend",
			DetailedExplanation: "Means 'she' or 'girlfriend' depending on context. Third person pronoun for females or romantic partner.",
			ExampleSentences: models.ExampleSentences{
				"彼女は学生です - She is a student",
				"彼女ができた - I got a girlfriend",
				"彼女の名前 - Her name",
			},
			UsageNotes:    "Context determines if it means 'she' or 'girlfriend'. Opposite is 彼 (he/boyfriend).",
			JLPTLevel:     "N4",
			IndexPosition: 38,
		},
		{
			Word:                "通う",
			Reading:             "かよう",
			ShortMeaning:        "to commute, to attend",
			DetailedExplanation: "To regularly go to a place (school, work, etc.). Implies repeated travel to the same destination.",
			ExampleSentences: models.ExampleSentences{
				"学校に通う - Attend school",
				"会社に通っている - Commuting to work",
				"ジムに通う - Go to the gym regularly",
			},
			UsageNotes:    "Implies regular, repeated visits. Different from 行く (to go once).",
			JLPTLevel:     "N4",
			IndexPosition: 39,
		},
		{
			Word:                "乾く",
			Reading:             "かわく",
			ShortMeaning:        "to dry, to get thirsty",
			DetailedExplanation: "To become dry. Can refer to laundry drying, throat getting dry (thirsty), or things drying out.",
			ExampleSentences: models.ExampleSentences{
				"洗濯物が乾いた - The laundry dried",
				"喉が乾いた - I'm thirsty (lit: throat is dry)",
				"髪が乾く - Hair dries",
			},
			UsageNotes:    "Intransitive verb. Transitive form is 乾かす (to dry something).",
			JLPTLevel:     "N4",
			IndexPosition: 40,
		},
		{
			Word:                "可愛い",
			Reading:             "かわいい",
			ShortMeaning:        "cute, adorable",
			DetailedExplanation: "Cute, pretty, or adorable. Very commonly used in Japanese culture for anything appealing or charming.",
			ExampleSentences: models.ExampleSentences{
				"可愛い犬 - A cute dog",
				"この服可愛いね - This outfit is cute",
				"可愛い赤ちゃん - An adorable baby",
			},
			UsageNotes:    "Extremely common in Japanese. Used much more broadly than 'cute' in English.",
			JLPTLevel:     "N4",
			IndexPosition: 41,
		},
		{
			Word:                "代わり",
			Reading:             "かわり",
			ShortMeaning:        "substitute, replacement",
			DetailedExplanation: "A substitute, replacement, or alternative. Also used in the phrase 'instead of' (の代わりに).",
			ExampleSentences: models.ExampleSentences{
				"私の代わりに行ってください - Please go instead of me",
				"代わりの人 - A substitute person",
				"お茶の代わりにコーヒー - Coffee instead of tea",
			},
			UsageNotes:    "Often used with の and に to mean 'instead of' (の代わりに).",
			JLPTLevel:     "N4",
			IndexPosition: 42,
		},
		{
			Word:                "関係",
			Reading:             "かんけい",
			ShortMeaning:        "relation, connection",
			DetailedExplanation: "Relationship, connection, or relevance between things or people. Used in many contexts.",
			ExampleSentences: models.ExampleSentences{
				"関係ない - Not related/irrelevant",
				"良い関係 - A good relationship",
				"仕事に関係する - Related to work",
			},
			UsageNotes:    "Very common noun. Often used with する to mean 'to be related'.",
			JLPTLevel:     "N4",
			IndexPosition: 43,
		},
		{
			Word:                "感じる",
			Reading:             "かんじる",
			ShortMeaning:        "to feel, to sense",
			DetailedExplanation: "To feel, sense, or perceive something emotionally or physically.",
			ExampleSentences: models.ExampleSentences{
				"寒さを感じる - Feel the cold",
				"不安を感じた - Felt anxiety",
				"何も感じない - Don't feel anything",
			},
			UsageNotes:    "Often used with を particle for what is being felt.",
			JLPTLevel:     "N4",
			IndexPosition: 44,
		},
		{
			Word:                "簡単",
			Reading:             "かんたん",
			ShortMeaning:        "simple, easy",
			DetailedExplanation: "Simple, easy, or uncomplicated. Na-adjective describing lack of difficulty or complexity.",
			ExampleSentences: models.ExampleSentences{
				"簡単な問題 - An easy problem",
				"それは簡単だ - That's simple",
				"簡単に説明する - Explain simply",
			},
			UsageNotes:    "Na-adjective. Very commonly used to describe ease or simplicity.",
			JLPTLevel:     "N4",
			IndexPosition: 45,
		},
		{
			Word:                "気分",
			Reading:             "きぶん",
			ShortMeaning:        "mood, feeling",
			DetailedExplanation: "Mood, feeling, or state of mind. Can refer to physical condition or emotional state.",
			ExampleSentences: models.ExampleSentences{
				"気分が悪い - Feel sick/bad mood",
				"気分がいい - Feel good/in a good mood",
				"気分転換 - Change of mood/refreshment",
			},
			UsageNotes:    "Common in daily conversation. Often used with が for mood description.",
			JLPTLevel:     "N4",
			IndexPosition: 46,
		},
		{
			Word:                "決まる",
			Reading:             "きまる",
			ShortMeaning:        "to be decided",
			DetailedExplanation: "To be decided, settled, or fixed. Intransitive verb indicating that a decision has been made.",
			ExampleSentences: models.ExampleSentences{
				"予定が決まった - The plan has been decided",
				"まだ決まっていない - Not decided yet",
				"時間が決まる - The time is set",
			},
			UsageNotes:    "Intransitive verb. Transitive form is 決める (to decide).",
			JLPTLevel:     "N4",
			IndexPosition: 47,
		},
		{
			Word:                "着物",
			Reading:             "きもの",
			ShortMeaning:        "kimono, traditional clothing",
			DetailedExplanation: "Traditional Japanese clothing (kimono). Can also mean 'clothes' in general in some contexts.",
			ExampleSentences: models.ExampleSentences{
				"着物を着る - Wear a kimono",
				"美しい着物 - A beautiful kimono",
				"結婚式で着物を着た - Wore a kimono at the wedding",
			},
			UsageNotes:    "Primarily refers to traditional Japanese clothing.",
			JLPTLevel:     "N4",
			IndexPosition: 48,
		},
		{
			Word:                "禁煙",
			Reading:             "きんえん",
			ShortMeaning:        "no smoking",
			DetailedExplanation: "No smoking or smoking prohibition. Common sign/rule in public places.",
			ExampleSentences: models.ExampleSentences{
				"禁煙席 - Non-smoking seat",
				"ここは禁煙です - This is a non-smoking area",
				"禁煙を始めた - Started quitting smoking",
			},
			UsageNotes:    "Very common on signs. Opposite is 喫煙 (smoking allowed).",
			JLPTLevel:     "N4",
			IndexPosition: 49,
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
