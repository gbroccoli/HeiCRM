package routes

import (
	"github.com/gbroccoli/HeiCRM/services/apigateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

// ServiceConfig holds the configuration for each microservice
type ServiceConfig struct {
	AuthServiceURL string
	// Add more services here as they are created:
	// UserServiceURL string
	// TicketServiceURL string
}

// Mount configures all API Gateway routes
func Mount(r *gin.Engine, config ServiceConfig) {
	// API v1 group
	api := r.Group("/api/v1")
	{
		// Auth service routes - proxy to auth microservice
		// All /api/v1/auth/* requests are forwarded to auth service
		auth := api.Group("/auth")
		{
			auth.Any("/*path", proxy.ReverseProxy(config.AuthServiceURL))
		}

		// Health check endpoint for API Gateway itself
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "api-gateway",
			})
		})

		// Future microservices will be added here:
		// users := api.Group("/users")
		// {
		//     users.Any("/*path", proxy.ReverseProxy(config.UserServiceURL))
		// }
		//
		// tickets := api.Group("/tickets")
		// {
		//     tickets.Any("/*path", proxy.ReverseProxy(config.TicketServiceURL))
		// }
	}
}
