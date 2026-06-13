package db

import (
	"github.com/user/taskapp/internal/models"
)

func CreateUser(user *models.User) error {
	query := `
		INSERT INTO users (id, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := DB.Exec(query, user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	return err
}

func GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE email = $1`
	err := DB.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func GetUserByID(id string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, email, password, created_at, updated_at FROM users WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}
