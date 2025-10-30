package handler

import (
	"log"
	"net/http"

	"github.com/gbroccoli/HeiCRM/services/auth/internal/tools"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TgSend   *bool  `json:"tg_send"`
}

// IsTgSend validate params on null
func (register *RegisterRequest) IsTgSend() bool {
	if register.TgSend == nil {
		return false
	}

	return *register.TgSend
}

func (h *Handler) Register(c *gin.Context) {

	tokenH, err := tools.ExtractToken(c)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	_, err = h.JWT.Verify(tokenH)
	if err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}

	var candidate RegisterRequest
	if err := c.ShouldBindJSON(&candidate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректные данные запроса."})
		return
	}

	token, err := h.JWT.GenerateAccessToken("vovo.r", 2)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"data": gin.H{
			"name":    candidate.Name,
			"tg_send": candidate.IsTgSend(),
			"token":   token,
		},
	})
}
