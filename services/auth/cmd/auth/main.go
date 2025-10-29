package main

import (
	"log"
	"os"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
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

	log.Printf("PID=%d", os.Getpid())

	// create default
	g := gin.Default()

	log.Println("Connecting to database")
	dbx.Open()
	defer func() {
		err := dbx.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	h := handler.New(dbx.G())

	routes.Mount(g, h)

	log.Println("starting http server")
	err = g.Run(":8080")
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
