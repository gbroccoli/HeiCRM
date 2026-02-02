package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	natsHandler "github.com/gbroccoli/HeiCRM/services/users/internal/nats"
	"github.com/gin-gonic/gin"
)

// DeleteUser allows admin to delete a user
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID пользователя")
		return
	}

	// Prevent self-deletion
	email, exists := c.Get("email")
	if exists {
		currentUserID, err := getUserIDByEmail(h.DB, email.(string))
		if err == nil && currentUserID == userID {
			response.Forbidden(c, "Нельзя удалить самого себя")
			return
		}
	}

	// Get user email before deletion for NATS event
	var userEmail string
	h.DB.QueryRow("SELECT email FROM users WHERE id = $1", userID).Scan(&userEmail)

	result, err := h.DB.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		log.Printf("failed to delete user %d: %v", userID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if rowsAffected == 0 {
		response.NotFoundError(c, "Пользователь не найден")
		return
	}

	// Publish NATS event
	natsHandler.PublishUserDeactivated(h.NC, userID, userEmail)

	c.JSON(http.StatusOK, gin.H{
		"code":    response.Deleted,
		"message": "Пользователь удалён",
	})
}
