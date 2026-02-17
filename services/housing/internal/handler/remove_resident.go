package handler

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// RemoveResident sets move_out_date for a resident (soft delete)
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

	var req models.RemoveResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	residentID := req.ResidentID

	exists, err := roomExists(h.DB, buildingID, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Комната не найдена")
		return
	}

	// Check if resident exists and belongs to this room
	var currentRoomID uint64
	err = h.DB.QueryRow("SELECT room_id FROM residents WHERE id = $1 AND move_out_date IS NULL", residentID).Scan(&currentRoomID)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Резидент не найден или уже выселен")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	if currentRoomID != roomID {
		response.BadRequest(c, "Резидент не проживает в этой комнате")
		return
	}

	// Set move_out_date
	moveOutDate := time.Now().Format("2006-01-02")
	var resident models.Resident
	err = h.DB.QueryRow(
		`UPDATE residents SET move_out_date = $1, updated_at = NOW()
		 WHERE id = $2
		 RETURNING id, room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date, move_out_date, created_at, updated_at`,
		moveOutDate, residentID,
	).Scan(&resident.ID, &resident.RoomID, &resident.FullName, &resident.BirthDate,
		&resident.PassportSeries, &resident.PassportNumber, &resident.Email, &resident.Phone,
		&resident.MoveInDate, &resident.MoveOutDate, &resident.CreatedAt, &resident.UpdatedAt)

	if err != nil {
		log.Printf("failed to remove resident %d: %v", residentID, err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Update room status
	if err := updateRoomStatus(h.DB, roomID); err != nil {
		log.Printf("failed to update room status: %v", err)
	}

	response.SuccessOK(c, resident)
}
