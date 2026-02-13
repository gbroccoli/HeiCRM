package handler

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   uint   `json:"role_id"`
	TgSend   *bool  `json:"tg_send"`
}

// hashToken creates SHA256 hash of the token for use as Redis key
// This prevents storing full JWT tokens in Redis keys and supports multiple sessions
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// isRequestHTTPS checks if the original request used HTTPS
// Handles requests through reverse proxies and tunnels (ngrok, cloudflare)
// that set X-Forwarded-Proto header
func isRequestHTTPS(c *gin.Context) bool {
	// Check X-Forwarded-Proto header (set by reverse proxies/tunnels)
	if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	// Check direct TLS connection
	if c.Request.TLS != nil {
		return true
	}
	return false
}

func getUser(db *sql.DB, email, password string) (*User, error) {
	user := &User{}

	query := `SELECT id, name, email, password, role_id, tg_send FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.TgSend)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) Login(c *gin.Context) {

	ctx := context.Background()

	var userParams models.LoginRequest
	if err := c.ShouldBindJSON(&userParams); err != nil {
		response.ValidationError(c, "Некорректные данные запроса")
		return
	}

	user, err := getUser(h.DB, userParams.Email, userParams.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": response.AuthRequired,
			"msg":  "Неверный email или пароль",
		})
		log.Print(err.Error())
		return
	}

	checkPassword := h.PasswordManager.CheckHash(user.Password, userParams.Password)
	if !checkPassword {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": response.AuthRequired,
			"msg":  "Неверный email или пароль",
		})
		log.Print("invalid email or password")
		return
	}

	tokenAccess, err := h.JWT.GenerateAccessToken(userParams.Email, user.RoleID, "access")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "Непредвиденная ошибка",
		})
		log.Printf("Failed to generate access token: %v", err)
		return
	}

	refreshToken, err := h.JWT.GenerateRefreshToken(userParams.Email, user.RoleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "Непредвиденная ошибка",
		})
		log.Printf("Failed to generate refresh token: %v", err)
		return
	}

	// Store refresh token in Redis with hash-based key
	// Key format: refresh_token:{sha256_hash} enables multiple sessions per user
	// TTL matches JWT expiry (30 days) - sliding window on each /refresh
	tokenHash := hashToken(refreshToken.Token)
	redisKey := fmt.Sprintf("refresh_token:%s", tokenHash)

	err = h.R.Set(
		ctx,
		redisKey,
		refreshToken.Token,
		time.Until(refreshToken.ExpiresAt),
	).Err()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.InternalError,
			"msg":  "Не удалось сохранить сессию",
		})
		log.Printf("Failed to save refresh token to Redis: %v", err)
		return
	}

	// Detect HTTPS - check both env and actual request protocol
	// This handles tunnels (ngrok, cloudflare) that use HTTPS even in dev
	cfg := config.G()
	isProduction := cfg.Env == "production" || cfg.Env == "prod"
	isHTTPS := isProduction || isRequestHTTPS(c)

	// SameSite strategy:
	// - Lax: works for most cases, allows cross-site GET (good for dev)
	// - None: required for cross-origin requests with credentials (needs Secure=true)
	if isHTTPS {
		c.SetSameSite(http.SameSiteNoneMode) // Requires Secure=true (HTTPS)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // Works without HTTPS in dev
	}

	// Domain for cross-subdomain cookies (e.g., ".yourdomain.com")
	// Leave empty for same-domain cookies
	cookieDomain := cfg.Cookie.Domain

	// Delete old cookies with incorrect path (cleanup from previous versions)
	// This prevents duplicate cookies when path was changed from /api/v1/auth to /
	c.SetCookie(
		"refresh",
		"",
		-1,             // MaxAge -1 deletes the cookie
		"/api/v1/auth", // Old path
		cookieDomain,
		isHTTPS,
		true,
	)

	// Set new refresh token cookie
	c.SetCookie(
		"refresh",
		refreshToken.Token,
		int(time.Until(refreshToken.ExpiresAt).Seconds()),
		"/", // Path: "/" allows cookie to be sent to all endpoints
		cookieDomain,
		isHTTPS, // secure: true for HTTPS (production or tunnels)
		true,    // httpOnly: always true for security
	)

	c.JSON(http.StatusOK, gin.H{
		"code":  response.OK,
		"msg":   "Успешный вход",
		"token": tokenAccess.Token,
	})
}
