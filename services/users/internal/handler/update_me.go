package handler

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	natsHandler "github.com/gbroccoli/HeiCRM/services/users/internal/nats"
	"github.com/gin-gonic/gin"
)

// UpdateMe updates the current authenticated user's profile
func (h *Handler) UpdateMe(c *gin.Context) {
	email, exists := c.Get("email")
	if !exists {
		response.Unauthorized(c, "Email не найден в контексте")
		return
	}

	emailStr, ok := email.(string)
	if !ok || emailStr == "" {
		response.Unauthorized(c, "Некорректный email в токене")
		return
	}

	userID, err := getUserIDByEmail(h.DB, emailStr)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "Пользователь не найден")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректное тело запроса", err)
		return
	}

	// Upsert profile
	query := `
		INSERT INTO user_profiles (user_id, first_name, last_name, middle_name, phone, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id)
		DO UPDATE SET
			first_name = COALESCE($2, user_profiles.first_name),
			last_name = COALESCE($3, user_profiles.last_name),
			middle_name = COALESCE($4, user_profiles.middle_name),
			phone = COALESCE($5, user_profiles.phone),
			updated_at = $6
	`

	now := time.Now()
	_, err = h.DB.Exec(query, userID, req.FirstName, req.LastName, req.MiddleName, req.Phone, now)
	if err != nil {
		log.Printf("failed to update profile for user %d: %v", userID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Publish NATS event
	natsHandler.PublishProfileUpdated(h.NC, userID, emailStr)

	// Fetch updated profile
	user, err := getUserWithProfile(h.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "Пользователь не найден")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    response.Updated,
		"message": "Профиль обновлён",
		"data":    user,
	})
}
