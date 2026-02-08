package main

import (
	"log"
	"os"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/gbroccoli/HeiCRM/pkg/dbx"
	"github.com/gbroccoli/HeiCRM/pkg/logx"
	"github.com/gbroccoli/HeiCRM/pkg/natsx"
	"github.com/gbroccoli/HeiCRM/services/notification/internal/email"
	natsub "github.com/gbroccoli/HeiCRM/services/notification/internal/nats"
	"github.com/gbroccoli/HeiCRM/services/notification/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// init logger
	err := logx.Init("logs/notification.log")
	if err != nil {
		log.Fatalf("failed to init log: %v", err)
	}

	// read config
	log.Println("starting notification service")
	config.MustLoad("config.yaml")

	// get pid
	log.Printf("PID=%d", os.Getpid())

	// create gin engine
	g := gin.Default()
	g.Use(gin.Recovery())
	g.Use(gin.Logger())

	// init database
	log.Println("Connecting to database")
	dbx.Open()
	defer func() {
		if err := dbx.Close(); err != nil {
			log.Fatalf("failed to close db: %v", err)
		}
	}()

	// init nats
	log.Println("Connecting to NATS")
	natsx.Open()
	defer natsx.Close()

	// init email sender
	sender := email.NewSender()

	// start NATS subscribers
	sub := natsub.NewSubscriber(natsx.G(), dbx.G(), sender)
	if err := sub.Subscribe(); err != nil {
		log.Fatalf("failed to subscribe to NATS: %v", err)
	}

	// mount routes
	routes.Mount(g)

	// run http server
	log.Println("starting http server on :8084")
	if err := g.Run(":8084"); err != nil {
		log.Fatalf("failed to start http server: %v", err)
	}
}
