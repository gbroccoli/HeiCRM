package handler

import "github.com/gin-gonic/gin"

func (h *Handler) MeUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "success",
	})
}
