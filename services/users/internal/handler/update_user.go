package handler

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// UpdateUser allows admin to update any user's data
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid user id")
		return
	}

	var req models.AdminUpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "invalid request body", err)
		return
	}

	tx, err := h.DB.Begin()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer tx.Rollback()

	now := time.Now()

	// Update users table if name or role_id provided
	if req.Name != nil || req.RoleID != nil {
		if req.Name != nil {
			_, err = tx.Exec("UPDATE users SET name = $1, updated_at = $2 WHERE id = $3", *req.Name, now, userID)
			if err != nil {
				log.Printf("failed to update user name for user %d: %v", userID, err)
				response.DatabaseErrorResponse(c, err)
				return
			}
		}
		if req.RoleID != nil {
			_, err = tx.Exec("UPDATE users SET role_id = $1, updated_at = $2 WHERE id = $3", *req.RoleID, now, userID)
			if err != nil {
				log.Printf("failed to update user role for user %d: %v", userID, err)
				response.DatabaseErrorResponse(c, err)
				return
			}
		}
	}

	// Parse date of birth if provided
	var dob *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err != nil {
			response.BadRequest(c, "invalid date_of_birth format, expected YYYY-MM-DD")
			return
		}
		dob = &parsed
	}

	// Upsert profile
	profileQuery := `
		INSERT INTO user_profiles (user_id, first_name, last_name, middle_name, phone, student_id, room_id, avatar_url, date_of_birth, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (user_id)
		DO UPDATE SET
			first_name = COALESCE($2, user_profiles.first_name),
			last_name = COALESCE($3, user_profiles.last_name),
			middle_name = COALESCE($4, user_profiles.middle_name),
			phone = COALESCE($5, user_profiles.phone),
			student_id = COALESCE($6, user_profiles.student_id),
			room_id = COALESCE($7, user_profiles.room_id),
			avatar_url = COALESCE($8, user_profiles.avatar_url),
			date_of_birth = COALESCE($9, user_profiles.date_of_birth),
			updated_at = $10
	`

	_, err = tx.Exec(profileQuery,
		userID, req.FirstName, req.LastName, req.MiddleName, req.Phone,
		req.StudentID, req.RoomID, req.AvatarURL, dob, now,
	)
	if err != nil {
		log.Printf("failed to update profile for user %d: %v", userID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	if err := tx.Commit(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Fetch updated user
	user, err := models.GetUserWithProfile(h.DB, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			response.NotFoundError(c, "user not found")
			return
		}
		response.DatabaseErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    response.Updated,
		"message": "user updated",
		"data":    user,
	})
}
