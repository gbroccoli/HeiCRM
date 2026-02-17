package main

import (
	"log"
	"os"
	"time"

	apigatewayConf "github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
	"github.com/gbroccoli/HeiCRM/services/apigateway/internal/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	err := logx.Init("logs/apigateway.log")
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}

	apigatewayConf.MustLoad("config.yaml")

	conf := apigatewayConf.G()

	log.Println("Starting API Gateway")
	log.Printf("PID=%d", os.Getpid())

	// Create Gin router
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	// CORS middleware - centralized for all microservices
	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true // Allow all origins in development
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	authUrl := conf.Serves.Auth
	if authUrl == "" {
		authUrl = "http://localhost:8080"
	}

	usersUrl := conf.Serves.Users
	if usersUrl == "" {
		usersUrl = "http://localhost:8081"
	}

	housingUrl := conf.Serves.Housing
	if housingUrl == "" {
		housingUrl = "http://localhost:8082"
	}

	tasksUrl := conf.Serves.Tasks
	if tasksUrl == "" {
		tasksUrl = "http://localhost:8083"
	}

	notificationUrl := conf.Serves.Notification
	if notificationUrl == "" {
		notificationUrl = "http://localhost:8084"
	}

	// Service configuration - microservice URLs
	serviceConfig := routes.ServiceConfig{
		AuthServiceURL:         authUrl,
		UserServiceURL:         usersUrl,
		HousingServiceURL:      housingUrl,
		TasksServiceURL:        tasksUrl,
		NotificationServiceURL: notificationUrl,
	}

	// Mount all routes
	routes.Mount(r, serviceConfig)

	// Start server
	port := getEnv("GATEWAY_PORT", "8000")
	log.Printf("API Gateway listening on :%s", port)
	log.Printf("Proxying auth service: %s", serviceConfig.AuthServiceURL)
	log.Printf("Proxying users service: %s", serviceConfig.UserServiceURL)
	log.Printf("Proxying housing service: %s", serviceConfig.HousingServiceURL)
	log.Printf("Proxying tasks service: %s", serviceConfig.TasksServiceURL)
	log.Printf("Proxying notification service: %s", serviceConfig.NotificationServiceURL)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
