package natsx

import (
	"fmt"
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/config"
	"github.com/nats-io/nats.go"
)

var conn *nats.Conn

// GetAddr formats NATS address from config
func GetAddr() string {
	cfg := config.G()
	return fmt.Sprintf("nats://%s:%s", cfg.NATS.Host, cfg.NATS.Port)
}

// Open connects to NATS server
func Open() {
	addr := GetAddr()

	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalf("failed to connect to nats: %v", err)
	}

	conn = nc
	log.Printf("connected to nats at %s", addr)
}

// Close closes the NATS connection
func Close() {
	if conn != nil {
		conn.Close()
	}
}

// G returns the global NATS connection (singleton pattern)
func G() *nats.Conn {
	if conn == nil {
		log.Panic("nats connection is nil")
	}
	return conn
}

// Conn returns the global NATS connection (alias for G)
func Conn() *nats.Conn {
	return G()
}
