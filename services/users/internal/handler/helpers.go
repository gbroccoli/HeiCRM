package handler

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gin-gonic/gin"
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

// logActivity записывает действие пользователя в user_activity_log.
// При ошибке логирует в stdout, не возвращает ошибку клиенту.
func logActivity(db *sql.DB, userID uint64, action string, details interface{}, c *gin.Context) {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		log.Printf("logActivity: failed to marshal details: %v", err)
		return
	}

	_, err = db.Exec(
		`INSERT INTO user_activity_log (user_id, action, details, ip_address, user_agent) VALUES ($1, $2, $3, $4, $5)`,
		userID, action, detailsJSON, c.ClientIP(), c.Request.UserAgent(),
	)
	if err != nil {
		log.Printf("logActivity: failed to insert activity for user %d: %v", userID, err)
	}
}
