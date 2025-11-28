package routes

import (
	"github.com/gbroccoli/HeiCRM/services/auth/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/midleware"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler) {
	auth := r.Group("/")
	{
		auth.POST("/login", h.Login)

		auth.POST("/refresh", midleware.RefreshTokenMiddleware(h.JWT, h.R), h.RefreshToken)

		middleGroup := auth.Group("/")
		middleGroup.Use(midleware.AuthMiddleware(h.JWT))
		{
			middleGroup.POST("/register", h.Register)
			middleGroup.POST("/logout", h.Logout)
			middleGroup.GET("/me", h.MeUsers)
			middleGroup.GET("/role")
		}
	}
}
