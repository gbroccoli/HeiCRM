package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/password"
)

type Handler struct {
	DB              *sql.DB
	JWT             *jwt.JWT
	PasswordManager *password.Password
}

func New(DB *sql.DB, JWT *jwt.JWT) *Handler {
	return &Handler{
		DB:              DB,
		JWT:             JWT,
		PasswordManager: &password.Password{},
	}
}
