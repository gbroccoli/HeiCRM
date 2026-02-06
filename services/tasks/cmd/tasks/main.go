package main

import (
	"log"
	"os"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
	"github.com/gbroccoli/HeiCRM/services/tasks/internal/handler"
	"github.com/gbroccoli/HeiCRM/services/tasks/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// init logger
	err := logx.Init("logs/tasks.log")
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}

	// read config
	log.Println("starting tasks service")
	config.MustLoad("config.yaml")

	// get pid process for find in system monitor
	log.Printf("PID=%d", os.Getpid())

	// create default gin engine
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	// init connect to database
	log.Println("Connecting to database")
	dbx.Open()
	defer func() {
		err := dbx.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	// init jwt
	j := jwt.New([]byte(config.G().Jwt.SecretKey))

	// init handler
	h := handler.New(dbx.G(), j)

	// mount routes
	routes.Mount(g, h, j)

	// run api server
	log.Println("starting http server on :8083")
	err = g.Run(":8083")
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
