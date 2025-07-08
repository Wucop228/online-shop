package http

import (
	"database/sql"
	"github.com/Wucop228/online-shop/auth-service/internal/config"
	"github.com/Wucop228/online-shop/auth-service/internal/models"
	"github.com/Wucop228/online-shop/auth-service/internal/repository"
	"github.com/Wucop228/online-shop/pkg/hash"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type AuthHandler struct {
	db         *sql.DB
	AuthConfig *config.AuthConfig
}

func NewAuthHandler(db *sql.DB, authConfig *config.AuthConfig) *AuthHandler {
	return &AuthHandler{
		db:         db,
		AuthConfig: authConfig,
	}
}

type RequestRegister struct {
	Username string
	Password string
	Email    string
}

type RequestLogin struct {
	Username string
	Password string
}

type RequestRefreshToken struct {
	RefreshToken string
}

func setTokenCookie(c echo.Context, name_token string, token string, expiry time.Time, path string) {
	cookie := new(http.Cookie)
	cookie.Name = name_token
	cookie.Value = token
	cookie.Expires = expiry
	cookie.Path = path
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.SameSite = http.SameSiteLaxMode

	c.SetCookie(cookie)
}

func GenerateAccessToken(h *AuthHandler, c echo.Context, user models.User) (string, error) {
	accessExp := time.Now().Add(h.AuthConfig.AccessTokenTTL).Unix()
	accessClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": accessExp,
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenSigned, err := accessToken.SignedString([]byte(h.AuthConfig.JWTSecret))
	if err != nil {
		return "", c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return accessTokenSigned, err
}

func GenerateRefreshToken(h *AuthHandler, c echo.Context, user models.User) (string, error) {
	refreshExp := time.Now().Add(h.AuthConfig.RefreshTokenTTL).Unix()
	refreshClaims := jwt.MapClaims{
		"sub": user.ID,
		"exp": refreshExp,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenSigned, err := refreshToken.SignedString([]byte(h.AuthConfig.JWTSecret))
	if err != nil {
		return "", c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	refreshModel := models.RefreshToken{
		Token:  refreshTokenSigned,
		UserID: user.ID,
		Expiry: time.Duration(refreshExp),
	}

	if err := repository.CreateRefreshToken(h.db, refreshModel); err != nil {
		return "", c.JSON(http.StatusBadGateway, echo.Map{"error": err.Error()})
	}

	return refreshTokenSigned, nil
}

func (h *AuthHandler) Register(c echo.Context) error {
	var req RequestRegister
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "username or password is empty"})
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

	return c.JSON(http.StatusCreated, echo.Map{
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req RequestLogin
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var user models.User
	user, err := repository.LoginUser(h.db, req.Username)
	if err != nil {
		return c.JSON(http.StatusBadGateway, echo.Map{"error": err.Error()})
	}

	if !hash.CheckPasswordHash(req.Password, user.PasswordHash) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "wrong username or password"})
	}

	accessTokenSigned, err := GenerateAccessToken(h, c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	refreshTokenSigned, err := GenerateRefreshToken(h, c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	setTokenCookie(c, "access_token", accessTokenSigned, time.Now().Add(h.AuthConfig.AccessTokenTTL), "/")
	setTokenCookie(c, "refresh_token", refreshTokenSigned, time.Now().Add(h.AuthConfig.RefreshTokenTTL), "/")

	return c.JSON(http.StatusOK, echo.Map{
		"access_token":  accessTokenSigned,
		"refresh_token": refreshTokenSigned,
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var reqToken RequestRefreshToken
	if err := c.Bind(&reqToken); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if reqToken.RefreshToken == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "refresh token is required"})
	}

	var user models.User
	user, err := repository.GetUserByRefreshToken(h.db, reqToken.RefreshToken)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid refresh token"})
	}

	accessTokenSigned, err := GenerateAccessToken(h, c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	refreshTokenSigned, err := GenerateRefreshToken(h, c, user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	setTokenCookie(c, "access_token", accessTokenSigned, time.Now().Add(h.AuthConfig.AccessTokenTTL), "/")
	setTokenCookie(c, "refresh_token", refreshTokenSigned, time.Now().Add(h.AuthConfig.RefreshTokenTTL), "/")

	return c.JSON(http.StatusCreated, echo.Map{
		"access_token": accessTokenSigned,
	})
}
