package routes

import (
	"github.com/gbroccoli/HeiCRM/services/apigateway/internal/proxy"
	"github.com/gin-gonic/gin"
)

// ServiceConfig holds the configuration for each microservice
type ServiceConfig struct {
	AuthServiceURL         string
	UserServiceURL         string
	HousingServiceURL      string
	TasksServiceURL        string
	NotificationServiceURL string
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

		// User service routes - proxy to users microservice
		users := api.Group("/users")
		{
			users.Any("/*path", proxy.ReverseProxy(config.UserServiceURL))
		}

		// Housing service routes - proxy to housing microservice
		housing := api.Group("/housing")
		{
			housing.Any("/*path", proxy.ReverseProxy(config.HousingServiceURL))
		}

		// Tasks service routes - proxy to tasks microservice
		tasks := api.Group("/tasks")
		{
			tasks.Any("/*path", proxy.ReverseProxy(config.TasksServiceURL))
		}

		// Notification service routes - proxy to notification microservice
		notifications := api.Group("/notifications")
		{
			notifications.Any("/*path", proxy.ReverseProxy(config.NotificationServiceURL))
		}

		// Health check endpoint for API Gateway itself
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "api-gateway",
			})
		})
	}
}
