package handler

import "database/sql"

type Handler struct {
	DB *sql.DB
}

func New(DB *sql.DB) *Handler {
	return &Handler{
		DB: DB,
	}
}
