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
	// Add your production frontend domain before deployment
	cfg := apigatewayConf.G()
	allowedOrigins := []string{
		"http://localhost:3000", // React/Next.js dev server
		"http://localhost:5173", // Vite dev server
		"http://localhost:4200", // Angular dev server
		"http://localhost:8081", // Alternative frontend port
	}

	// Add production origins if not in dev mode
	if cfg.Env == "production" || cfg.Env == "prod" {
		// TODO: Replace with your actual production domains
		allowedOrigins = append(allowedOrigins,
			"https://crm.yourdomain.com", // Production frontend
			"https://app.yourdomain.com", // Alternative domain
		)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true, // Required for HTTPOnly cookies
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

	buildingsUrl := conf.Serves.Buildings
	if buildingsUrl == "" {
		buildingsUrl = "http://localhost:8082"
	}

	// Service configuration - microservice URLs
	serviceConfig := routes.ServiceConfig{
		AuthServiceURL:     authUrl,
		UserServiceURL:     usersUrl,
		BuildingServiceURL: buildingsUrl,
	}

	// Mount all routes
	routes.Mount(r, serviceConfig)

	// Start server
	port := getEnv("GATEWAY_PORT", "8000")
	log.Printf("API Gateway listening on :%s", port)
	log.Printf("Proxying auth service: %s", serviceConfig.AuthServiceURL)
	log.Printf("Proxying users service: %s", serviceConfig.UserServiceURL)
	log.Printf("Proxying buildings service: %s", serviceConfig.BuildingServiceURL)

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
