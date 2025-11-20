package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) Login(c *gin.Context) {

	var userParams LoginRequest
	if err := c.ShouldBindJSON(&userParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса."})
		return
	}

	// logic no database users

	tokenAccess, err := h.JWT.GenerateAccessToken(userParams.Email, 1, "access")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
	}

	refreshToken, err := h.JWT.GenerateRefreshToken(userParams.Email)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
	}

	expires := 30 * 24 * time.Hour
	c.SetCookie(
		"refresh",
		refreshToken,
		int(expires),
		"/auth/refresh",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "success login",
		"token":   tokenAccess,
	})
}
