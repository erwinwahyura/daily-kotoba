package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/yourusername/kotoba-api/internal/db"
	"github.com/yourusername/kotoba-api/internal/models"
)

type GrammarRepository struct {
	db *db.DB
}

func NewGrammarRepository(db *db.DB) *GrammarRepository {
	return &GrammarRepository{db: db}
}

func (r *GrammarRepository) GetByLevelAndIndex(level string, index int) (*models.GrammarPattern, error) {
	pattern := &models.GrammarPattern{}
	query := `
		SELECT id, pattern, plain_form, meaning, detailed_explanation,
		       conjugation_rules, usage_examples, nuance_notes, jlpt_level,
		       related_patterns, common_mistakes, index_position, created_at
		FROM grammar_patterns
		WHERE jlpt_level = $1 AND index_position = $2
	`
	err := r.db.QueryRow(query, level, index).Scan(
		&pattern.ID,
		&pattern.Pattern,
		&pattern.PlainForm,
		&pattern.Meaning,
		&pattern.DetailedExplanation,
		&pattern.ConjugationRules,
		&pattern.UsageExamples,
		&pattern.NuanceNotes,
		&pattern.JLPTLevel,
		&pattern.RelatedPatterns,
		&pattern.CommonMistakes,
		&pattern.IndexPosition,
		&pattern.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("grammar pattern not found")
	}
	if err != nil {
		return nil, err
	}

	return pattern, nil
}

func (r *GrammarRepository) GetByPattern(patternName string, level string) (*models.GrammarPattern, error) {
	p := &models.GrammarPattern{}
	query := `
		SELECT id, pattern, plain_form, meaning, detailed_explanation,
		       conjugation_rules, usage_examples, nuance_notes, jlpt_level,
		       related_patterns, common_mistakes, index_position, created_at
		FROM grammar_patterns
		WHERE pattern = $1 AND jlpt_level = $2
	`
	err := r.db.QueryRow(query, patternName, level).Scan(
		&p.ID, &p.Pattern, &p.PlainForm, &p.Meaning, &p.DetailedExplanation,
		&p.ConjugationRules, &p.UsageExamples, &p.NuanceNotes, &p.JLPTLevel,
		&p.RelatedPatterns, &p.CommonMistakes, &p.IndexPosition, &p.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("grammar pattern not found")
	}
	return p, err
}

func (r *GrammarRepository) GetByID(id string) (*models.GrammarPattern, error) {
	pattern := &models.GrammarPattern{}
	query := `
		SELECT id, pattern, plain_form, meaning, detailed_explanation,
		       conjugation_rules, usage_examples, nuance_notes, jlpt_level,
		       related_patterns, common_mistakes, index_position, created_at
		FROM grammar_patterns
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&pattern.ID,
		&pattern.Pattern,
		&pattern.PlainForm,
		&pattern.Meaning,
		&pattern.DetailedExplanation,
		&pattern.ConjugationRules,
		&pattern.UsageExamples,
		&pattern.NuanceNotes,
		&pattern.JLPTLevel,
		&pattern.RelatedPatterns,
		&pattern.CommonMistakes,
		&pattern.IndexPosition,
		&pattern.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("grammar pattern not found")
	}
	if err != nil {
		return nil, err
	}

	return pattern, nil
}

func (r *GrammarRepository) GetByLevel(level string, page, limit int) ([]models.GrammarPattern, int, error) {
	offset := (page - 1) * limit

	var total int
	countQuery := `SELECT COUNT(*) FROM grammar_patterns WHERE jlpt_level = $1`
	err := r.db.QueryRow(countQuery, level).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, pattern, plain_form, meaning, detailed_explanation,
		       conjugation_rules, usage_examples, nuance_notes, jlpt_level,
		       related_patterns, common_mistakes, index_position, created_at
		FROM grammar_patterns
		WHERE jlpt_level = $1
		ORDER BY index_position
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, level, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var patterns []models.GrammarPattern
	for rows.Next() {
		var p models.GrammarPattern
		err := rows.Scan(
			&p.ID, &p.Pattern, &p.PlainForm, &p.Meaning, &p.DetailedExplanation,
			&p.ConjugationRules, &p.UsageExamples, &p.NuanceNotes, &p.JLPTLevel,
			&p.RelatedPatterns, &p.CommonMistakes, &p.IndexPosition, &p.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		patterns = append(patterns, p)
	}

	return patterns, total, nil
}

func (r *GrammarRepository) GetTotalCountByLevel(level string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM grammar_patterns WHERE jlpt_level = $1`
	err := r.db.QueryRow(query, level).Scan(&count)
	return count, err
}

func (r *GrammarRepository) Create(pattern *models.GrammarPattern) error {
	query := `
		INSERT INTO grammar_patterns (pattern, plain_form, meaning, detailed_explanation,
		                              conjugation_rules, usage_examples, nuance_notes,
		                              jlpt_level, related_patterns, common_mistakes, index_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query, pattern.Pattern, pattern.PlainForm, pattern.Meaning,
		pattern.DetailedExplanation, pattern.ConjugationRules, pattern.UsageExamples,
		pattern.NuanceNotes, pattern.JLPTLevel, pattern.RelatedPatterns,
		pattern.CommonMistakes, pattern.IndexPosition,
	).Scan(&pattern.ID, &pattern.CreatedAt)
}

func (r *GrammarRepository) BulkCreate(patterns []models.GrammarPattern) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO grammar_patterns (pattern, plain_form, meaning, detailed_explanation,
		                              conjugation_rules, usage_examples, nuance_notes,
		                              jlpt_level, related_patterns, common_mistakes, index_position)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, p := range patterns {
		_, err := stmt.Exec(p.Pattern, p.PlainForm, p.Meaning, p.DetailedExplanation,
			p.ConjugationRules, p.UsageExamples, p.NuanceNotes, p.JLPTLevel,
			p.RelatedPatterns, p.CommonMistakes, p.IndexPosition)
		if err != nil {
			return fmt.Errorf("failed to insert pattern %s: %w", p.Pattern, err)
		}
	}

	return tx.Commit()
}
