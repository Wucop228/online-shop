package app

import (
	"database/sql"
	"fmt"
	"github.com/Wucop228/online-shop/auth-service/internal/config"
	"github.com/Wucop228/online-shop/auth-service/internal/delivery/http"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	"log"
)

func Run(cfg *config.Config) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s port=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.DBPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB connection error: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("DB ping error: ", err)
	}

	authCfg := config.AuthConfig{AccessTokenTTL: cfg.AccessTokenTTL, RefreshTokenTTL: cfg.RefreshTokenTTL,
		JWTSecret: cfg.JWTSecret, RefreshSecret: cfg.RefreshSecret}

	e := echo.New()
	h := http.NewAuthHandler(db, &authCfg)

	e.POST("/create-user", h.Register)
	e.POST("/login", h.Login)
	e.POST("/refresh", h.RefreshToken)

	if err := e.Start(fmt.Sprintf(":%s", cfg.ServerPort)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
