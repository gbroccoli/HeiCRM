package routes

import (
	"github.com/gbroccoli/HeiCRM/services/auth/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/midleware"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)

		middleGroup := auth.Group("/")
		middleGroup.Use(midleware.AuthMiddleware(h.JWT))
		{
			middleGroup.POST("/register", h.Register)
		}
	}
}
