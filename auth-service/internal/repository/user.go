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

func LoginUser(db *sql.DB, username string) (models.User, error) {
	query := "SELECT *  FROM users WHERE username = $1"
	user := models.User{}
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	if err != nil {
	}
	return user, err
}

func GetUserByRefreshToken(db *sql.DB, refreshToken string) (models.User, error) {
	query := `
		SELECT users.*
		FROM users
		JOIN refresh_tokens ON users.id = refresh_tokens.user_id
		WHERE refresh_tokens.token = $1;
	`
	user := models.User{}
	err := db.QueryRow(query, refreshToken).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash)
	return user, err
}
