package handler

import (
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// RemoveResident removes a user from a room (sets user_profiles.room_id = NULL)
func (h *Handler) RemoveResident(c *gin.Context) {
	buildingID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID здания")
		return
	}

	roomID, err := strconv.ParseUint(c.Param("roomId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID комнаты")
		return
	}

	userID, err := strconv.ParseUint(c.Param("userId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID пользователя")
		return
	}

	exists, err := roomExists(h.DB, buildingID, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	// Remove user from room only if they are in this exact room
	result, err := h.DB.Exec(
		"UPDATE user_profiles SET room_id = NULL, updated_at = NOW() WHERE user_id = $1 AND room_id = $2",
		userID, roomID,
	)
	if err != nil {
		log.Printf("failed to remove resident %d from room %d: %v", userID, roomID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.NotFoundError(c, "Пользователь не найден в этой комнате")
		return
	}

	response.SuccessDeleted(c, nil)
}
