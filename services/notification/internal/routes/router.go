package routes

import (
	"github.com/gin-gonic/gin"
)

// Mount configures notification service routes
func Mount(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "notification",
		})
	})
}
