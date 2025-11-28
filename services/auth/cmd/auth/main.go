package main

import (
	"log"
	"os"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
	"github.com/gbroccoli/HeiCRM/pkg/redisx"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/auth/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// init logout
	err := logx.Init("logs/auth.log")
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}

	// read config
	log.Println("starting auth service")
	config.MustLoad("config.yaml")

	// get pid process for find in system monitor
	log.Printf("PID=%d", os.Getpid())

	// create default
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	// CORS is now handled by API Gateway
	// Requests should come through: Frontend -> Gateway (:8000) -> Auth Service (:8080)

	// init connect to database
	log.Println("Connecting to database")
	dbx.Open()
	defer func() {
		err := dbx.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	// init redis
	log.Println("Connecting to Redis")
	redisx.Open()
	defer func() {
		if err := redisx.Close(); err != nil {
			log.Printf("failed to close redis: %v", err)
		}
	}()

	// init jwt
	j := jwt.New([]byte(config.G().Jwt.SecretKey))

	// init base model handler
	h := handler.New(dbx.G(), j, redisx.G())

	// mount routers
	routes.Mount(g, h)

	// run api servers
	log.Println("starting http server")
	err = g.Run(":8080")
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
