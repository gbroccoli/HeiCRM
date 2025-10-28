package main

import (
	"log"
	"os"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
	"github.com/gin-gonic/gin"
)

func main() {
	err := logx.Init("logs/auth.log")
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}

	log.Println("starting auth service")
	config.MustLoad("config.yaml")

	log.Printf("PID=%d", os.Getpid())

	g := gin.Default()

	log.Println("Connecting to database")
	dbx.Open()
	defer func() {
		err := dbx.Close()
		if err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	log.Println("starting http server")
	err = g.Run(":8080")
	if err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
