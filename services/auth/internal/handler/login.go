package handler

import (
	"log"
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

	tokenAccess, err := h.JWT.GenerateAccessToken(userParams.Email, 1, "access")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		log.Fatal(err.Error())
		return
	}

	c.Header("Authorization", "Bearer "+tokenAccess)

	c.JSON(http.StatusOK, gin.H{
		"message": "success login",
		"data": gin.H{
			"email":    userParams.Email,
			"password": userParams.Password,
		},
	})
}
