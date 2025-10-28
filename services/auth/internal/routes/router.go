package routes

import "github.com/gin-gonic/gin"

func Mount(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login")
		auth.POST("/logout")
		auth.POST("/refresh")
		auth.POST("/refresh")
		auth.GET("/me")
	}
}
