package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AssignResident assigns a user to a room (sets user_profiles.room_id)
func (h *Handler) AssignResident(c *gin.Context) {
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

	exists, err := roomExists(h.DB, buildingID, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	var req models.AssignResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Check capacity
	capacity, err := getRoomCapacity(h.DB, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	occupancy, err := getRoomCurrentOccupancy(h.DB, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if occupancy >= capacity {
		response.BadRequest(c, fmt.Sprintf("Комната заполнена (%d/%d)", occupancy, capacity))
		return
	}

	// Check if user profile exists
	var profileExists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM user_profiles WHERE user_id = $1)", req.UserID).Scan(&profileExists)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !profileExists {
		response.NotFoundError(c, "Профиль пользователя не найден")
		return
	}

	// Check if user is already in this room
	var currentRoomID *uint64
	h.DB.QueryRow("SELECT room_id FROM user_profiles WHERE user_id = $1", req.UserID).Scan(&currentRoomID)
	if currentRoomID != nil && *currentRoomID == roomID {
		response.BadRequest(c, "Пользователь уже заселён в эту комнату")
		return
	}

	// Assign user to room
	result, err := h.DB.Exec(
		"UPDATE user_profiles SET room_id = $1, updated_at = NOW() WHERE user_id = $2",
		roomID, req.UserID,
	)
	if err != nil {
		log.Printf("failed to assign resident %d to room %d: %v", req.UserID, roomID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		response.NotFoundError(c, "Профиль пользователя не найден")
		return
	}

	response.SuccessOK(c, gin.H{
		"user_id": req.UserID,
		"room_id": roomID,
		"message": "Пользователь заселён",
	})
}
