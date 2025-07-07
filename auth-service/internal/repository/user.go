package repository

import (
	"database/sql"
	"github.com/Wucop228/online-shop/auth-service/internal/models"
)

func CreateUser(db *sql.DB, user *models.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)"
	_, err := db.Exec(query, user.Username, user.Email, user.PasswordHash)
	return err
}
