package models

import "time"

type ReadingArticle struct {
	ID         string    `json:"id" db:"id"`
	JLPTLevel  string    `json:"jlpt_level" db:"jlpt_level"`
	Title      string    `json:"title" db:"title"`
	Text       string    `json:"text" db:"text"`
	Topic      string    `json:"topic" db:"topic"`
	Difficulty string    `json:"difficulty" db:"difficulty"`
	WordCount  int       `json:"word_count" db:"word_count"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	Questions  []ReadingQuestion `json:"questions,omitempty"`
}

type ReadingQuestion struct {
	ID           string   `json:"id" db:"id"`
	ArticleID    string   `json:"article_id" db:"article_id"`
	QuestionNum  int      `json:"question_num" db:"question_num"`
	Question     string   `json:"question" db:"question"`
	Options      []string `json:"options" db:"options"`
	CorrectIndex int      `json:"correct_index" db:"correct_index"`
	Explanation  string   `json:"explanation" db:"explanation"`
}

type ReadingSession struct {
	ArticleID string         `json:"article_id"`
	Article   *ReadingArticle `json:"article"`
	Answers   map[string]int `json:"answers,omitempty"`
	Score     int            `json:"score,omitempty"`
	Total     int            `json:"total,omitempty"`
}

type SubmitReadingRequest struct {
	ArticleID string         `json:"article_id"`
	Answers   map[string]int `json:"answers"`
}
