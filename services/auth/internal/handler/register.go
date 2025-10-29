package handler

import (
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TgSend   *bool  `json:"tg_send"`
}

func (h *Handler) Register(c *gin.Context) {
	var candidate RegisterRequest
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса."})
		return
	}

	var db = dbx.G()

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"name":    candidate.Name,
			"tg_send": candidate.TgSend,
		},
	})
}
