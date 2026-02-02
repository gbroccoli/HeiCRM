package midleware

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	_ "github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/tools"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// hashToken creates SHA256 hash for Redis key lookups
func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func AuthMiddleware(j *jwt.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Отсутствует заголовок авторизации",
			})
			c.Abort()
			return
		}

		token, err := tools.ExtractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Недействительный токен",
			})
			c.Abort()
			return
		}

		claims, err := j.VerifyAccessToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Subject)
		c.Set("role", claims.Role)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// RefreshTokenMiddleware validates refresh tokens with two-layer security:
// 1. JWT signature verification (cryptographic validity)
// 2. Redis presence check (not revoked via rotation)
// This prevents reuse of rotated tokens and enables token revocation
func RefreshTokenMiddleware(j *jwt.JWT, r *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()

		// Get all cookies named "refresh" (there might be duplicates with different paths)
		// This handles migration from /api/v1/auth to / path
		cookies := c.Request.Cookies()
		var refreshTokens []string
		for _, cookie := range cookies {
			if cookie.Name == "refresh" && cookie.Value != "" {
				refreshTokens = append(refreshTokens, cookie.Value)
			}
		}

		if len(refreshTokens) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": "Refresh-токен отсутствует",
			})
			return
		}

		// Try each refresh token until we find a valid one
		// This handles the case where old cookies with wrong path still exist
		var validClaims *jwt.FieldsClaims
		var validToken string

		for _, token := range refreshTokens {
			// Verify JWT signature and expiry
			claims, err := j.VerifyRefreshToken(token)
			if err != nil {
				continue // Try next token
			}

			// Check if token exists in Redis (not revoked)
			tokenHash := hashToken(token)
			redisKey := fmt.Sprintf("refresh_token:%s", tokenHash)

			exists, err := r.Exists(ctx, redisKey).Result()
			if err != nil {
				continue // Try next token
			}

			if exists > 0 {
				// Found valid token!
				validClaims = claims
				validToken = token
				break
			}
		}

		if validClaims == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": "Токен отозван или недействителен",
			})
			return
		}

		// Store valid token for the handler to use (for rotation)
		c.Set("validRefreshToken", validToken)
		c.Set("email", validClaims.Subject)
		c.Set("userID", validClaims.ID)
		c.Set("role", validClaims.Role)
		c.Next()
	}
}
