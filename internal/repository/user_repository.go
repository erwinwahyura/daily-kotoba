package repository

import (
	"database/sql"
	"errors"

	"github.com/yourusername/kotoba-api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (email, password_hash, current_jlpt_level)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(query, user.Email, user.PasswordHash, user.CurrentLevel).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}

	// Create initial user progress
	progressQuery := `
		INSERT INTO user_progress (user_id, current_vocab_index)
		VALUES ($1, 0)
	`
	_, err = r.db.Exec(progressQuery, user.ID)
	return err
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, password_hash, current_jlpt_level, created_at, updated_at
		FROM users
		WHERE email = $1
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
		WHERE id = $1
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
