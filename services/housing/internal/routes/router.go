package routes

import (
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/services/housing/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/housing/internal/middleware"
	"github.com/gin-gonic/gin"
)

func Mount(r *gin.Engine, h *handler.Handler, j *jwt.JWT) {
	buildings := r.Group("/")
	buildings.Use(middleware.AuthMiddleware(j))
	{
		// Buildings CRUD
		buildings.GET("/", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.ListBuildings)
		buildings.POST("/", middleware.RoleMiddleware(middleware.RoleAdmin), h.CreateBuilding)
		buildings.GET("/:id", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.GetBuilding)
		buildings.PUT("/:id", middleware.RoleMiddleware(middleware.RoleAdmin), h.UpdateBuilding)
		buildings.DELETE("/:id", middleware.RoleMiddleware(middleware.RoleAdmin), h.DeleteBuilding)

		// Rooms CRUD
		buildings.GET("/:id/rooms", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.ListRooms)
		buildings.POST("/:id/rooms", middleware.RoleMiddleware(middleware.RoleAdmin), h.CreateRoom)
		buildings.GET("/:id/rooms/:roomId", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.GetRoom)
		buildings.PUT("/:id/rooms/:roomId", middleware.RoleMiddleware(middleware.RoleAdmin), h.UpdateRoom)
		buildings.DELETE("/:id/rooms/:roomId", middleware.RoleMiddleware(middleware.RoleAdmin), h.DeleteRoom)

		// Residents
		buildings.POST("/:id/rooms/:roomId/residents", middleware.RoleMiddleware(middleware.RoleAdmin), h.AssignResident)
		buildings.GET("/:id/rooms/:roomId/residents", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.ListResidents)
		buildings.GET("/:id/rooms/:roomId/residents/:residentId", middleware.RoleMiddleware(middleware.RoleAdmin, middleware.RoleManager), h.GetResident)
		buildings.PUT("/:id/rooms/:roomId/residents/:residentId", middleware.RoleMiddleware(middleware.RoleAdmin), h.UpdateResident)
		buildings.POST("/:id/rooms/:roomId/residents/:residentId/transfer", middleware.RoleMiddleware(middleware.RoleAdmin), h.TransferResident)
		buildings.DELETE("/:id/rooms/:roomId/residents/:residentId", middleware.RoleMiddleware(middleware.RoleAdmin), h.RemoveResident)
	}
}
