package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/nats-io/nats.go"
)

type Handler struct {
	DB  *sql.DB
	JWT *jwt.JWT
	NC  *nats.Conn
}

func New(db *sql.DB, jwt *jwt.JWT, nc *nats.Conn) *Handler {
	return &Handler{
		DB:  db,
		JWT: jwt,
		NC:  nc,
	}
}
