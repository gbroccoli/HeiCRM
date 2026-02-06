package routes

import (
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/services/tasks/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/tasks/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler, j *jwt.JWT) {
	tasks := r.Group("/")
	tasks.Use(middleware.AuthMiddleware(j))
	{
		// Tasks CRUD
		tasks.GET("/", h.ListTasks)
		tasks.POST("/", h.CreateTask)
		tasks.GET("/:id", h.GetTask)
		tasks.PUT("/:id", h.UpdateTask)
		tasks.DELETE("/:id", middleware.RoleMiddleware(middleware.RoleAdmin), h.DeleteTask)

		// Task status management
		tasks.PUT("/:id/status", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.UpdateTaskStatus)

		// Task assignment
		tasks.PUT("/:id/assign", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.AssignTask)
		tasks.POST("/:id/take", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.TakeTask)

		// Comments
		tasks.GET("/:id/comments", h.GetComments)
		tasks.POST("/:id/comments", h.AddComment)

		// History
		tasks.GET("/:id/history", h.GetHistory)
	}
}
