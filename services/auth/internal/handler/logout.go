package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// hashToken creates SHA256 hash of the token for Redis key lookup
func hashTokenLogout(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// isRequestHTTPSLogout checks if the original request used HTTPS
func isRequestHTTPSLogout(c *gin.Context) bool {
	if proto := c.GetHeader("X-Forwarded-Proto"); proto == "https" {
		return true
	}
	if c.Request.TLS != nil {
		return true
	}
	return false
}

// Logout invalidates the user's session by:
// 1. Deleting refresh token from Redis (revokes session)
// 2. Clearing refresh token cookie
func (h *Handler) Logout(c *gin.Context) {
	ctx := context.Background()

	// Get all cookies named "refresh" (handles multiple cookies with different paths)
	cookies := c.Request.Cookies()
	var refreshTokens []string
	for _, cookie := range cookies {
		if cookie.Name == "refresh" && cookie.Value != "" {
			refreshTokens = append(refreshTokens, cookie.Value)
		}
	}

	// Delete all refresh tokens from Redis
	deletedCount := 0
	for _, token := range refreshTokens {
		tokenHash := hashTokenLogout(token)
		redisKey := fmt.Sprintf("refresh_token:%s", tokenHash)

		err := h.R.Del(ctx, redisKey).Err()
		if err != nil {
			log.Printf("Failed to delete refresh token from Redis: %v", err)
			// Continue to delete cookie even if Redis fails
		} else {
			deletedCount++
		}
	}

	// Clear refresh token cookies (delete cookies with both old and new paths)
	cfg := config.G()
	isProduction := cfg.Env == "production" || cfg.Env == "prod"
	isHTTPS := isProduction || isRequestHTTPSLogout(c)

	if isHTTPS {
		c.SetSameSite(http.SameSiteNoneMode)
	} else {
		c.SetSameSite(http.SameSiteLaxMode)
	}

	cookieDomain := cfg.Cookie.Domain

	// Delete cookie with new path (/)
	c.SetCookie(
		"refresh",
		"",
		-1, // MaxAge -1 deletes the cookie
		"/",
		cookieDomain,
		isHTTPS,
		true,
	)

	// Delete cookie with old path (/api/v1/auth) for backwards compatibility
	c.SetCookie(
		"refresh",
		"",
		-1,
		"/api/v1/auth",
		cookieDomain,
		isHTTPS,
		true,
	)

	log.Printf("User logged out successfully. Deleted %d tokens from Redis", deletedCount)

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"msg":  "Успешный выход",
	})
}
