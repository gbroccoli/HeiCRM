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
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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

func getUser(db *sql.DB, email, password string) (*User, error) {
	user := &User{}

	query := `SELECT id, name, email, password, role_id, tg_send FROM users WHERE email = $1`
	err := db.QueryRow(query, email).Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.RoleID, &user.TgSend)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (h *Handler) Login(c *gin.Context) {

	ctx := context.Background()

	var userParams LoginRequest
	if err := c.ShouldBindJSON(&userParams); err != nil {
		response.ValidationError(c, "Invalid request data")
		return
	}

	user, err := getUser(h.DB, userParams.Email, userParams.Password)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code": response.AuthRequired,
			"msg":  "invalid email or password",
		})
		log.Print(err.Error())
		return
	}

	checkPassword := h.PasswordManager.CheckHash(user.Password, userParams.Password)
	if !checkPassword {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code": response.AuthRequired,
			"msg":  "invalid email or password",
		})
		log.Print("invalid email or password")
		return
	}

	tokenAccess, err := h.JWT.GenerateAccessToken(userParams.Email, 1, "access")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "an unexpected problem",
		})
		log.Printf("Failed to generate access token: %v", err)
		return
	}

	refreshToken, err := h.JWT.GenerateRefreshToken(userParams.Email, user.RoleID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "an unexpected problem",
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
			"msg":  "failed to save session",
		})
		log.Printf("Failed to save refresh token to Redis: %v", err)
		return
	}

	// Secure cookie only in production (requires HTTPS)
	// For local development (http://localhost), use secure=false
	cfg := config.G()
	isProduction := cfg.Env == "production" || cfg.Env == "prod"

	// SameSite=Lax allows cookie in cross-site GET and same-site POST
	// For production with HTTPS, use SameSite=None for full cross-origin support
	if isProduction {
		c.SetSameSite(http.SameSiteNoneMode) // Requires Secure=true (HTTPS)
	} else {
		c.SetSameSite(http.SameSiteLaxMode) // Works without HTTPS in dev
	}

	// Domain for cross-subdomain cookies (e.g., ".yourdomain.com")
	// Leave empty for same-domain cookies
	cookieDomain := cfg.Cookie.Domain

	c.SetCookie(
		"refresh",
		refreshToken.Token,
		int(time.Until(refreshToken.ExpiresAt).Seconds()),
		"/api/v1/auth",
		cookieDomain,
		isProduction, // secure: true only in production
		true,         // httpOnly: always true for security
	)

	c.JSON(http.StatusOK, gin.H{
		"code":  response.OK,
		"msg":   "successfully logged in",
		"token": tokenAccess.Token,
	})
}
