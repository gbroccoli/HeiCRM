package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
)

type Handler struct {
	DB  *sql.DB
	JWT *jwt.JWT
}

func New(db *sql.DB, jwt *jwt.JWT) *Handler {
	return &Handler{
		DB:  db,
		JWT: jwt,
	}
}
