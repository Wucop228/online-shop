package http

import (
	"database/sql"
	"github.com/Wucop228/online-shop/auth-service/internal/models"
	"github.com/Wucop228/online-shop/auth-service/internal/repository"
	"github.com/Wucop228/online-shop/pkg/hash"
	"github.com/labstack/echo/v4"
	"net/http"
)

type AuthHandler struct {
	db            *sql.DB
	JWTSecret     []byte
	RefreshSecret []byte
}

func NewAuthHandler(db *sql.DB, jwtSecret string, refreshSecret string) *AuthHandler {
	return &AuthHandler{
		db:            db,
		JWTSecret:     []byte(jwtSecret),
		RefreshSecret: []byte(refreshSecret),
	}
}

type Request struct {
	Username string
	Password string
	Email    string
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req Request
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	hash, err := hash.HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": err.Error()})
	}

	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
	}
	if err := repository.CreateUser(h.db, &user); err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{"status": "ok"})
}
