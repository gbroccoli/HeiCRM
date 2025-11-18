package handler

import (
	"net/http"

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
	}

	newAccessToken, err := h.JWT.GenerateAccessToken(email.(string), role.(uint), "access")

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "No token",
		})
	}

	c.Header("Authorization", "Bearer "+newAccessToken)

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "OK",
	})
}
