package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetRole(c *gin.Context) {
	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		c.Abort()
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"role": role,
	})
}
