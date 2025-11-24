package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/jwt"
	"github.com/gbroccoli/HeiCRM/pkg/password"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	DB              *sql.DB
	JWT             *jwt.JWT
	PasswordManager *password.Password
	R               *redis.Client
}

func New(DB *sql.DB, JWT *jwt.JWT, r *redis.Client) *Handler {
	return &Handler{
		DB:              DB,
		JWT:             JWT,
		PasswordManager: &password.Password{},
		R:               r,
	}
}
