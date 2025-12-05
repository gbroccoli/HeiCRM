package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// RefreshToken implements token rotation with sliding window strategy:
// 1. Validates old refresh token (JWT signature + Redis presence)
// 2. Deletes old token from Redis (rotation - prevents reuse)
// 3. Issues new access + refresh tokens with extended expiry (sliding window)
// Result: Active users never re-login, stolen tokens work only once
func (h *Handler) RefreshToken(c *gin.Context) {
	ctx := context.Background()

	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
		return
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
		return
	}

	// Get old refresh token from cookie (for deletion from Redis)
	oldRefreshToken, err := c.Cookie("refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.InvalidToken,
			"msg":  "No refresh token",
		})
		return
	}

	// Generate new access token
	newAccessToken, err := h.JWT.GenerateAccessToken(email.(string), role.(uint), "access")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "Failed to generate access token",
		})
		log.Printf("Failed to generate access token: %v", err)
		return
	}

	// Generate new refresh token
	newRefreshToken, err := h.JWT.GenerateRefreshToken(email.(string), role.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": response.TokenGenerationFailed,
			"msg":  "Failed to generate refresh token",
		})
		log.Printf("Failed to generate refresh token: %v", err)
		return
	}

	// Delete old refresh token from Redis (rotation strategy)
	oldTokenHash := hashToken(oldRefreshToken)
	oldRedisKey := fmt.Sprintf("refresh_token:%s", oldTokenHash)

	err = h.R.Del(ctx, oldRedisKey).Err()
	if err != nil {
		log.Printf("Warning: Failed to delete old refresh token from Redis: %v", err)
		// Continue anyway - old token will expire naturally
	}

	// Save new refresh token to Redis with sliding window
	// Each /refresh call generates a new token with ExpiresAt = now + 30 days
	// This means active users never need to re-login (sliding window strategy)
	newTokenHash := hashToken(newRefreshToken.Token)
	newRedisKey := fmt.Sprintf("refresh_token:%s", newTokenHash)

	err = h.R.Set(
		ctx,
		newRedisKey,
		newRefreshToken.Token,
		time.Until(newRefreshToken.ExpiresAt), // Always ~30 days from now
	).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": response.InternalError,
			"msg":  "Failed to save session",
		})
		log.Printf("Failed to save new refresh token to Redis: %v", err)
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

	// Set new refresh token cookie
	c.SetCookie(
		"refresh",
		newRefreshToken.Token,
		int(time.Until(newRefreshToken.ExpiresAt).Seconds()),
		"/api/v1/auth",
		cookieDomain,
		isProduction, // secure: true only in production
		true,         // httpOnly: always true for security
	)

	c.JSON(http.StatusOK, gin.H{
		"code":  response.OK,
		"token": newAccessToken.Token,
	})
}
