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
				"error": "missing authorization header",
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
				"error": "invalid token",
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

		c.Set("email", claims.Email)
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

		refreshToken, err := c.Cookie("refresh")
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": "invalid token",
			})
			c.Abort()
			return
		}
		if refreshToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": "missing authorization header",
			})
			c.Abort()
			return
		}

		token := refreshToken
		claims, err := j.VerifyRefreshToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": err.Error(),
			})
			return
		}

		// Check if token exists in Redis (not revoked)
		tokenHash := hashToken(token)
		redisKey := fmt.Sprintf("refresh_token:%s", tokenHash)

		exists, err := r.Exists(ctx, redisKey).Result()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":  response.InternalError,
				"error": "failed to verify token",
			})
			return
		}

		if exists == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.InvalidToken,
				"error": "token has been revoked or is invalid",
			})
			c.Abort()
			return
		}

		c.Set("email", claims.Email)
		c.Set("userID", claims.ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
