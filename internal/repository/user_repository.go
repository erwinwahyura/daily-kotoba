package repository

import (
	"database/sql"
	"errors"

	"github.com/erwinwahyura/daily-kotoba/internal/db"
	"github.com/erwinwahyura/daily-kotoba/internal/models"
)

type UserRepository struct {
	db *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	// Generate UUID in application code for SQLite compatibility
	if user.ID == "" {
		user.ID = r.db.GenerateUUID()
	}
	
	query := `
		INSERT INTO users (id, email, password_hash, current_jlpt_level)
		VALUES (` + r.db.Placeholder(1) + `, ` + r.db.Placeholder(2) + `, ` + r.db.Placeholder(3) + `, ` + r.db.Placeholder(4) + `)
		RETURNING created_at, updated_at
	`
	err := r.db.QueryRow(query, user.ID, user.Email, user.PasswordHash, user.CurrentLevel).
		Scan(&user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	// Create initial user progress
	progressQuery := `
		INSERT INTO user_progress (user_id, current_vocab_index)
		VALUES (` + r.db.Placeholder(1) + `, 0)
	`
	_, err = r.db.Exec(progressQuery, user.ID)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, current_jlpt_level, created_at, updated_at
		FROM users
		WHERE email = ` + r.db.Placeholder(1) + `
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CurrentLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, current_jlpt_level, created_at, updated_at
		FROM users
		WHERE id = ` + r.db.Placeholder(1) + `
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CurrentLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) EmailExists(email string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`
	err := r.db.QueryRow(query, email).Scan(&exists)
	return exists, err
}

func (r *UserRepository) UpdateUserLevel(userID, level string) error {
	query := `
		UPDATE users
		SET current_jlpt_level = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(query, level, userID)
	return err
}
