package handler

import (
	"net/http"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RefreshToken(c *gin.Context) {
	email, ok := c.Get("email")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
	}

	role, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
		return
	}

	newAccessToken, err := h.JWT.GenerateAccessToken(email.(string), role.(uint), "access")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
		return
	}

	c.SetCookie(
		"refresh",
		"",
		-1,
		"/auth/refresh",
		"",
		true,
		true,
	)

	newRefreshToken, err := h.JWT.GenerateRefreshToken(email.(string), role.(uint))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": response.InvalidToken,
			"msg":  "No token",
		})
		return
	}

	expires := time.Now().Add(30 * 24 * time.Hour)
	c.SetCookie(
		"refresh",
		newRefreshToken,
		int(expires.Unix()),
		"/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(200, gin.H{
		"code":  200,
		"token": newAccessToken,
	})
}
