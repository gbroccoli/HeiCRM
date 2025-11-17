package handler

import (
	"log"
	"net/http"

	"github.com/gbroccoli/HeiCRM/services/auth/internal/tools"
	"github.com/gin-gonic/gin"
)

func (h *Handler) RefreshToken(c *gin.Context) {

	token, err := tools.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
		log.Fatal("Token extract error:", err)
		return
	}

	//newAccessToken := h.JWT.GenerateAccessToken()

	c.JSON(200, gin.H{
		"code": 200,
	})
}
