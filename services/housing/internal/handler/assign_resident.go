package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// AssignResident creates a new resident and assigns them to a room
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

	var req models.CreateResidentRequest
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

	// Create resident
	var resident models.Resident
	err = h.DB.QueryRow(
		`INSERT INTO residents (room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date, move_out_date, created_at, updated_at`,
		roomID, req.FullName, req.BirthDate, req.PassportSeries, req.PassportNumber, req.Email, req.Phone, req.MoveInDate,
	).Scan(&resident.ID, &resident.RoomID, &resident.FullName, &resident.BirthDate,
		&resident.PassportSeries, &resident.PassportNumber, &resident.Email, &resident.Phone,
		&resident.MoveInDate, &resident.MoveOutDate, &resident.CreatedAt, &resident.UpdatedAt)

	if err != nil {
		log.Printf("failed to create resident: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Update room status
	if err := updateRoomStatus(h.DB, roomID); err != nil {
		log.Printf("failed to update room status: %v", err)
	}

	response.SuccessCreated(c, resident)
}
