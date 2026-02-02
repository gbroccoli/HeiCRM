package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/password"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	DB              *sql.DB
	JWT             *jwt.JWT
	PasswordManager *password.Password
	R               *redis.Client
	NC              *nats.Conn
}

func New(DB *sql.DB, JWT *jwt.JWT, r *redis.Client, nc *nats.Conn) *Handler {
	return &Handler{
		DB:              DB,
		JWT:             JWT,
		PasswordManager: &password.Password{},
		R:               r,
		NC:              nc,
	}
}
