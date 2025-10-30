package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
)

type Handler struct {
	DB  *sql.DB
	JWT *jwt.JWT
}

func New(DB *sql.DB, JWT *jwt.JWT) *Handler {
	return &Handler{
		DB:  DB,
		JWT: JWT,
	}
}
