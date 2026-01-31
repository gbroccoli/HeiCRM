package routes

import (
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/services/users/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/users/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler, j *jwt.JWT) {
	users := r.Group("/")
	users.Use(middleware.AuthMiddleware(j))
	{
		// /me routes must be registered BEFORE /:id to avoid Gin interpreting "me" as an id
		users.GET("/me", h.GetMe)
		users.PUT("/me", h.UpdateMe)

		// Admin/staff routes
		users.GET("/", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.ListUsers)

		// User by ID routes
		// GetUser handles self-access check internally (admin/manager or own profile)
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", middleware.RoleMiddleware(middleware.RoleAdmin), h.UpdateUser)
		users.DELETE("/:id", middleware.RoleMiddleware(middleware.RoleAdmin), h.DeleteUser)
		users.GET("/:id/activity", middleware.RoleMiddleware(middleware.RoleAdmin), h.GetActivity)
	}
}
