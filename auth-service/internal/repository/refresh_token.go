package repository

import (
	"database/sql"
	"github.com/Wucop228/online-shop/auth-service/internal/models"
)

func CreateRefreshToken(db *sql.DB, refreshToken models.RefreshToken) error {
	query := "INSERT INTO refresh_tokens (token, user_id, expiry) values ($1, $2, $3)"
	_, err := db.Exec(query, refreshToken.Token, refreshToken.UserID, refreshToken.Expiry)
	return err
}
