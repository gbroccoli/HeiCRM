package routes

import (
	"github.com/gbroccoli/HeiCRM/services/auth/internal/handler"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/register", h.Register)
	}
}
