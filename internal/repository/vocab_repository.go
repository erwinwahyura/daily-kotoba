package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/yourusername/kotoba-api/internal/models"
)

type VocabRepository struct {
	db *sql.DB
}

func NewVocabRepository(db *sql.DB) *VocabRepository {
	return &VocabRepository{db: db}
}

func (r *VocabRepository) GetByLevelAndIndex(level string, index int) (*models.Vocabulary, error) {
	vocab := &models.Vocabulary{}
	query := `
		SELECT id, word, reading, short_meaning, detailed_explanation,
		       example_sentences, usage_notes, jlpt_level, index_position, created_at
		FROM vocabulary
		WHERE jlpt_level = $1 AND index_position = $2
	`
	err := r.db.QueryRow(query, level, index).Scan(
		&vocab.ID,
		&vocab.Word,
		&vocab.Reading,
		&vocab.ShortMeaning,
		&vocab.DetailedExplanation,
		&vocab.ExampleSentences,
		&vocab.UsageNotes,
		&vocab.JLPTLevel,
		&vocab.IndexPosition,
		&vocab.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("vocabulary not found")
	}
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (r *VocabRepository) GetByID(id string) (*models.Vocabulary, error) {
	vocab := &models.Vocabulary{}
	query := `
		SELECT id, word, reading, short_meaning, detailed_explanation,
		       example_sentences, usage_notes, jlpt_level, index_position, created_at
		FROM vocabulary
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&vocab.ID,
		&vocab.Word,
		&vocab.Reading,
		&vocab.ShortMeaning,
		&vocab.DetailedExplanation,
		&vocab.ExampleSentences,
		&vocab.UsageNotes,
		&vocab.JLPTLevel,
		&vocab.IndexPosition,
		&vocab.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("vocabulary not found")
	}
	if err != nil {
		return nil, err
	}

	return vocab, nil
}

func (r *VocabRepository) GetByLevel(level string, page, limit int) ([]models.Vocabulary, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM vocabulary WHERE jlpt_level = $1`
	err := r.db.QueryRow(countQuery, level).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	query := `
		SELECT id, word, reading, short_meaning, detailed_explanation,
		       example_sentences, usage_notes, jlpt_level, index_position, created_at
		FROM vocabulary
		WHERE jlpt_level = $1
		ORDER BY index_position
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, level, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var vocabList []models.Vocabulary
	for rows.Next() {
		var vocab models.Vocabulary
		err := rows.Scan(
			&vocab.ID,
			&vocab.Word,
			&vocab.Reading,
			&vocab.ShortMeaning,
			&vocab.DetailedExplanation,
			&vocab.ExampleSentences,
			&vocab.UsageNotes,
			&vocab.JLPTLevel,
			&vocab.IndexPosition,
			&vocab.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		vocabList = append(vocabList, vocab)
	}

	return vocabList, total, nil
}

func (r *VocabRepository) GetTotalCountByLevel(level string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM vocabulary WHERE jlpt_level = $1`
	err := r.db.QueryRow(query, level).Scan(&count)
	return count, err
}

func (r *VocabRepository) Create(vocab *models.Vocabulary) error {
	query := `
		INSERT INTO vocabulary (word, reading, short_meaning, detailed_explanation,
		                       example_sentences, usage_notes, jlpt_level, index_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		vocab.Word,
		vocab.Reading,
		vocab.ShortMeaning,
		vocab.DetailedExplanation,
		vocab.ExampleSentences,
		vocab.UsageNotes,
		vocab.JLPTLevel,
		vocab.IndexPosition,
	).Scan(&vocab.ID, &vocab.CreatedAt)
}

func (r *VocabRepository) BulkCreate(vocabList []models.Vocabulary) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO vocabulary (word, reading, short_meaning, detailed_explanation,
		                       example_sentences, usage_notes, jlpt_level, index_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, vocab := range vocabList {
		_, err := stmt.Exec(
			vocab.Word,
			vocab.Reading,
			vocab.ShortMeaning,
			vocab.DetailedExplanation,
			vocab.ExampleSentences,
			vocab.UsageNotes,
			vocab.JLPTLevel,
			vocab.IndexPosition,
		)
		if err != nil {
			return fmt.Errorf("failed to insert vocab %s: %w", vocab.Word, err)
		}
	}

	return tx.Commit()
}
