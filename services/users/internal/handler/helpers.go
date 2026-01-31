package handler

import (
	"database/sql"

	"github.com/gbroccoli/HeiCRM/pkg/models"
)

// getUserIDByEmail looks up the database user ID by email
func getUserIDByEmail(db *sql.DB, email string) (uint64, error) {
	var id uint64
	err := db.QueryRow("SELECT id FROM users WHERE email = $1", email).Scan(&id)
	return id, err
}

// getUserWithProfile fetches user joined with profile by user ID
func getUserWithProfile(db *sql.DB, userID uint64) (*models.UserWithProfile, error) {
	return models.GetUserWithProfile(db, userID)
}

// getCurrentUserID resolves the current user's database ID from the email stored in gin context
func getCurrentUserID(h *Handler, email string) (uint64, error) {
	return getUserIDByEmail(h.DB, email)
}
