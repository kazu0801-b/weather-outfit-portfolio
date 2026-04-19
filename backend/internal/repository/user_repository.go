package repository

import (
	"database/sql"

	"weather-outfit-backend/internal/model"
)

type UserRepository struct {
	DB *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user model.User) error {
	query := `
		INSERT INTO users (username, email, password_hash)
		VALUES (?, ?, ?)
	`

	_, err := r.DB.Exec(query, user.Username, user.Email, user.PasswordHash)
	return err
}

func (r *UserRepository) FindUserByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	var user model.User

	err := r.DB.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}