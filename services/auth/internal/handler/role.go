package handler

import (
	"net/http"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) Role(c *gin.Context) {
	role, ok := c.Get("role")
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Не авторизован",
		})
		return
	}

	user := models.User{RoleID: role.(uint)}
	roleName, err := user.GetRole(h.DB)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Не авторизован",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": response.OK,
		"role": roleName,
	})
}
