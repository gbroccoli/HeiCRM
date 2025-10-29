package handler

import (
	"net/http"

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

	c.JSON(http.StatusOK, gin.H{
		"message": "success login",
	})
}
