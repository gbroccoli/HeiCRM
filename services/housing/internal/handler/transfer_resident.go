package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gbroccoli/HeiCRM/pkg/models"
	"github.com/gbroccoli/HeiCRM/pkg/response"
	"github.com/gin-gonic/gin"
)

// TransferResident moves a resident from one room to another
func (h *Handler) TransferResident(c *gin.Context) {
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

	residentID, err := strconv.ParseUint(c.Param("residentId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "Некорректный ID резидента")
		return
	}

	exists, err := roomExists(h.DB, buildingID, roomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !exists {
		response.NotFoundError(c, "Исходная комната не найдена")
		return
	}

	var req models.TransferResidentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequestError(c, "Некорректные данные", err)
		return
	}

	// Check target room exists
	targetExists, err := roomExists(h.DB, req.NewBuildingID, req.NewRoomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if !targetExists {
		response.NotFoundError(c, "Целевая комната не найдена")
		return
	}

	// Check target room capacity
	capacity, err := getRoomCapacity(h.DB, req.NewRoomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	occupancy, err := getRoomCurrentOccupancy(h.DB, req.NewRoomID)
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	if occupancy >= capacity {
		response.BadRequest(c, fmt.Sprintf("Целевая комната заполнена (%d/%d)", occupancy, capacity))
		return
	}

	// Check resident exists and is active in source room
	var fullName, birthDate, moveInDate string
	var passportSeries, passportNumber, email, phone *string
	err = h.DB.QueryRow(
		`SELECT full_name, birth_date, passport_series, passport_number, email, phone, move_in_date
		 FROM residents WHERE id = $1 AND room_id = $2 AND move_out_date IS NULL`,
		residentID, roomID,
	).Scan(&fullName, &birthDate, &passportSeries, &passportNumber, &email, &phone, &moveInDate)
	if errors.Is(err, sql.ErrNoRows) {
		response.NotFoundError(c, "Резидент не найден или уже выселен")
		return
	}
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Transaction: close old record + create new record
	tx, err := h.DB.Begin()
	if err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}
	defer tx.Rollback()

	moveOutDate := time.Now().Format("2006-01-02")
	newMoveInDate := time.Now().Format("2006-01-02")

	// Soft-close current record
	_, err = tx.Exec(
		"UPDATE residents SET move_out_date = $1, updated_at = NOW() WHERE id = $2",
		moveOutDate, residentID,
	)
	if err != nil {
		log.Printf("failed to close resident record: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Insert new record in target room
	var newResident models.Resident
	err = tx.QueryRow(
		`INSERT INTO residents (room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		 RETURNING id, room_id, full_name, birth_date, passport_series, passport_number, email, phone, move_in_date, move_out_date, created_at, updated_at`,
		req.NewRoomID, fullName, birthDate, passportSeries, passportNumber, email, phone, newMoveInDate,
	).Scan(&newResident.ID, &newResident.RoomID, &newResident.FullName, &newResident.BirthDate,
		&newResident.PassportSeries, &newResident.PassportNumber, &newResident.Email, &newResident.Phone,
		&newResident.MoveInDate, &newResident.MoveOutDate, &newResident.CreatedAt, &newResident.UpdatedAt)
	if err != nil {
		log.Printf("failed to create new resident record: %v", err)
		response.DatabaseErrorResponse(c, err)
		return
	}

	if err := tx.Commit(); err != nil {
		response.DatabaseErrorResponse(c, err)
		return
	}

	// Update room statuses for both rooms
	if err := updateRoomStatus(h.DB, roomID); err != nil {
		log.Printf("failed to update source room status: %v", err)
	}
	if err := updateRoomStatus(h.DB, req.NewRoomID); err != nil {
		log.Printf("failed to update target room status: %v", err)
	}

	response.SuccessCreated(c, newResident)
}
